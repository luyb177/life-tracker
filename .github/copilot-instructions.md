# Copilot Instructions

## Build & Code Generation

```bash
# Build the backend
cd backend && go build ./...

# Sync dependencies
cd backend && go mod tidy

# Generate code from API definition (run after editing .api file)
cd backend && goctl api go -api life_tracker.api -dir . --style go_zero
# Note: goctl skips existing files — manually remove files you want regenerated first.
# After generation, delete the stub lifetracker.go if life_tracker.go already exists.
```

## Architecture

This is a go-zero REST API backend for a personal life tracker (activity & expense logging).

### Layered pattern (handler → logic → repo)

Every API endpoint follows this flow:
- **`internal/handler/`** — HTTP layer. Thin: parses request, calls logic, writes `respx.OkCtx`/`respx.ErrorCtx`. goctl-generated stubs.
- **`internal/logic/`** — Business logic. Validates input, orchestrates repo calls, handles auth flows.
- **`internal/repo/`** — Data access. Each sub-package (user, verify, token) has its own `Repository` interface and `repo` struct (GORM + Redis). The `repo.go` file aggregates all into `Repositories` with a `Transaction()` helper.

### ServiceContext (DI container)

`internal/svc/service_context.go` wires all dependencies at startup: MySQL, Redis, Mailer, JWT handler, EmailSender, Repos, and both middlewares (IP, JWT). It panics on failure — the app won't start without its dependencies.

### API-first with goctl

`life_tracker.api` is the single source of truth for types and routes. **Always add new request/response types here first**, then run `goctl api go` to regenerate `types.go`, handler stubs, and `routes.go`. goctl will NOT overwrite existing files — delete files you want refreshed before running it.

## Key Conventions

### Error handling — unified envelope

All API responses use the common envelope `{code, msg, data}` via `common/respx` and `common/errorx`:
- Success: `respx.OkCtx(ctx, w, data)` → `{code: 0, msg: "ok", data: ...}`
- Error: `respx.ErrorCtx(ctx, w, err)` — expects `*errorx.AppError` with a business code (400–xxxxxx)

Logic functions return `(*types.SomeResp, error)` and NEVER write HTTP responses directly. The handler is the only layer that calls respx.

### JWT refresh token rotation

Access tokens expire in 15min. Refresh tokens expire in 7 days but are **rotated on every use**:
- Each token contains a unique `JTI` (hex UUID, generated in `common/jwtx`)
- On login/refresh, the JTI is stored in Redis (`refresh:{userID}`) with the refresh token's TTL
- On refresh: validate JTI matches Redis → generate new tokens with new JTIs → store new JTI (overwrites old, immediately invalidating the previous refresh token)
- On password change: Redis key is deleted, revoking all refresh tokens

### Config files — never commit real credentials

All `*.yaml` files are gitignored EXCEPT `*.example.yaml`. The example file uses placeholder values. The actual `life_tracker.yaml` is developer-local.

### Middleware files

goctl generates stub middleware files (e.g., `jwtmiddleware.go`). The real implementation lives in `*_middleware.go` (e.g., `jwt_middleware.go`). When goctl regenerates, it creates fresh stubs — keep the real impl files and delete the goctl stubs if they collide.

### Verification codes — atomic Lua

Verification codes are stored in Redis with a Lua script (`repo/verify/lua/verify.lua`) that atomically checks and deletes the code, preventing replay attacks.

### Data files

IP geo-location databases (`data/ip2region_v*.xdb`) are committed to the repo and loaded by the IP middleware at startup.
