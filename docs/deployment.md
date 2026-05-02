<!--
# Deployment

## Docker Compose (single host)

```bash
docker compose up --build
docker compose down -v        # tear down, removes the postgres volume
```

Services exposed:

| Service | Port | URL |
| --- | --- | --- |
| `web` | 3000 | <http://localhost:3000> |
| `api` | 9000 | <http://localhost:9000> |
| `postgres` | 5432 | internal + 5432 published |
| `redis` | 6379 | internal + 6379 published |

`docker-compose.yml` wires:

- `web` → `BACKEND_URL=http://api:9000` (Docker DNS)
- `api` → `DATABASE_URL=postgres://...@postgres:5432/...`,
  `REDIS_URL=redis://redis:6379`
- Healthchecks on Postgres and Redis so the API doesn't start before they're
  ready

### Image sizes

| Image | Approx size |
| --- | --- |
| `slugify-api` | ~15 MB (multi-stage → distroless) |
| `slugify-web` | ~150 MB (Next.js standalone on Node 20-alpine) |

---

## Kubernetes (work in progress)

Manifests live in `k8s/` and apply with:

```bash
# 1. Build local images so the cluster can pull them
cd api && docker build -t slugify-api:dev . && cd ..
cd web && docker build -t slugify-web:dev . && cd ..

# 2. Apply manifests in order
kubectl apply -f k8s/
kubectl get pods -n slugify -w
```

Once everything is `Running`, install ingress-nginx and open
`http://slugify.localhost`.

### Layout

```
k8s/
  00-namespace.yaml             # Namespace
  01-config.yaml                # ConfigMap (non-secret config)
  02-secrets.yaml               # Secret (DB URL, JWT secret)
  10-postgres-pvc.yaml          # PersistentVolumeClaim
  11-postgres-deployment.yaml   # Postgres pod
  12-postgres-service.yaml      # ClusterIP svc
  20-redis-deployment.yaml
  21-redis-service.yaml
  30-api-deployment.yaml        # 2 replicas, envFrom ConfigMap+Secret
  31-api-service.yaml
  40-web-deployment.yaml
  41-web-service.yaml
  50-ingress.yaml               # nginx Ingress, host slugify.localhost
```

### Concepts each file demonstrates

| File | Concept |
| --- | --- |
| `00-namespace.yaml` | Namespace isolation |
| `01-config.yaml` / `02-secrets.yaml` | ConfigMap + Secret |
| `10-postgres-pvc.yaml` | Persistent Volume Claim |
| `11-...-deployment.yaml` | Deployment + readiness probes |
| `12-...-service.yaml` | ClusterIP Service + label selectors |
| `30-api-deployment.yaml` | `envFrom` injection, replicas, probes |
| `50-ingress.yaml` | Ingress rules, host routing |

---

## Environment variables

| Var | Required | Used by | Purpose |
| --- | --- | --- | --- |
| `PORT` | yes | API | Bind port |
| `APP_NAME` | no | API | Logging label (default `url-shortener`) |
| `DATABASE_URL` | yes | API | Postgres connection string |
| `REDIS_URL` | yes | API | Redis connection string |
| `DOMAINURL` | yes | API | Used to build the returned `short_url` |
| `JWT_SECRET` | yes | API | Symmetric secret for HS256 JWT signing |
| `BACKEND_URL` | yes | Web | Where the Next.js proxy forwards to |

---

## Production checklist (when you're ready)

- [ ] Replace `JWT_SECRET` with 32+ random bytes from a secret manager
- [ ] Put TLS in front (nginx, Caddy, or an Ingress with cert-manager)
- [ ] Set `DOMAINURL` to your real domain
- [ ] Disable `imagePullPolicy: Never` and push images to a registry
- [ ] Wire `/metrics` (roadmap) to Prometheus
- [ ] Add a backup strategy for Postgres (`pg_basebackup`, snapshots, etc.)
- [ ] Switch from in-memory token bucket logging to structured JSON logs

-->