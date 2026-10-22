# task_api vs. real service repos: gap analysis and roadmap

## Summary

The learning project already uses the same core stack as the real repos: Echo v5, GORM + Postgres, go-redis v9, and slog. The gaps are not about libraries. They are about the scaffolding that surrounds the business code: config loading, interface-based DI, a central error/response envelope, versioned migrations, the middleware stack (request ID, request logging, recover, timeout, metrics, tracing), graceful shutdown, mock-based unit tests, linting, and CI.

Fastest route to "looks like a real repo": imitate **timeproduct_be** for layout and wiring, and **prowifi-be** for error handling, migrations, and the Makefile. Both are the cleanest examples I read.

Repos read in depth: timeproduct_be, prowifi-be, tapsales-be, noblock-be. Repos sampled: timevoucher-be, backoffice-v2-be, selfcare-be (go.mod and layout only), terraform, basecamp-fe.

## (a) House style observed

| Concern | Common pattern | Evidence |
|---|---|---|
| Layout | `cmd/<bin>/main.go`, `config/`, `internal/{handlers,services,data_access,entities,models,middleware,clients,database,utils/{response,errors,validator}}`. One package per domain, often one file per operation. | timeproduct_be/internal/*, noblock-be/internal/*, tapsales-be/internal/* |
| Layering | handler → service → data_access/repository (+ clients). Every layer is an interface (`IPromotionsService`, `IPromotionsDB`) with a `NewX(...) IX` constructor. | timeproduct_be/internal/services/promotions/promotions.go, data_access/promotions/promotions.go; prowifi-be/.kiro/steering/structure.md ("Layering Rule") |
| DI wiring | Aggregator structs (`DataAccess`, `Services`, `Handlers`) built in one place. Each handler has `RegisterRoutes(*echo.Group)`. prowifi-be uses a `server.Deps` struct wired in `cmd/app/api.go`. | timeproduct_be/internal/route/route.go, handlers/handlers.go, services/services.go |
| HTTP | Echo v5 (labstack/echo/v5 v5.2–v5.4). Older repos are still on v4 (backoffice-v2-be, timevoucher-be). | go.mod of all repos |
| DB | GORM v1.30–1.31 + gorm.io/driver/postgres, `gorm.io/plugin/dbresolver` for a read replica, connection pool settings in config, a custom GORM→slog logger. | timeproduct_be/internal/database/postgres.go, database/gorm_logger.go; config/config.go `DbDetail` |
| Migrations | Newer repos use goose v3 with SQL files embedded via `embed.FS`, plus a separate `migrate` binary/subcommand that runs before the server. Older repos use `AutoMigrate`. | timeproduct_be/cmd/migrate/main.go, internal/database/migrations/embed.go + sql/0000N_*.sql; prowifi-be/cmd/app/migrate.go, migrations/; AutoMigrate in noblock-be/internal/databases/database.go, selfcare-be/internal/database/postgres.go |
| Config | spf13/viper reads `config.yml`/`config.json` plus `AutomaticEnv`, maps into typed structs, uses a singleton `GetConfig()`. A `config.example.*` file is committed and the real one is gitignored. | timeproduct_be/config/config.go, config.example.json; tapsales-be/config.example.yml |
| Logging | slog JSON handler wrapped in a context-aware handler (private org module `gotime/observability/log`), so request_id/trace_id get attached to every `*Context` log call. tapsales/timevoucher use logrus/lumberjack (legacy). | timeproduct_be/internal/observability/log.go |
| Middleware stack | echoprometheus → OTel tracing → RequestID (uuid) that injects IDs into the ctx logger → RequestLogger (skips health) → Recover → Timeout → CORS/BodyLimit. | timeproduct_be/internal/server/server.go `Start`; prowifi-be/internal/server/server.go `New` |
| Metrics / tracing | Separate Prometheus `/metrics` server on its own port (default 9999). OTel via the internal time-trace/gotime trace module. Terraform's ECS module has `otel_collector_enabled`. | timeproduct_be/internal/metrics/server.go; prowifi-be/cmd/app/api.go; terraform/svc/timeproduct/base/main.tf |
| Errors / response | A typed app error carrying status + code + message, and one place that writes the error JSON: `response.Make` wrapper (timeproduct), a central `HTTPErrorHandler` (prowifi), or `AppError{Type,Code}` (tapsales). JSON envelopes like `{statusCode,message,data}` or `{success,data,metadata}`. | timeproduct_be/internal/utils/response/response.go; prowifi-be/internal/apperror/apperror.go, internal/server/errors.go; tapsales-be/internal/utils/errors/errors.go; noblock-be/pkg/response/response.go |
| Validation | go-playground/validator v10 registered as `echo.Validator`, `c.Bind` then `c.Validate`, custom rules and translations. | timeproduct_be/internal/server/server.go `CustomValidator`; tapsales-be/internal/utils/validator/validator.go |
| Auth | API-key middleware backed by a DB table and an in-memory TTL cache (timeproduct, prowifi). JWT HS256 access/refresh tokens stored hashed in Redis (tapsales). Keycloak JWT + Casbin RBAC (noblock, backoffice). | timeproduct_be/internal/middleware/auth_middleware.go; tapsales-be/pkg/jwt/jwt.go, internal/middleware/access_token.go; noblock-be/internal/middlewares/casbin_middleware.go, casbin_model.conf |
| Data | Pagination via `?page&page_size` (capped at 100) with returned metadata, sort allowlisting to prevent SQL injection, `CreatedAt/UpdatedAt/DeletedAt` (soft delete), `WithContext(ctx)` on every query, Postgres error code mapping (23505 → domain error). | noblock-be/pkg/helper/pagination.go, sorting.go; timeproduct_be/internal/models/promotion.go; prowifi-be/internal/repository/order.go `Create` |
| Outbound HTTP | A wrapped `HTTPClient` interface with functional options (timeout, logger), request-ID forwarding, structured request logs, and exponential-backoff retry helpers. | tapsales-be/pkg/httpclient/httpclient.go, internal/utils/retry/retry.go; prowifi-be/internal/httpx/ |
| Async | Separate worker/scheduler binaries. Postgres-backed queue worker with ack/nack/retry (prowifi). AWS SQS consumer + robfig/cron scheduler (timevoucher). | prowifi-be/internal/worker/worker.go, cmd/app/worker/; timevoucher-be/internal/worker/queue.go, internal/scheduler/scheduler.go, Dockerfile.worker |
| Testing | testify (assert/require/mock). mockery v2 (`.mockery.yml`, testify template) generates mocks for every interface. Handler tests use `httptest` + `e.NewContext` with a mocked service. Client tests use `httptest.NewServer`. Integration tests use in-memory SQLite (noblock). | timeproduct_be/.mockery.yml, internal/mocks/, handlers/promotions/promotions_test.go, services/promotions/promotions_test.go; prowifi-be/.kiro/steering/testing.md; noblock-be/internal/testutil/testdb.go |
| Tooling | `.golangci.yml` v2 (errcheck, govet, staticcheck, gosec, unused, bodyclose, …), `.air.toml` hot reload, Makefile (run/dev/build/test -race -cover/lint/migrate-*), git pre-commit lint hook, AGENTS.md/.kiro steering docs. | timeproduct_be/.golangci.yml; prowifi-be/Makefile; timevoucher-be/scripts/pre-commit, Makefile |
| CI/CD | bitbucket-pipelines.yml: unit tests → go vet + golangci-lint → docker build (arm64) → push to ECR via OIDC. Config is pulled from SSM Parameter Store at deploy time. | timeproduct_be/bitbucket-pipelines.yml |
| Docker | Multi-stage: golang builder (go.mod copied first for cache, `CGO_ENABLED=0 -trimpath`) → alpine runtime + tzdata. An entrypoint runs `migrate && server`. | timeproduct_be/Dockerfile, docker-entrypoint.sh; prowifi-be/Dockerfile |
| Infra | One terraform root per service (`svc/<name>/{base,stg,uat,prod}`) using the shared `modules/aws/ecs/v4` module: ECS service behind an ALB, health-check path, autoscaling, OTel sidecar, config from SSM. | terraform/svc/timeproduct/base/main.tf, terraform/AGENTS.md |

## (b) Gap table

| Area | Real repos do | task_api does | Suggested change |
|---|---|---|---|
| Config | viper + typed `Config` + `config.example.*` | `config/config.go` is empty. DSN and Redis address are hardcoded in `internal/database/database.go`. The `DATABASE_URL` set in compose.yaml is never read. | Implement viper config with server/db/redis/log/timeout sections. Commit `config.example.yml`, gitignore `config.yml`. |
| go.mod | Direct deps listed directly | Every dependency is marked `// indirect`, and both echo v4 and v5 are present | Run `go mod tidy`. Drop echo v4. |
| DI / interfaces | `IXService`/`IXDB` interfaces, aggregator structs, `RegisterRoutes(group)` | `TaskService` is a concrete struct. Repository and cache are free functions that take `*gorm.DB`/`*redis.Client`. Only the handler has an interface. | Add `ITaskRepository`, `ITaskCache`, `ITaskService`. Inject them via constructors. Add `internal/route` with an `api/v1` group. |
| Context propagation | `db.WithContext(ctx)` everywhere | Repository ignores ctx. `CreateTask`/`GetTasks` take no ctx. | Pass ctx through every layer. Use `WithContext`. |
| Errors / responses | Typed app error + one central writer + JSON envelope | `c.String` with plain-text errors, a `switch` repeated in every handler, mixed messages (`"database connection error"`) | Add `internal/utils/response` (timeproduct style) or an `apperror` + `HTTPErrorHandler` (prowifi style). |
| Validation | validator v10 + `c.Validate` | Manual `slices.Contains` in the service. `UpdateTask` doesn't validate status. The handler builds the update map (business logic in the handler). | Add `validate:"required,oneof=..."` tags and a CustomValidator. Move map building into the service. |
| Middleware | RequestID, RequestLogger, Recover, Timeout, Prometheus, OTel, configurable CORS | `CORS("*")` only | Add the timeproduct middleware stack (see Tier 1/2). |
| Health | `/v1/health` skipped by logs and auth; ECS uses it as its health check | `/health` checks the DB only. Redis isn't checked. Uptime is the string `"Placeholder"`. | Add a liveness endpoint and a readiness endpoint (DB + Redis). Compute real uptime. |
| Lifecycle | `signal.NotifyContext`, graceful shutdown, close DB, `ReadHeaderTimeout` | `e.Start` error ignored, `panic` on connect errors, no shutdown | Copy the main.go shape from timeproduct_be/cmd/api/main.go or prowifi-be/cmd/app/api.go. |
| Logging | slog JSON with request_id in context | slog default mixed with `fmt.Println` in service/cache | Replace every `fmt.Print*`. Add a context-attr slog handler (a small custom one, since the org module is private). |
| Data model | Audit fields, soft delete, indexes, typed IDs | `id`, `title`, `status` only | Add `created_at/updated_at/deleted_at`, `description`, `due_date`, a DB `CHECK` on status. |
| Migrations | goose + embed + `migrate` binary | `init.sql` mounted into the docker entrypoint. An empty untracked `sql/` dir exists. | Add goose migrations and a `cmd/migrate`. Remove init.sql. |
| List endpoint | page/page_size/sort allowlist/filter + metadata | `GetTasks` returns all rows | Add pagination, `?status=` filter, and sort (noblock-be/pkg/helper). |
| Caching | Prefixed keys, wrapper client, TTL from config | Cache-aside on GetTask only, hardcoded TTL and key prefix | Wrap in a `RedisClient` with prefix + TTL from config (timeproduct_be/internal/redis/redis.go). |
| Auth | API key or JWT + Redis, RBAC | None | Tier 2: API key middleware. Tier 3: JWT access/refresh + Casbin. |
| Async | Worker/scheduler binaries, queue with ack/nack/retry | None | Tier 3: Postgres-queue worker or cron job. |
| Outbound clients | `HTTPClient` interface + retry | None | Tier 3: a client for an external API with retry + `httptest` tests. |
| Testing | Unit tests with mockery mocks per layer, table-driven, `-race -cover` | `main_test.go` at module root hits a real Postgres/Redis with a hardcoded DSN. A shared ID assertion depends on TRUNCATE. | Move tests next to their code. Mock the service in handler tests and the repo/cache in service tests. Keep the DB tests as tagged integration tests. |
| Tooling | Makefile, .golangci.yml, .air.toml, .mockery.yml, pre-commit, AGENTS.md | None. Root README is 23 bytes. | Add all of them (Tier 1). |
| Docker | Multi-stage, mod-cache layer, alpine runtime, entrypoint runs migrations | Single stage `golang:1.27.1`, `COPY . .` before `go mod download`, no .dockerignore | Rewrite as multi-stage (tapsales-be/Dockerfile). Add .dockerignore. |
| Compose | Healthchecks, env defaults | `api` depends on `db` only (Redis isn't listed). No healthchecks. Credentials are inline. | Add `depends_on: condition: service_healthy`, redis, and `${VAR:-default}` (prowifi-be/docker-compose.yml). |
| CI | bitbucket-pipelines test → lint → build → push | None | The repo is local; add GitHub Actions or Bitbucket pipelines that mirror timeproduct_be's steps (excluding the ECR push). |

## (c) Roadmap

### Tier 1: quick wins (each about an hour or less)
1. **Hygiene + go.mod.** Remove the binaries and .DS_Store, add `.gitignore`, run `go mod tidy`. See (d).
2. **viper config.** Port `timeproduct_be/config/config.go` (Server/Db/Redis/Logging/Timeout) and add `config.example.yml`. `database.ConnectDB(cfg)` builds the DSN from config.
3. **Makefile + golangci-lint.** Copy the target set from `prowifi-be/Makefile` (run, dev, build, test `-race -cover`, lint, migrate-*, tidy). Copy the linter list from `timeproduct_be/.golangci.yml`. Fix whatever errcheck flags (e.g. ignored `e.Start`).
4. **Middleware basics.** RequestID (uuid), RequestLogger skipping health, Recover, Timeout. Imitate `timeproduct_be/internal/server/server.go` and `internal/middleware/timeout_middleware.go`.
5. **Response envelope.** Add `internal/utils/response` and the `response.Make` wrapper from `timeproduct_be/internal/utils/response/response.go`. Every handler returns JSON errors, and the duplicated `switch` blocks go away.
6. **Validator.** Register go-playground/validator as `e.Validator`, use `validate:"required,oneof=todo inprogress ..."` on the request entities, and return 422 with field errors.
7. **Logging cleanup.** Replace every `fmt.Println` in `services/tasks/service.go` and `data_access/cache.go` with `slog.*Context`. Make log level configurable (`timeproduct_be/internal/observability/log.go`).

### Tier 2: medium (half a day to a day each)
8. **Interfaces + DI refactor.** `ITaskRepository` (struct with `db`, every method takes ctx and uses `WithContext`), `ITaskCache`, `ITaskService`. Add `DataAccess`/`Services`/`Handlers` aggregators and `route.APIRoute` with an `api/v1` group. Imitate `timeproduct_be/internal/{route,handlers,services,data_access}/*.go`. Optionally use one file per operation, as in timeproduct_be.
9. **Mockery + unit tests.** Add `.mockery.yml` (copy the timeproduct_be format) and generate `internal/mocks`. Write a table-driven service test (repo/cache mocked, including cache hit/miss/Redis error) and a handler test (service mocked, asserting status + JSON). Follow `timeproduct_be/internal/services/promotions/promotions_test.go`, `handlers/promotions/promotions_test.go`, and `prowifi-be/.kiro/steering/testing.md`. Move `main_test.go` into a `//go:build integration` test.
10. **goose migrations.** Embed `00001_create_tasks.sql` and `00002_add_audit_fields.sql` (created_at/updated_at/deleted_at, description, due_date, status CHECK, index on status). Add a `cmd/migrate` and an entrypoint `./migrate && ./server`. Imitate `timeproduct_be/cmd/migrate/main.go`, `internal/database/migrations/`, and `docker-entrypoint.sh`.
11. **Pagination/filter/sort.** `GET /api/v1/tasks?page=&page_size=&status=&sort=&order=` returns `{data, metadata{page,page_size,total_items,total_pages}}`. Copy `noblock-be/pkg/helper/pagination.go` and `sorting.go` (keep the allowlist comment, which prevents SQL injection).
12. **Graceful shutdown + health/readiness.** Use `signal.NotifyContext`, `echo.StartConfig.Start(ctx)`, `ReadHeaderTimeout`, and close the DB and Redis on exit (`prowifi-be/cmd/app/api.go`). Split `/v1/health` (liveness) from `/v1/ready` (DB + Redis ping with a timeout) and compute real uptime.
13. **Prometheus.** Add the `echoprometheus` middleware and a separate metrics server on its own port (`timeproduct_be/internal/metrics/server.go`). Add a Prometheus container to compose to scrape it.
14. **Central error type.** `internal/apperror` with Status/Code/Message + `Wrap`/`Unwrap`, plus one `HTTPErrorHandler` that includes `request_id` in the body. Map Postgres 23505 to a 409 (e.g. a unique task title per user). Imitate `prowifi-be/internal/apperror/apperror.go`, `internal/server/errors.go`, and `internal/repository/order.go`.
15. **Docker/compose.** Multi-stage Dockerfile (`tapsales-be/Dockerfile`), `.dockerignore`, compose healthchecks + redis dependency + env defaults (`prowifi-be/docker-compose.yml`), and `.air.toml` for hot reload.
16. **API key auth.** An `api_keys` table plus middleware with an in-memory TTL cache that skips health (`timeproduct_be/internal/middleware/auth_middleware.go`).

### Tier 3: advanced (multi-day)
17. **Users + JWT.** Register/login with bcrypt, HS256 access/refresh tokens with hashed refresh tokens stored in Redis, rotation, and revocation. Tasks become owned by users (`user_id` FK, and every query filters by owner). Imitate `tapsales-be/pkg/jwt/jwt.go` and `internal/middleware/access_token.go`.
18. **RBAC with Casbin.** Admin/member roles. `methodToAction` + `Enforce(domain)` middleware, with policies embedded from CSV. Imitate `noblock-be/internal/middlewares/casbin_middleware.go`, `casbin_model.conf`, and `internal/casbin/`.
19. **Background worker + scheduler.** Make the binary a single `app api|migrate|worker` CLI (`prowifi-be/cmd/app/main.go`). Add a Postgres-backed queue (e.g. "task due reminder" jobs) with ack/nack/retry/visibility timeout (`prowifi-be/internal/worker/worker.go`), and a cron job that marks overdue tasks (`timevoucher-be/internal/scheduler/scheduler.go`, robfig/cron).
20. **Outbound client + retry.** Call an external API (e.g. a webhook or notification service) through an `HTTPClient` interface with functional options, request-ID forwarding, and exponential-backoff retry. Test it with `httptest.NewServer`. Imitate `tapsales-be/pkg/httpclient/httpclient.go`, `internal/utils/retry/retry.go`, and `prowifi-be/internal/httpx/`.
21. **OpenTelemetry tracing.** Use the public `go.opentelemetry.io/otel` SDK + otelecho (the org's time-trace module is private). Put trace_id/span_id into slog context. Add Jaeger or an OTel collector to compose.
22. **Read replica + pool config.** Pool settings with `cmp.Or` defaults, plus dbresolver (`timeproduct_be/internal/database/postgres.go`, config `DbDetail`).
23. **OpenAPI.** Either a spec-first design with oapi-codegen echo strict server (`selfcare-be/server.cfg.yaml`, `api/*.gen.go`), or a hand-written `openapi.yaml` that replaces the Postman collection.
24. **CI pipeline.** test → `go vet` + golangci-lint → docker build, mirroring `timeproduct_be/bitbucket-pipelines.yml` (GitHub Actions is fine for a personal repo). Optionally add a pre-commit hook (`timevoucher-be/scripts/pre-commit`).
25. **Infra (stretch).** Read `terraform/svc/timeproduct/base/main.tf` to see how a service maps to ECS + ALB + health check + SSM config. Locally, an equivalent exercise is a k8s/kind manifest with liveness/readiness probes on the endpoints from item 12. Do not run anything against the real terraform repo.

### Frontend (brief)
basecamp-fe uses Vite + TanStack Router + TanStack Query (`src/hooks/use-*.ts`, `src/lib/query/keys.ts`), a shared `apiFetch` with an `ApiError` class and 401 token-renew retry (`src/lib/api/client.ts`), react-hook-form + zod schemas (`src/schemas/`), radix/shadcn UI (`components.json`, `src/components/ui`), OIDC auth (`src/lib/auth/`), vitest + Testing Library, oxlint + prettier, and `.env.local.example`.

The learning frontend hardcodes `API_URL` (in two files), uses raw `fetch` without error handling, manages forms with `useState`, and has no tests.

Suggested exercises:
- Read the API base URL from an env var.
- Write a single typed API client that understands the new response envelope.
- Use TanStack Query for task list and mutations (Next.js is fine; you don't need to switch to Vite).
- Use zod + react-hook-form in createTask.
- Add vitest tests for `lib/tasks.ts`.
- Add pagination UI once item 11 lands.

Per frontend/AGENTS.md, read `node_modules/next/dist/docs/` before changing Next.js code.

## (d) Hygiene issues in the learning project
- `task_api/task_api` (~16 MB binary) is **committed**. `task_api/main` (~21 MB) is untracked. `task_api/cmd/api/main` is listed in `git ls-files`. Remove them with `git rm --cached` and gitignore `/main`, `/task_api`, `bin/`, `tmp/`.
- `.DS_Store` and `task_api/.DS_Store` are committed. Gitignore them.
- There's no `.gitignore` at the root or in task_api (only `frontend/.gitignore`).
- DB credentials are hardcoded in `internal/database/database.go`, `main_test.go`, and `compose.yaml`. They're dev-only, but they should move to config/env.
- `go.mod` declares `go 1.27.1` and marks everything indirect. Run tidy.
- `config/config.go` is an empty package. `sql/` is an empty untracked dir.
- There are two Postman folders (`.postman/` and `postman/`). Consolidate them, or replace them with OpenAPI.
- Root `README.md` is a stub. Add setup/run/test/migrate instructions in the style of timeproduct_be/README.md.
- Code smells to fix along the way:
  - `DeleteTask` lacks the `id <= 0` check.
  - `ErrRedisError` is unused.
  - `parseTaskID` returns an `echo.HTTPError` (a different shape from the other errors).
  - `NewTaskService` returns a value while the handler stores it by value, with pointer receivers on the service.
  - Model formatting isn't gofmt'd (`internal/models/task.go`).
- There are 13 uncommitted modified files plus 3 new ones (cache.go, errors.go, types.go). Commit the current work before starting the roadmap.

## What I could not verify
- I did not build or run task_api or its tests (read-only scope). The build status with the current go.mod is unknown.
- selfcare-be and backoffice-v2-be were sampled only through go.mod, layout, and grep. timevoucher-be's internals and test coverage weren't examined beyond the scheduler/worker/Makefile/pre-commit. basecamp-fe was reviewed from package.json, layout, and the start of `src/lib/api/client.ts`.
- I did not check test coverage numbers or whether each repo's CI passes.
- The private org modules (gotime observability log/trace, time-trace, pg/queue) weren't inspected. Their behavior is inferred from call sites.
- I didn't check whether the AutoMigrate repos also have external migration tooling.
- No secrets, hostnames, or tokens from the reference repos are reproduced here. Config keys are referenced by name only.
