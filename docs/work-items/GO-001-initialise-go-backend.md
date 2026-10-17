# GO-001: Initialise the Go backend

- **Status:** CLOSED
- **Created:** 2026-10-17
- **Approved by owner:** yes
- **Depends on:** none
- **Revised:** 2026-10-17 — added loading of `go-backend/.env` via `github.com/joho/godotenv` (sections 4, 9, 10, 11, 13, 14)
- **Revised:** 2026-10-17 — owner decision: the process exits when a database health check fails (sections 8, 9, 10 AC7, 11, 14, open question 3)
- **Revised:** 2026-10-17 — owner decision: no Docker Compose healthcheck in GO-001; AC11 removed and the healthcheck and restart policy moved to GO-002 (sections 4, 5, 10, 12, 14)
- **Revised:** 2026-10-17 — owner decision: fake ping deprecated; no `Pinger` interface; health tests use a real pool (sections 9, 12, 13)
- **Revised:** 2026-10-17 — added `TestHealthReturns503WhenDatabaseIsSlow` (section 12); AC6 verified (section 10); `ConnectDB` no longer takes a context (no doc section describes its signature)
- **Revised:** 2026-10-17 — closed on owner authorization; status set to CLOSED; learning-log entry appended

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
- `.env` loading in `cmd/api/main.go`: before config is read, call `godotenv.Load()` from `github.com/joho/godotenv` to read `go-backend/.env` into the process environment [Confirmed: owner, "read from .env in main.go"]. A missing `.env` file is not an error, so the service still runs when variables come from the shell or the environment [Assumption].
- New dependency `github.com/joho/godotenv`, added to `go.mod` and `go.sum`. Reason: the standard library has no `.env` parser [Confirmed: owner agreed to the dependency in discussion; reason per `CLAUDE.md` rule 9].
- Database package `go-backend/internal/database`: creates a `pgxpool` connection pool from `DATABASE_URL` and exposes a ping check [Assumption: package names].
- HTTP health handler `GET /health`, registered on `net/http` `ServeMux`. It reports whether the database answers a ping.
- `go-backend/.env.example` with `PORT` and `DATABASE_URL`; `go-backend/.gitignore` ignoring `.env`.
- `go-backend/docker-compose.yml` containing only a `postgres` service and a named volume. No Docker healthcheck or restart policy (moved to GO-002). It is separate from the root `docker-compose.yml`, which runs the JS stack and is not changed by this item [Assumption: separate file].
- Unit tests for configuration and for the health handler.

## 5. Non-goals
- Any endpoint from `openapi.yaml` (signup, signin, transfers, and so on).
- Database migrations, schema, or any table (including `users`) [Confirmed: owner's scope choice].
- Migrating the root `docker-compose.yml`, the frontend, or `backend/` (JS). Nothing under `backend/` or `frontend/` changes.
- CORS configuration for the frontend origin `http://localhost:5173` [Assumption: deferred to the first endpoint item; see section 14].
- Authentication, JWT handling, logging frameworks, or a dependency-injection library.
- A Go Dockerfile for the backend image.
- Docker Compose healthcheck for `postgres` and restart policy for the Go service. Moved to GO-002 (owner decision, 2026-10-17).

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
- **[Confirmed: owner decision, 2026-10-17]** The process exits when a database health check fails. The check runs in the background 30 seconds after startup and then every 5 minutes. A failed check logs `[ERROR]: database health check failed` and exits with status 1. If the database is down at startup, the process exits about 30 seconds after starting. If it goes down while running, the process exits at the next scheduled check, at most 5 minutes later. While the process is alive, `/health` reports the database state on each request. The pool is created without blocking on a successful connection, so startup itself does not fail on an unreachable database. GO-001 adds no restart policy for the Go service.
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
- The health handler depends on `*pgxpool.Pool` directly. No `Ping(ctx) error` interface and no fake pinger (owner decision, 2026-10-17: fake ping deprecated). Handler tests use a real `pgxpool.Pool` pointed at an unreachable address.
- Configuration is read once in `main` and passed down as values. No global variables.
- Order in `main`: `godotenv.Load()` first, then `config` reads the environment. `Load()` reads `.env` relative to the current working directory, so the binary is run from `go-backend/` (`go run ./cmd/api` from that directory) [Confirmed: godotenv behaviour; verify with `go doc github.com/joho/godotenv` once the module is fetched].
- Precedence: `godotenv.Load()` does not overwrite variables that are already set in the shell. A `PORT` exported in the shell beats `PORT` in `.env` [Assumption: verify against godotenv docs; AC15 tests it]. `godotenv.Overload()` would reverse this, and is not used here.
- A missing `.env` is detected with `errors.Is(err, fs.ErrNotExist)` and ignored. Any other error from `Load()` (for example a malformed line) stops startup [Assumption].
- Routing uses `mux.HandleFunc("GET /health", ...)`. Go 1.22+ method patterns make a `GET`-only route, and `405` is returned for other methods automatically [Assumption: `405` behaviour, verify in tests].

**Alternatives considered:**
- `chi` or `gin` routers: rejected for GO-001 because the owner chose the standard library; they can be reconsidered when routing needs grow.
- `database/sql` with `lib/pq`: rejected by the owner in favour of `pgx/v5`.
- Failing at startup on the first database ping: not chosen. The owner decided the first check runs 30 seconds after startup and that a failed check exits the process (section 8).

**Affected files:** only `go-backend/`. Nothing else in the repository changes.

## 10. Acceptance criteria
- [x] **AC1:** `go build ./...` run in `go-backend/` exits 0. Verify: command output.
- [x] **AC2:** `go vet ./...` exits 0. Verify: command output.
- [x] **AC3:** `go test ./...` exits 0 and runs the config and health tests. Verify: command output with `-v` showing the test names.
- [x] **AC4:** With `DATABASE_URL` unset, `go run ./cmd/api` exits with a non-zero status and an error message containing `DATABASE_URL`, before listening on any port. Verify: command output and `echo $?`.
- [x] **AC5:** With `PORT` unset, the service listens on 4500. Verify: startup log line shows `4500`, and `curl -i localhost:4500/health` returns a response.
- [x] **AC6:** With Postgres running (`docker compose -f go-backend/docker-compose.yml up -d postgres`) and `DATABASE_URL` set, `curl -i localhost:4500/health` returns `200` and body `{"status":"ok","database":"ok"}`. Verify: curl output. **Evidence (2026-10-17):** `postgres-db` container (image `postgres:16`, compose project `go-backend`) accepted connections (`pg_isready`). The built service, run with `DATABASE_URL=postgres://myuser:…@localhost:5432/mydb?sslmode=disable`, logged `listening on port 4500`. `curl -i localhost:4500/health` returned `HTTP/1.1 200 OK`, `Content-Type: application/json`, body `{"status":"ok","database":"ok"}`.
- [ ] **AC7:** With Postgres stopped and the service running, `curl -i localhost:4500/health` returns `503` and body `{"status":"unavailable","database":"unreachable"}` while the process is alive. The process then exits with status 1 and logs `database health check failed` at the next scheduled check: about 30 seconds after startup if Postgres is already down, otherwise at most 5 minutes after Postgres stops. Verify: curl output, the startup and exit log lines, and `echo $?`.
- [x] **AC8:** `curl -i -X POST localhost:4500/health` returns `405`. Verify: curl output.
- [x] **AC9:** `curl -i localhost:4500/nope` returns `404`. Verify: curl output.
- [x] **AC10:** Pressing Ctrl+C on a running service logs a shutdown message and the process exits with status 0. Verify: terminal output and `echo $?`.
- ~~**AC11:** Compose `postgres` service shows `healthy`.~~ Removed 2026-10-17 (owner decision): the Docker healthcheck moved to GO-002. The number is kept so other AC references stay stable.
- [x] **AC12:** `go-backend/.env.example` lists `PORT` and `DATABASE_URL` with no real secret values, and `go-backend/.gitignore` contains a line that ignores `.env`. Verify: file contents, and `git check-ignore go-backend/.env` prints the path.
- [x] **AC13:** No file outside `go-backend/` and `docs/work-items/` is changed by this item. Verify: `git status --short` shows only those paths.
- [x] **AC14:** With `go-backend/.env` containing `PORT=4600` and `PORT` not set in the shell, `go run ./cmd/api` logs port `4600` and `curl -i localhost:4600/health` returns a response. With no `.env` file present and `PORT`/`DATABASE_URL` exported in the shell, the service starts normally. Verify: startup log and curl output for both runs.
- [x] **AC15:** With `go-backend/.env` containing `PORT=4600` and `PORT=4700` exported in the shell, the service logs port `4700`. Verify: startup log output.
- [x] **AC16:** `go.mod` lists `github.com/joho/godotenv` as a direct requirement, and `go build ./...` succeeds with `go.sum` committed alongside it. Verify: `grep godotenv go-backend/go.mod` and `go build ./...` output.

## 11. Edge cases and error handling
| Case | Expected status/response | Source |
| --- | --- | --- |
| `DATABASE_URL` unset | Process exits non-zero, error names `DATABASE_URL` | Owner decision + AC4 |
| `DATABASE_URL` malformed (e.g. `not a url`) | Process exits non-zero with the parse error | [Assumption] pgx `ParseConfig` error at startup |
| `PORT` non-numeric (e.g. `abc`) | Process exits non-zero, error names `PORT` | [Assumption] |
| `PORT` out of range (e.g. `70000`) | Process exits non-zero, error names `PORT` | [Assumption] |
| Database down at startup | Process starts and `/health` returns 503; the first check fails about 30 seconds after startup, and the process exits with status 1 | Section 8 (owner decision), AC7 |
| Database goes down while running | `/health` returns 503 until the next scheduled check; the process then exits with status 1 (at most 5 minutes later) | Section 8 (owner decision), AC7 |
| `GET /health` with a hung database | `/health` responds within 2 seconds with 503 (context timeout) | [Assumption] 2-second timeout |
| `POST /health` | 405 | AC8, Go 1.22 ServeMux behaviour |
| Unknown path | 404 | AC9 |
| `.env` file missing | Not an error; service uses shell/environment variables | Section 9 assumption, AC14 |
| `.env` file malformed (e.g. a line with no `=`) | Process exits non-zero with the godotenv error | [Assumption] godotenv error handling, section 9 |
| Variable set in both shell and `.env` | Shell value wins | Section 9 assumption, AC15 |
| SIGTERM during an in-flight request | Request completes or is cut off after 10 seconds; process exits 0 | Section 8 assumption |

## 12. Testing requirements
- **Unit (config):** table-driven tests for: both variables set; `PORT` unset (defaults to 4500); `DATABASE_URL` unset (error); `PORT` non-numeric (error); `PORT` out of range (error). Covers AC3, AC4, AC5 and the edge cases in section 11.
- **Unit (health):** `httptest` tests with a real `pgxpool.Pool` (no fake), in `internal/health/handler_test.go`: `TestHealthReturns503WhenDatabaseUnreachable` (ping fails, unreachable address: 503 and body) and `TestHealthReturns503WhenDatabaseIsSlow` (ping slower than the 2-second timeout, using a local TCP listener that accepts connections but never replies: 503 within `pingTimeout + 1s` and the 503 body). Covers AC3 and AC7 at the handler level. Passing output: `go test ./internal/health/ -v -count=1`, 2.00 seconds for the slow case. The 200 case needs a live Postgres, so it is covered by AC6 (manual, section 12 Manual/integration), not by a unit test.
- **Manual / integration:** AC4 to AC10 run by the owner with the commands in section 10. Capture the output as evidence.
- **Commands:** `go build ./...`, `go vet ./...`, `go test ./... -v`, run from `go-backend/`.

## 13. Relevant Go concepts to learn
- **Modules and packages** — `go mod init`, the `go` directive, `cmd/` and `internal/` layout, why `internal/` blocks imports from outside the module.
- **Environment configuration** — `os.LookupEnv` vs `os.Getenv`, strconv parsing, returning errors instead of exiting from inside a library function.
- **Loading `.env` files** — what `godotenv.Load` does to the process environment (`os.Setenv`), why it does not override existing variables, the difference from `Overload`, and `errors.Is` with `fs.ErrNotExist` for an optional file.
- **Third-party modules** — `go get`, `require` lines in `go.mod`, what `go.sum` records, and `go mod tidy`.
- **Errors** — returning errors up the call stack and wrapping them with `%w` (used when the database or config error is reported in `main`).
- **Interfaces** — implicit satisfaction. Deprecated for GO-001 (no `Pinger` interface); revisit when a second implementation is needed.
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
- The process starts without a reachable database but exits when a scheduled health check fails (owner decision, section 8). It does not stay up to report the state.
- `/health` response bodies and the 2-second ping timeout.
- The `postgres` image tag: `postgres:17`.
- The Compose file lives at `go-backend/docker-compose.yml`, separate from the root file.
- Shutdown timeout of 10 seconds.
- `CORS` is deferred to the first endpoint work item.
- `.env` loading uses `godotenv.Load()` in `main`; a missing file is ignored; shell variables take precedence (section 9).

**Dependencies (new):**
- `github.com/joho/godotenv` — loads `.env` files into the process environment. Agreed by the owner for inclusion in GO-001. Its version is pinned by `go.mod` at implementation time; the owner confirms the exact version when approving.

**Open questions for the owner:**
1. Confirm the module name, or give the path you want.
2. Confirm the environment variable names and the default port.
3. Resolved 2026-10-17: the process exits when a database health check fails (section 8).
4. Confirm the `/health` response shape.
5. Confirm the separate `go-backend/docker-compose.yml` is acceptable, rather than adding Postgres to the root Compose file.
6. Confirm the `.env` behaviour in section 9: missing file ignored, shell values win.

**Risks:**
- Dependencies `pgx/v5`, `godotenv`, and their transitive modules are fetched at `go mod tidy` time; no network means GO-001 cannot be built.
- `godotenv` is an extra dependency for a local-only convenience. Docker Compose can inject variables without it, so it is needed only for running the binary directly.
- Docker is installed here (`/usr/bin/docker`), and AC6 ran locally on 2026-10-17 against the `postgres:16` container (section 10, AC6 evidence). The daemon responded to `docker info` (server version 29.2.0). The `postgres-db` container was left running after the check.

## 15. Definition of Done
- [x] All acceptance criteria verified with evidence (AC7: the 5-minute exit path was verified by reading only; the owner accepted this at closure, see closure record)
- [x] Tests pass (output observed)
- [x] OpenAPI contract respected; no contract endpoints implemented in this item, so no discrepancies apply
- [x] Review outcome PASS WITH NON-BLOCKING OBSERVATIONS (final pass 2026-10-17); observations recorded in the closure record
- [x] Documentation updated where needed (`go-backend/AI-WORKFLOW.md` is unchanged; a `go-backend/README.md` is out of scope unless the owner asks)
- [x] Owner can explain key Go concepts and design decisions (learning check 2026-10-17; items still needing practice are listed in the closure record)
- [x] Owner authorized closure (2026-10-17)

---
## Closure record (filled at closure)
Drafted 2026-10-17 after the review and learning check. Closed 2026-10-17 on the owner's authorization; status set to `CLOSED`.

- **Verification evidence:**
  - Build, vet, format: `go build ./...` and `go vet ./...` exit 0; `gofmt -l .` prints nothing (final run, 2026-10-17).
  - Tests: `go test ./... -count=1 -v` passes. Config: `TestLoad` with 8 subtests. Health: `TestHealthReturns503WhenDatabaseUnreachable` and `TestHealthReturns503WhenDatabaseIsSlow` (2.00 s against a silent TCP listener).
  - AC4: with `DATABASE_URL` unset, `invalid configuration: DATABASE_URL is not set`, exit 1, no listening line.
  - AC5, AC8, AC9: on port 4500 with the database up, `/health` returned 200, `POST /health` returned 405, and `/nope` returned 404.
  - AC6: `postgres:16` container `postgres-db` (compose project `go-backend`) accepted connections. `curl -i localhost:4500/health` returned `HTTP/1.1 200 OK` with `{"status":"ok","database":"ok"}`.
  - AC10: SIGINT logged `shutdown complete` and exited 0.
  - AC14, AC15: `.env` with `PORT=4600` and shell `PORT` unset listened on 4600; no `.env` with exported variables listened on 4600; shell `PORT=4700` overrode `.env` `PORT=4600` and listened on 4700.
  - AC16: `go.mod` lists `github.com/joho/godotenv v1.5.1` as a direct requirement.
  - AC7 (partly verified): with the database down, `/health` returned 503 while the process was alive, and the process exited with status 1 at the first check (about 28 s after start, observed on an earlier build). The 5-minute path is verified by reading only.
  - Malformed `DATABASE_URL` exits 1 with the parse error. Hung database: 503 in 2.00 s.
  - Review: PASS WITH NON-BLOCKING OBSERVATIONS (final pass, 2026-10-17). Blocking findings from earlier passes are resolved: the slow-ping test is added, AC6 has evidence, `log.Fatal` is removed from the database package, and the handler logs encode errors instead of exiting.

- **Design decisions:**
  - Owner decision: the process exits with status 1 when a scheduled database health check fails. The first check runs 30 s after start, then every 5 minutes. A database that is down at startup therefore ends the process about 30 s after it starts.
  - No Docker healthcheck or restart policy in GO-001. Both are moved to GO-002 (DOING).
  - `/health` bounds the database ping with a 2-second timeout derived from `r.Context()`, so a client that disconnects cancels the ping.
  - The health handler depends on `*pgxpool.Pool` directly. The `Pinger` interface and fake ping are deprecated for this item.
  - `ConnectDB` takes no context and returns an error. `pgxpool.New` only parses the URL and builds the pool; the database is not dialed at startup.
  - A failed response write in the handler is logged, not fatal, so one client disconnect cannot stop the process.
  - Configuration: `DATABASE_URL` is required; `PORT` defaults to 4500. `.env` is loaded with `godotenv.Load()`, a missing file is ignored, and shell variables win.
  - Shutdown: on SIGINT or SIGTERM, `Shutdown` waits up to 10 s for in-flight requests. Deferred calls then close the pool and release the signal context.
  - The Compose file is separate from the root Compose file, and contains only `postgres` (image `postgres:16`).

- **Go concepts learned:**
  - Solid (learning check, 2026-10-17): context deadlines and cancellation (`WithTimeout`, `r.Context()`, `select` on `ctx.Done()`); returning or logging errors versus `log.Fatal`/`os.Exit`; deferred calls run in reverse order and are skipped by `os.Exit`; the order of graceful shutdown.
  - Needs more practice: what `Shutdown` does when a request outlasts its deadline; `pool.Close()` blocking until connections are returned; checking library source or `go doc` before predicting behaviour, instead of assuming.
  - Also covered in the code: `signal.NotifyContext`, `httptest`, table-driven tests with `t.Setenv`, `ServeMux` method patterns (405 and 404), and `godotenv` precedence.

- **Remaining technical debt:**
  - AC7's 5-minute exit path has not been run. Verifying it needs a 5-minute wait or a stopped database. Owner decision (2026-10-17): accepted for now, not verified before closure.
  - The background database check goroutine in [main.go](go-backend/cmd/api/main.go) sleeps 30 s without listening for shutdown. It ends with the process, so this is harmless today, but it should be tidied when the health check is revisited.
  - `GET /health` logs an info line on every request, which will be noisy under polling.
  - The slow-ping test has not been mutation-checked, so it has not been shown to fail when the timeout is removed.
  - Docker: the `postgres` service has no healthcheck and no restart policy (GO-002). The `postgres-db` container was left running after the AC6 check.
  - Naming: the database parameter is `dbUrl`; Go style usually writes `dbURL`. Not changed.

- **Suggested next task (not started):** GO-002, Compose healthcheck and restart policy (DOING). Its open questions need the owner's answers before approval.
