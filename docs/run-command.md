 # Run Commands

A cheatsheet of every command you'll need while developing on Slugify.

## Prerequisites

- Go 1.25+
- Node.js 20+
- Docker (optional but easier than installing Postgres + Redis locally)
- `kubectl` (only if you want to try the Kubernetes setup)

---

## Start the stack with Docker (easiest)

```bash
docker compose up --build
docker compose down            # stop
docker compose down -v         # stop AND wipe the postgres volume
```

URLs:
- API → <http://localhost:9000>
- Web → <http://localhost:3000>

---

## Start the stack manually

### 1. Backend (Go API)

```bash
cd api
cp .env.example .env
# Edit .env to fill:
#   PORT=9000
#   APP_NAME=slugify
#   DATABASE_URL=postgres://postgres:postgres@localhost:5432/slugify?sslmode=disable
#   REDIS_URL=redis://localhost:6379
#   DOMAINURL=http://localhost:9000
#   JWT_SECRET=any-long-random-string

# If you don't have Postgres/Redis locally, run just those two via compose:
docker compose up postgres redis -d

go run ./cmd/server
```

Migrations run automatically on boot.

### 2. Frontend (Next.js)

In a second terminal:

```bash
cd web
npm install

# PowerShell:
$env:BACKEND_URL="<http://localhost:9000>"
# bash / zsh:
# export BACKEND_URL=http://localhost:9000

npm run dev
```

Open <http://localhost:3000>.

The Next.js dev server proxies:
- `/backend/*` → `${BACKEND_URL}/api/*` (X-API-Key routes)
- `/auth-api/*` → `${BACKEND_URL}/auth/*` (JWT routes)

So you never hit CORS in dev.

---

## Tests

### Go (backend)

```bash
cd api
go test ./... -count=1            # local
go test ./... -count=1 -race      # in CI (Linux)
go test ./... -count=1 -v         # verbose, prints each test name
```

The race detector requires CGO + a 64-bit C toolchain. On Windows with the
default MinGW, drop `-race` locally — CI on Linux still runs it.

### Lint and build (frontend)

```bash
cd web
npm run lint
npm run build
```

There are no React component tests yet (on the roadmap).

---

## Database

### Apply / roll back migrations manually

Migrations apply automatically on `go run ./cmd/server`. To run them by hand:

```bash
cd api
# Install once
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

migrate -path ./migrations -database "$env:DATABASE_URL" up
migrate -path ./migrations -database "$env:DATABASE_URL" down 1
migrate -path ./migrations -database "$env:DATABASE_URL" version
```

### Connect to Postgres (Docker)

```bash
docker compose exec postgres psql -U postgres slugify
```

Common queries:
```sql
\\dt                          -- list tables
SELECT id, email FROM users;
SELECT short_code, clicks FROM links ORDER BY clicks DESC LIMIT 5;
TRUNCATE links;              -- wipe links during dev
```

### Connect to Redis (Docker)

```bash
docker compose exec redis redis-cli
KEYS url:*
GET url:abc
FLUSHALL                      -- wipe cache during dev
```

---

## Adding a new endpoint

1. Define the request/response in `internal/models/*.go`.
2. Add the repository method to the interface in `internal/db/repository.go`.
3. Implement it in `internal/db/postgres.go`.
4. Add the business logic in `internal/services/*.go`.
5. Write the handler in `internal/handlers/*.go`. Read `user_id` from context.
6. Register the route in `cmd/server/main.go` under the right group
   (`/auth/*` JWT or `/api/*` X-API-Key).
7. Add the matching client function in `web/src/lib/auth-api.ts` or
   `web/src/lib/slugify-api.ts`.

## Adding a new migration

```bash
cd api
# Use the next free 6-digit slot (e.g. 000009)
echo "ALTER TABLE links ADD COLUMN expires_at TIMESTAMP;" > migrations/000009_add_expires_to_links.up.sql
echo "ALTER TABLE links DROP COLUMN expires_at;"          > migrations/000009_add_expires_to_links.down.sql
```

Restart the server. Migrations apply on boot.

---

## Kubernetes (optional)

```bash
# Build local images
cd api && docker build -t slugify-api:dev . && cd ..
cd web && docker build -t slugify-web:dev . && cd ..

# Apply all manifests
kubectl apply -f k8s/
kubectl config set-context --current --namespace=slugify

# Watch pods come up
kubectl get pods -w

# Logs / shell
kubectl logs deploy/api --tail=100 -f
kubectl exec -it deploy/postgres -- psql -U postgres slugify

# Tear down
kubectl delete -f k8s/
```

---

## CI

`.github/workflows/ci.yml` runs on every push and PR:

- **api** job: `gofmt`, `go vet`, `go test -race`, `go build`
- **web** job: `npm ci` (or `npm install` fallback), `npm run lint`,
  `npm run build`
- **docker** job: builds both images (does not push)

Run the same checks locally before pushing:

```bash
cd api && gofmt -l . && go vet ./... && go test ./... -count=1 && go build ./...
cd ../web && npm ci && npm run lint && npm run build
```