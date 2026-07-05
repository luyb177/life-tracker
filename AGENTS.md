# AGENTS.md

This file provides guidance to Codex (Codex.ai/code) when working with code in this repository.

## Commands

```bash
# Build / check compilation
cd backend && go build ./...

# Sync dependencies
cd backend && go mod tidy

# Run the API server
cd backend && go run life_tracker.go -f etc/life_tracker.yaml

# Generate code from API definition (AFTER editing life_tracker.api)
cd backend && goctl api go -api life_tracker.api -dir . --style go_zero
# WARNING: goctl does NOT overwrite existing files.
# Delete files you want regenerated BEFORE running goctl, or you'll get duplicate symbols.
# Also delete the goctl-generated stub `lifetracker.go` if `life_tracker.go` already exists.

# Run a cron summary job manually
cd backend && go run cmd/cron/main.go cron summary -t 1   # 1=day 2=week 3=month 4=year
```

There are currently no tests (`*_test.go` files do not exist).

## Module & Stack

- Module: `github.com/luyb177/life-tracker/backend`
- Go 1.26, go-zero REST framework, GORM + MySQL, Redis, DeepSeek AI
- Frontend (Vue 3) is in a separate repository, not in this workspace.

## Architecture

### Handler → Logic → Repo layering

Every endpoint follows this pattern:
- **`internal/handler/`** — Thin HTTP layer. Parses request, calls logic, writes response via `respx.OkCtx`/`respx.ErrorCtx`. goctl-generated stubs.
- **`internal/logic/`** — Business logic. Validates input, orchestrates repo calls. Returns `(*types.SomeResp, error)` — never writes HTTP responses.
- **`internal/repo/`** — Data access. Each sub-package defines a `Repository` interface and `repo` struct (GORM + Redis). `repo/repo.go` aggregates all into `Repositories` with a `Transaction()` helper.

### ServiceContext (DI container)

`internal/svc/service_context.go` wires all dependencies at startup and **panics on failure**. Everything flows through it: MySQL, Redis, Mailer, JWT handler, EmailSender, Repos, and both middlewares.

### API-first with goctl

`life_tracker.api` is the **single source of truth** for types and routes. New endpoints MUST start here. Then run `goctl api go` to regenerate `types.go`, handler stubs, and `routes.go`. The `.api` file uses goctl's DSL — types get compiled into `internal/types/types.go`.

### Two entry points

| Entry | File | Purpose |
|-------|------|---------|
| API server | `life_tracker.go` | REST API on `:8888` |
| Cron scheduler | `cmd/cron/main.go` | Cobra CLI with `summary` subcommand for AI summary jobs |

Both share the same `ServiceContext` and config file.

## Key Design Decisions

### Life logs and summaries are separate layers

Daily activity records live in the `life_logs` table and are exposed through `/api/v1/life_log/*`.
The `summaries` table is only for periodic review content:
- `source=2` (user): manual period summaries. One user summary per user+period_type+period_start.
- `source=1` (AI): auto-generated period summaries. One AI summary per user+period_type+period_start.
- `last_updated_by` / `last_updated_at` record who last changed the content and when. Scheduled AI jobs use `last_updated_by=0`.

### Summary period model

All `period_start`/`period_end` use `YYYY-MM-DD` strings. `period_end` is an **open interval** (exclusive end date).

| PeriodType | Meaning | period_start rule | period_end rule |
|------------|---------|-------------------|-----------------|
| 1 | 日报 | Any date | `start + 1 day` |
| 2 | 周报 | Must be Monday | `start + 7 days` |
| 3 | 月报 | Must be 1st of month | 1st of next month |
| 4 | 年报 | Must be Jan 1 | Next Jan 1 |

PeriodType 5 (人生总结) was intentionally removed — do not add it back.

### AI summary context hierarchy

Higher-level summaries consume lower-level ones as context:
```
日报 ← (no child)
周报 ← 日报
月报 ← 周报
年报 ← 月报
```

The core AI summary flow lives in `internal/logic/cron/summary.go:Run()`. It gathers expense data, location breakdown, user journals, and child-period summaries, then sends them to DeepSeek.

### JWT refresh token rotation

- Access tokens: short-lived (default 2h).
- Refresh tokens: longer-lived (default 7d), **rotated on every use**.
- Each token has a unique JTI. Redis stores the current valid JTI at `refresh:{userID}`.
- On refresh: old JTI is atomically replaced (Lua script at `repo/token/lua/rotate.lua`). The old refresh token is immediately invalid.
- Password change deletes the Redis key, revoking all refresh tokens.

### Verification codes

Codes are stored in Redis with SHA256-hashed target keys. A Lua script (`repo/verify/lua/verify.lua`) atomically checks and deletes the code to prevent replay. Configurable cooldown via TTL.

## Error Handling

All responses use the unified envelope `{code, msg, data}` via `common/respx`:
- `code=0` = success
- Non-zero = business error (codes defined in `common/errorx/errors.go`)

Logic returns `(*Resp, error)` where `error` is `*errorx.AppError`. Only handlers call `respx.OkCtx`/`respx.ErrorCtx`.

## Important Conventions

- **Never commit real credentials.** Config files (`*.yaml`) except `*.example.yaml` are gitignored. Use `etc/life_tracker.example.yaml` for defaults.
- **Middleware naming:** goctl generates stubs like `jwtmiddleware.go`. The real impl lives in `jwt_middleware.go`. After goctl regeneration, delete goctl stubs that collide with real files.
- **Soft delete:** All GORM models use `soft_delete` plugin with nanosecond-precision delete markers.
- **Cursor pagination:** Expense and summary list endpoints use cursor pagination (page_token), not offset.
- **Tags:** Summary tags are comma-separated strings (`varchar(500)`), not a normalized relation. Split/trim on the frontend.
- **Timezone:** Database uses `Asia/Shanghai`. Config DSN must include `loc=Asia%2FShanghai`.

## Frontend Notes

The frontend (Vue 3 + Naive UI) is being rewritten. The backend API is stable. Key frontend integration points:
- `GET /api/v1/summary/day` returns mixed user+AI records — group by `source` when displaying.
- For the same user+period_type+period_start, there can be at most one record per source.
- For week/month/year summaries that already exist, guide users to edit instead of re-creating.
