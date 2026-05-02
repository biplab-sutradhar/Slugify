# Slugify

[![CI](<https://github.com/your-username/slugify/actions/workflows/ci.yml/badge.svg>)](<https://github.com/your-username/slugify/actions/workflows/ci.yml>)
![Go](<https://img.shields.io/badge/Go-1.25-00ADD8?logo=go&logoColor=white>)
![Next.js](<https://img.shields.io/badge/Next.js-15-000?logo=nextdotjs>)
![Postgres](<https://img.shields.io/badge/Postgres-16-336791?logo=postgresql&logoColor=white>)
![Redis](<https://img.shields.io/badge/Redis-7-DC382D?logo=redis&logoColor=white>)
![License](<https://img.shields.io/badge/license-MIT-blue.svg>)

A self-hostable URL shortener with a polished admin dashboard. Built end-to-end
to demonstrate full-stack engineering: Go services, Postgres + Redis, JWT auth,
multi-tenant data isolation, a Next.js dashboard, Docker, and CI.

> **Note:** this is a learning / portfolio project. The architecture is real,
> but performance numbers come from local benchmarks — not production traffic.

## Features

**Core**
- Short links with sub-10 ms Redis cache hits and Postgres fallback
- Custom **base62 ID generator** backed by a range-based ticket server
- **Per-user data isolation** — every link and API key is scoped to its owner
- **Two coexisting auth schemes**:
  - JWT (HS256, bcrypt) for the dashboard
  - API keys (`X-API-Key`) for external integrations
- **Token-bucket rate limiting** per API key (Redis)

**Dashboard**
- Modern Next.js 15 + Tailwind v4 UI
- Landing → signup → dashboard with auto-redirect when authenticated
- Pages: Overview, Links, API keys, Analytics
- Auto-provisioned API key on signup, manual mint on demand

**Operations**
- Multi-stage Docker builds (~15 MB API image, distroless)
- `docker compose` for the full stack
- GitHub Actions CI: `gofmt` / `vet` / `go test` + `next lint` + `next build` + Docker image build

## Tech stack

| Layer | Stack |
| --- | --- |
| API | Go 1.25, Gin, `golang-migrate`, `golang-jwt/v5`, `bcrypt` |
| Storage | PostgreSQL 16, Redis 7 |
| Frontend | Next.js 15 (App Router), React 19, Tailwind CSS v4, TypeScript |
| Infra | Docker, docker-compose, GitHub Actions |
| Tests | Go `testing` + table-driven unit tests with hand-rolled fakes |

## Quick start

```bash
docker compose up --build
```

- API → <http://localhost:9000> (`GET /health` for a sanity check)
- Web → <http://localhost:3000>

For non-Docker setup and every other command you'll ever need, see
[`docs/run-command.md`](docs/run-command.md).

## Repository layout

```
.
├── api/                       # Go backend
│   ├── cmd/server/            # main entry point
│   ├── internal/
│   │   ├── auth/              # JWT, API-key generation
│   │   ├── cache/             # Redis client + Cache interface
│   │   ├── config/            # env loading
│   │   ├── db/                # Postgres repositories + interfaces
│   │   ├── handlers/          # Gin HTTP handlers
│   │   ├── idgen/             # base62 + ticket server
│   │   ├── middleware/        # X-API-Key + JWT middlewares
│   │   ├── models/            # domain types
│   │   └── services/          # business logic, no HTTP / DB code
│   ├── migrations/            # golang-migrate SQL files
│   └── Dockerfile
│
├── web/                       # Next.js frontend
│   ├── src/
│   │   ├── app/               # App Router pages (route groups)
│   │   ├── components/        # UI + feature components
│   │   └── lib/               # API clients + auth context
│   └── Dockerfile
│
├── k8s/                       # Kubernetes manifests (work in progress)
├── docs/                      # Architecture, API, deployment, commands
├── docker-compose.yml
└── .github/workflows/ci.yml
```

## Documentation

- **[Architecture](docs/architecture.md)** — request flow, ID generation, multi-tenancy, schema
- **[API Contract](docs/api-contract.md)** — every REST endpoint with auth + payload
- **[Deployment](docs/deployment.md)** — Docker Compose and Kubernetes
- **[Run Commands](docs/run-command.md)** — local setup, testing, lint, migrations

## Roadmap

- [ ] Real time-series click events table (currently a counter)
- [ ] Link expiration (TTL)
- [ ] OpenAPI / Swagger spec at `/docs`
- [ ] Graceful shutdown on `SIGTERM`
- [ ] Mobile sidebar for the dashboard
- [ ] Helm chart for the Kubernetes setup
- [ ] Prometheus `/metrics` endpoint

## License

MIT — see [`LICENSE`](LICENSE).