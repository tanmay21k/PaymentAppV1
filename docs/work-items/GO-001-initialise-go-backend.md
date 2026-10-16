# GO-001: Initialise the Go backend

- **Status:** Awaiting Approval
- **Created:** 2026-10-17
- **Approved by owner:** no
- **Depends on:** none

## 1. Title
Initialise `go-backend/` as a runnable Go service with configuration, a PostgreSQL connection, and a health endpoint.

## 2. Problem statement
`go-backend/` contains only `AI-WORKFLOW.md` [Confirmed: directory listing]. The previous Go code, including a PostgreSQL user repository, sqlc setup and a JWT helper, was deleted in commit `e1e207f` as part of a reinitialisation [Confirmed: `git show e1e207f`]. There is no Go module, entrypoint, configuration loading, or database connection to build the migration on.

This item creates the smallest working foundation: a Go module that builds, starts an HTTP server, reads configuration from environment variables, connects to PostgreSQL, and exposes `GET /health`. No API endpoints from `openapi.yaml` are implemented here.

## 3. Business and technical context
- The migration target is Go with PostgreSQL, replacing the JS backend's MongoDB [Confirmed: `CLAUDE.md`, "Project"].
- The owner is learning Go by writing this code personally. Each choice below should be small enough to understand and explain [Confirmed: `CLAUDE.md`].
- Decisions confirmed by the owner in discussion:
  - Scope: skeleton plus PostgreSQL connection and health endpoint; no migrations [Confirmed: owner].
  - Router: standard library `net/http` `ServeMux` [Confirmed: owner].
  - Database driver: `pgx/v5` [Confirmed: owner].
  - Migration tool (goose): chosen for a later work item, not GO-001 [Confirmed: owner; scope is the skeleton, so GO-001 adds no migration tool].
- Existing environment conventions that the Go service should follow:
  - API port 4500 [Confirmed: `.env.example:1`, `openapi.yaml:8-9`].
  - Configuration via environment variables loaded from `backend/.env`-style files [Confirmed: `.env.example:1-3`, `README.md` "Configure backend environment variables"].

## 4. Scope
- Go module at `go-backend/go.mod`, with the Go toolchain directive set to the installed toolchain (go1.27.2) [Assumption: module name and Go directive; see section 14].
- Entrypoint `go-backend/cmd/api/main.go`: loads config, opens the database pool, starts the HTTP server, and shuts down gracefully on SIGINT/SIGTERM.
- Configuration package `go-backend/internal/config`: reads `PORT` and `DATABASE_URL` from the environment and validates them.
- Database package `go-backend/internal/database`: creates a `pgxpool` connection pool from `DATABASE_URL` and exposes a ping check [Assumption: package names].
- HTTP health handler `GET /health`, registered on `net/http` `ServeMux`. It reports whether the database answers a ping.
- `go-backend/.env.example` with `PORT` and `DATABASE_URL`; `go-backend/.gitignore` ignoring `.env`.
- `go-backend/docker-compose.yml` containing only a `postgres` service with a health check and a named volume. It is separate from the root `docker-compose.yml`, which runs the JS stack and is not changed by this item [Assumption: separate file].
- Unit tests for configuration and for the health handler.

## 5. Non-goals
- Any endpoint from `openapi.yaml` (signup, signin, transfers, and so on).
- Database migrations, schema, or any table (including `users`) [Confirmed: owner's scope choice].
- Migrating the root `docker-compose.yml`, the frontend, or `backend/` (JS). Nothing under `backend/` or `frontend/` changes.
- CORS configuration for the frontend origin `http://localhost:5173` [Assumption: deferred to the first endpoint item; see section 14].
- Authentication, JWT handling, logging frameworks, or a dependency-injection library.
- A Go Dockerfile for the backend image.

## 6. Relevant OpenAPI operations and schemas
- None. `GET /health` is an operational endpoint, not in `openapi.yaml` [Confirmed: no `/health` path in `openapi.yaml`, paths listed at lines 15-328].
- Server URL used for the default port: `http://localhost:4500` (`openapi.yaml:8-9`) [Confirmed].

## 7. Relevant JavaScript implementation references
- `backend/server.js:10` — `PORT` defaults to `3000` when unset. The Go default is a separate decision; see section 8 and section 14.
- `backend/server.js:13-15` — CORS allows only `http://localhost:5173`. Not in GO-001 scope.
- `backend/config/connectDb.js` — MongoDB connection. Replaced by the PostgreSQL connection in GO-001; not read in detail for this item.
- `.env.example:1-3` — `PORT=4500`, `MONGO_URI`, `JWT_SECRET`. The Go service uses `DATABASE_URL` instead of `MONGO_URI` [Assumption: name; see section 14].

## 8. Expected behavior
- **[Confirmed]** The service listens on the port given by `PORT`, which defaults to `4500` when unset [Assumption on default: `.env.example:1` and `openapi.yaml:8-9` both use 4500; the JS default of 3000 is not carried over].
- **[Confirmed]** `DATABASE_URL` is required. The process refuses to start without it.
- **[Assumption]** The server starts even when the database is unreachable at startup, so that `/health` can report the problem. The pool is created without blocking on a successful connection.
- **[Assumption]** `GET /health` returns JSON with a `status` field and a `database` field. Success body: `{"status":"ok","database":"ok"}`. Database failure body: `{"status":"unavailable","database":"unreachable"}`.
- **[Assumption]** Responses from `/health` use `Content-Type: application/json`.
- **[Confirmed]** Shutdown on SIGINT or SIGTERM stops accepting new requests, waits up to 10 seconds for in-flight requests, closes the database pool, and exits with status 0 [Assumption on the 10-second value].

## 9. Proposed design and architectural considerations
**Layout (proposed):**
```text
go-backend/
├── AI-WORKFLOW.md          (existing)
├── .env.example
├── .gitignore
├── docker-compose.yml      (postgres service only)
├── go.mod
├── go.sum
├── cmd/api/main.go         (wires config, pool, server; owns shutdown)
└── internal/
    ├── config/config.go    (+ config_test.go)
    ├── database/db.go      (pgxpool creation, Ping)
    └── health/health.go    (+ health_test.go)
```

**Design points:**
- The health handler depends on a small interface with one method, `Ping(ctx) error`, not on `*pgxpool.Pool` directly. Tests pass a fake implementation. This is the first use of an interface in the project.
- Configuration is read once in `main` and passed down as values. No global variables.
- Routing uses `mux.HandleFunc("GET /health", ...)`. Go 1.22+ method patterns make a `GET`-only route, and `405` is returned for other methods automatically [Assumption: `405` behaviour, verify in tests].

**Alternatives considered:**
- `chi` or `gin` routers: rejected for GO-001 because the owner chose the standard library; they can be reconsidered when routing needs grow.
- `database/sql` with `lib/pq`: rejected by the owner in favour of `pgx/v5`.
- Starting in `main` with `log.Fatal` on any error: simpler, but makes `/health` impossible to reach while the database is down. Rejected under the assumption in section 8.

**Affected files:** only `go-backend/`. Nothing else in the repository changes.

## 10. Acceptance criteria
- [ ] **AC1:** `go build ./...` run in `go-backend/` exits 0. Verify: command output.
- [ ] **AC2:** `go vet ./...` exits 0. Verify: command output.
- [ ] **AC3:** `go test ./...` exits 0 and runs the config and health tests. Verify: command output with `-v` showing the test names.
- [ ] **AC4:** With `DATABASE_URL` unset, `go run ./cmd/api` exits with a non-zero status and an error message containing `DATABASE_URL`, before listening on any port. Verify: command output and `echo $?`.
- [ ] **AC5:** With `PORT` unset, the service listens on 4500. Verify: startup log line shows `4500`, and `curl -i localhost:4500/health` returns a response.
- [ ] **AC6:** With Postgres running (`docker compose -f go-backend/docker-compose.yml up -d postgres`) and `DATABASE_URL` set, `curl -i localhost:4500/health` returns `200` and body `{"status":"ok","database":"ok"}`. Verify: curl output.
- [ ] **AC7:** With Postgres stopped and the service running, `curl -i localhost:4500/health` returns `503` and body `{"status":"unavailable","database":"unreachable"}`. Verify: curl output.
- [ ] **AC8:** `curl -i -X POST localhost:4500/health` returns `405`. Verify: curl output.
- [ ] **AC9:** `curl -i localhost:4500/nope` returns `404`. Verify: curl output.
- [ ] **AC10:** Pressing Ctrl+C on a running service logs a shutdown message and the process exits with status 0. Verify: terminal output and `echo $?`.
- [ ] **AC11:** `docker compose -f go-backend/docker-compose.yml up -d postgres` followed by `docker compose -f go-backend/docker-compose.yml ps` shows the `postgres` service as `healthy`. Verify: command output.
- [ ] **AC12:** `go-backend/.env.example` lists `PORT` and `DATABASE_URL` with no real secret values, and `go-backend/.gitignore` contains a line that ignores `.env`. Verify: file contents, and `git check-ignore go-backend/.env` prints the path.
- [ ] **AC13:** No file outside `go-backend/` and `docs/work-items/` is changed by this item. Verify: `git status --short` shows only those paths.

## 11. Edge cases and error handling
| Case | Expected status/response | Source |
| --- | --- | --- |
| `DATABASE_URL` unset | Process exits non-zero, error names `DATABASE_URL` | Owner decision + AC4 |
| `DATABASE_URL` malformed (e.g. `not a url`) | Process exits non-zero with the parse error | [Assumption] pgx `ParseConfig` error at startup |
| `PORT` non-numeric (e.g. `abc`) | Process exits non-zero, error names `PORT` | [Assumption] |
| `PORT` out of range (e.g. `70000`) | Process exits non-zero, error names `PORT` | [Assumption] |
| Database down at startup | Server starts; `/health` returns 503 | Section 8 assumption, AC7 |
| Database goes down while running | `/health` returns 503 | AC7 |
| `GET /health` with a hung database | `/health` responds within 2 seconds with 503 (context timeout) | [Assumption] 2-second timeout |
| `POST /health` | 405 | AC8, Go 1.22 ServeMux behaviour |
| Unknown path | 404 | AC9 |
| SIGTERM during an in-flight request | Request completes or is cut off after 10 seconds; process exits 0 | Section 8 assumption |

## 12. Testing requirements
- **Unit (config):** table-driven tests for: both variables set; `PORT` unset (defaults to 4500); `DATABASE_URL` unset (error); `PORT` non-numeric (error); `PORT` out of range (error). Covers AC3, AC4, AC5 and the edge cases in section 11.
- **Unit (health):** `httptest` tests with a fake pinger for: ping succeeds (200 and body); ping fails (503 and body); ping slower than the timeout (503). Covers AC3, AC6, AC7 at the handler level.
- **Manual / integration:** AC4 to AC11 run by the owner with the commands in section 10. Capture the output as evidence.
- **Commands:** `go build ./...`, `go vet ./...`, `go test ./... -v`, run from `go-backend/`.

## 13. Relevant Go concepts to learn
- **Modules and packages** — `go mod init`, the `go` directive, `cmd/` and `internal/` layout, why `internal/` blocks imports from outside the module.
- **Environment configuration** — `os.LookupEnv` vs `os.Getenv`, strconv parsing, returning errors instead of exiting from inside a library function.
- **Errors** — returning errors up the call stack and wrapping them with `%w` (used when the database or config error is reported in `main`).
- **Interfaces** — implicit satisfaction; defining a one-method interface where it is used (the health handler's `Pinger`).
- **`net/http` ServeMux** — method and path patterns (`GET /health`), what happens for unmatched methods and paths.
- **`context` and timeouts** — `context.WithTimeout` for the database ping inside the handler.
- **Goroutines and graceful shutdown** — `http.Server` running in a goroutine, `signal.NotifyContext`, `Server.Shutdown` with a deadline.
- **Testing** — table-driven tests, `t.Run`, `httptest.NewRecorder` and `httptest.NewRequest`.
- **pgx** — `pgxpool.New` and `ParseConfig`; `Pool.Ping` and `Pool.Close`.

## 14. Dependencies and risks
**Assumptions (not yet confirmed by the owner):**
- Module name `payflow` (`go.mod`). Alternative: a full repository path.
- Go directive matches the installed toolchain (go1.27.2). This may matter if the owner later builds in Docker with an older image.
- Configuration variable names: `PORT` and `DATABASE_URL` (instead of `MONGO_URI`).
- Default port 4500 when `PORT` is unset (the JS backend defaults to 3000).
- Startup succeeds without a reachable database, and `/health` reports the state.
- `/health` response bodies and the 2-second ping timeout.
- The `postgres` image tag: `postgres:17`.
- The Compose file lives at `go-backend/docker-compose.yml`, separate from the root file.
- Shutdown timeout of 10 seconds.
- `CORS` is deferred to the first endpoint work item.

**Open questions for the owner:**
1. Confirm the module name, or give the path you want.
2. Confirm the environment variable names and the default port.
3. Confirm whether the service should start without a database (section 8) or exit on a failed startup ping.
4. Confirm the `/health` response shape.
5. Confirm the separate `go-backend/docker-compose.yml` is acceptable, rather than adding Postgres to the root Compose file.

**Risks:**
- Dependencies `pgx/v5` and its transitive modules are fetched at `go mod tidy` time; no network means GO-001 cannot be built.
- Docker is installed here (`/usr/bin/docker`), so AC11 can be run locally [Confirmed: `which docker`]. Whether the Docker daemon is running is not yet checked.

## 15. Definition of Done
- [ ] All acceptance criteria verified with evidence
- [ ] Tests pass (output observed)
- [ ] OpenAPI contract respected; no contract endpoints implemented in this item, so no discrepancies apply
- [ ] Review outcome PASS or PASS WITH NON-BLOCKING OBSERVATIONS; findings addressed or explicitly accepted
- [ ] Documentation updated where needed (`go-backend/AI-WORKFLOW.md` is unchanged; a `go-backend/README.md` is out of scope unless the owner asks)
- [ ] Owner can explain key Go concepts and design decisions
- [ ] Owner authorized closure

---
## Closure record (filled at closure)
- **Verification evidence:**
- **Design decisions:**
- **Go concepts learned:**
- **Remaining technical debt:**
- **Suggested next task (not started):**
