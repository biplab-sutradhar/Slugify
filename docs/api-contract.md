# API Contract

Base URL (local): `http://localhost:9000`

All request/response bodies are JSON. Errors are returned as
`{ "error": "<message>" }` with an appropriate status code.

## Auth schemes

Two schemes are supported on different route groups.

| Group | Header | Used by |
| --- | --- | --- |
| `/auth/*` (protected) | `Authorization: Bearer <jwt>` | Dashboard |
| `/api/*` | `X-API-Key: <key>` | External integrations |

The dashboard talks to `/auth/*` for everything except the existing
`/backend/*` Next.js proxy, which targets `/api/*` using the auto-provisioned
key.

---

## Public endpoints

### `POST /auth/register`

Create a new user. Returns a JWT and an auto-minted default API key.

Request:
```json
{ "email": "alice@example.com", "password": "password123", "name": "Alice" }
```

Response `201`:
```json
{
  "token": "<jwt>",
  "user": {
    "id": "...",
    "email": "alice@example.com",
    "name": "Alice",
    "created_at": "..."
  },
  "api_key": "<43-char base64url string>"
}
```

Errors: `400` invalid body, `409` email already registered.

### `POST /auth/login`

Returns a JWT for an existing user. **No `api_key` in the response** — use
`POST /auth/api-key` if the dashboard needs one.

Request:
```json
{ "email": "alice@example.com", "password": "password123" }
```

Response `200`:
```json
{ "token": "<jwt>", "user": { ... } }
```

Errors: `401` invalid credentials.

### `GET /:shortCode`

Public redirect. Looks up the short code in Redis (then Postgres), returns
`302 Found` with the original URL, and increments the click counter
asynchronously.

Errors: `404` if the code is unknown.

### `GET /health`

Liveness probe. Always returns `200 {"status":"ok"}`.

---

## JWT-protected endpoints

Header on every request: `Authorization: Bearer <jwt>`

### `GET /auth/me`

Returns the current user.

Response `200`:
```json
{ "id": "...", "email": "...", "name": "...", "created_at": "..." }
```

### `POST /auth/api-key`

Mints a new API key for the current user. Used by the dashboard when a
logged-in user has no key yet.

Request: `{ "name": "Mobile" }` (optional, defaults to `"Default"`).

Response `201`: `{ "api_key": "<key>" }`

### `GET /auth/keys`

List the caller's API keys.

### `POST /auth/keys`

Create a key (alternative to `/auth/api-key`, used by the keys page).

Request: `{ "name": "...", "scope": "default" }`

### `DELETE /auth/keys/:id`

Revoke a key by ID. The key must belong to the caller, otherwise the server
returns the same response as a missing key (no enumeration).

---

## API-key-protected endpoints

Header on every request: `X-API-Key: <key>`

Rate-limited per key (100 req/min, refilled at 1/sec). Returns `429` once
the bucket is empty.

### `POST /api/shorten`

Request: `{ "long_url": "https://..." }`

Response `201`: `{ "short_url": "<http://localhost:9000/><code>" }`

### `GET /api/links?limit=20&offset=0`

List the caller's links. Default `limit=20`, max `100`.

Response `200`:
```json
[
  {
    "id": "...",
    "user_id": "...",
    "short_code": "abc",
    "long_url": "<https://example.com>",
    "is_active": true,
    "clicks": 42,
    "created_at": "..."
  }
]
```

### `GET /api/links/:id`

Get a single link by ID (must be owned by the caller).

### `PATCH /api/links/:id`

Activate or deactivate a link.

Request: `{ "is_active": true }`

### `DELETE /api/links/:id`

Permanently delete a link (only if owned by the caller).

### `GET /api/keys` / `POST /api/keys` / `DELETE /api/keys/:id`

Same semantics as the JWT versions. Useful for managing keys from external
scripts that already hold an admin key.

---

## Error format

```json
{ "error": "human-readable message" }
```

Common status codes:

| Code | Meaning |
| --- | --- |
| 400 | Malformed request body |
| 401 | Missing / invalid auth |
| 404 | Resource not found OR not owned by the caller |
| 409 | Conflict (e.g. duplicate email on register) |
| 429 | Rate limit exceeded |
| 500 | Server error |