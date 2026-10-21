# GO-001: Initialise the Go backend

- **Status:** CLOSED
- **Created:** 2026-10-17
- **Type:** feature
- **Depends on:** none
- **Note:** built and closed before the task format existed. The task files below were written afterwards, on 2026-10-21, from the closure evidence, so they are a logical order and not the order of the commits. Their redo checklists are drafts for you to confirm or edit.

## Problem
`go-backend/` had no Go module, entrypoint, configuration loading, or database connection to build the migration on. The earlier Go code (a PostgreSQL user repository, a sqlc setup and a JWT helper) was deleted in commit `e1e207f`. This item creates the smallest working foundation: a service that builds, reads configuration, connects to PostgreSQL, serves `GET /health`, and shuts down cleanly. No API endpoint from `openapi.yaml` is implemented here.

## Scope
- Go module `github.com/tanmay21k/PaymentAppV1/go-backend`, with `cmd/api/main.go` as the entrypoint and packages under `internal/` (`config`, `database`, `health`).
- Configuration from `PORT` (default 4500) and `DATABASE_URL` (required), validated at startup.
- `.env` loading with `github.com/joho/godotenv` (new dependency; the standard library has no `.env` parser). A missing file is ignored and shell variables take precedence.
- A `pgxpool` connection pool from `DATABASE_URL` (`pgx/v5`).
- `GET /health` on a `net/http` `ServeMux`, with a 2-second database ping timeout.
- A background database check that exits the process when it fails.
- Graceful shutdown on SIGINT and SIGTERM.
- `go-backend/.env.example`, `go-backend/.gitignore` (ignores `.env`), and `go-backend/docker-compose.yml` with a `postgres:16` service and a named volume.
- Unit tests for configuration and for the health handler.

## Non-goals
- Any endpoint from `openapi.yaml`, authentication, or JWT handling.
- Database migrations, schema, or any table.
- A Docker healthcheck for `postgres` or a restart policy for the Go service (moved to GO-002).
- CORS for `http://localhost:5173`, a Go Dockerfile, a logging framework, or a dependency-injection library.
- Changes to `backend/`, `frontend/`, or the root `docker-compose.yml`.

## Contract and behavior
- **Schemas and responses:** none.
- **OpenAPI:** none. `GET /health` is an operational endpoint that is not in `openapi.yaml`. The default port matches the server URL `http://localhost:4500` (`openapi.yaml:8-9`).
- **JS reference:** `backend/server.js:10` (JS `PORT` defaults to 3000; the Go default is 4500 and the JS default is not carried over); `backend/server.js:13-15` (CORS, not in scope); `.env.example:1-3` (`PORT`, `MONGO_URI`, `JWT_SECRET`; Go uses `DATABASE_URL` instead of `MONGO_URI`).
- **Behavior per case:**

| Case | Result | Source |
| --- | --- | --- |
| `PORT` unset | Listens on 4500 | [Confirmed] AC5 |
| `DATABASE_URL` unset | Exits 1, error names `DATABASE_URL`, before listening | [Confirmed] AC4 |
| `DATABASE_URL` malformed | Exits 1 with the parse error | [Confirmed] observed |
| `PORT` non-numeric or outside 1 to 65535 | Exits 1, error names `PORT` | [Confirmed] config tests |
| `.env` missing | Not an error; shell variables are used | [Confirmed] AC14 |
| `.env` malformed | Exits 1 with the godotenv error | [Confirmed] observed |
| Variable set in both shell and `.env` | The shell value wins | [Confirmed] AC15 |
| `GET /health`, database answers | 200 `{"status":"ok","database":"ok"}` | [Confirmed] AC6 |
| `GET /health`, database fails or hangs | 503 `{"status":"unavailable","database":"unreachable"}` within 2 seconds | [Confirmed] tests, AC7 |
| `POST /health` | 405 | [Confirmed] AC8 |
| Unknown path | 404 | [Confirmed] AC9 |
| Database check fails (30 s after start, then every 5 minutes) | Logs `database health check failed` and exits 1 | [Confirmed: owner decision] AC7 |
| SIGINT or SIGTERM | Stops new requests, drains up to 10 s, closes the pool, exits 0 | [Confirmed] AC10 |

- **Discrepancies (spec vs JS):** none.

## Acceptance criteria
- [x] **AC1:** `go build ./...` in `go-backend/` exits 0. Verify with: command output.
- [x] **AC2:** `go vet ./...` exits 0. Verify with: command output.
- [x] **AC3:** `go test ./...` exits 0 and runs the config and health tests. Verify with: `go test ./... -v` output.
- [x] **AC4:** With `DATABASE_URL` unset, `go run ./cmd/api` exits non-zero with a message containing `DATABASE_URL`, before listening on any port. Verify with: command output and `echo $?`.
- [x] **AC5:** With `PORT` unset, the service listens on 4500. Verify with: the startup log and `curl -i localhost:4500/health`.
- [x] **AC6:** With Postgres running and `DATABASE_URL` set, `curl -i localhost:4500/health` returns 200 and `{"status":"ok","database":"ok"}`. Verify with: curl output.
- [ ] **AC7:** With Postgres stopped and the service running, `/health` returns 503 and `{"status":"unavailable","database":"unreachable"}` while the process is alive, and the process then exits 1 with `database health check failed` at the next scheduled check (about 30 s after start if Postgres is already down, otherwise within 5 minutes). Verify with: curl output, log lines and `echo $?`. Partly verified; the 5-minute path was accepted unverified (see Closure).
- [x] **AC8:** `curl -i -X POST localhost:4500/health` returns 405. Verify with: curl output.
- [x] **AC9:** `curl -i localhost:4500/nope` returns 404. Verify with: curl output.
- [x] **AC10:** Ctrl+C on a running service logs a shutdown message and the process exits 0. Verify with: terminal output and `echo $?`.
- ~~**AC11:** Compose `postgres` service shows `healthy`.~~ Removed 2026-10-17: the healthcheck moved to GO-002. The number is kept so references stay stable.
- [x] **AC12:** `.env.example` lists `PORT` and `DATABASE_URL` with no real secrets, and `.gitignore` ignores `.env`. Verify with: file contents and `git check-ignore go-backend/.env`.
- [x] **AC13:** No file outside `go-backend/` and `docs/work-items/` is changed by this item. Verify with: `git status --short`.
- [x] **AC14:** With `.env` containing `PORT=4600` and `PORT` unset in the shell, the service logs port 4600. With no `.env` and `PORT` and `DATABASE_URL` exported, it starts normally. Verify with: the startup log and curl output.
- [x] **AC15:** With `.env` containing `PORT=4600` and `PORT=4700` exported in the shell, the service logs port 4700. Verify with: the startup log.
- [x] **AC16:** `go.mod` lists `github.com/joho/godotenv` as a direct requirement and `go build ./...` succeeds with `go.sum` present. Verify with: `grep godotenv go-backend/go.mod` and the build output.

## Open questions
- None. All were resolved during the item: module name and variable names confirmed, the process exits when a database check fails, the Compose healthcheck and restart policy moved to GO-002, and the fake ping (`Pinger` interface) was dropped.

## Tasks
Written afterwards from the closure evidence (see the Note above).

- [x] [01 Create the module and skeleton](tasks/01-module-and-skeleton.md)
- [x] [02 Read and validate configuration](tasks/02-configuration.md)
- [x] [03 Load .env with godotenv](tasks/03-dotenv.md)
- [x] [04 Create the connection pool](tasks/04-connection-pool.md)
- [x] [05 Add GET /health](tasks/05-health-endpoint.md)
- [x] [06 Shut down gracefully](tasks/06-graceful-shutdown.md)
- [x] [07 Exit when the database check fails](tasks/07-database-check-exit.md)
- [x] [08 Run Postgres in Compose](tasks/08-postgres-compose.md)

## Closure
- **Evidence:** (final run, 2026-10-17)
  - `go build ./...` and `go vet ./...` exit 0; `gofmt -l .` prints nothing; `go test ./... -count=1 -v` passes (`TestLoad` with 8 subtests, `TestHealthReturns503WhenDatabaseUnreachable`, and `TestHealthReturns503WhenDatabaseIsSlow` at 2.00 s).
  - AC4: `invalid configuration: DATABASE_URL is not set`, exit 1, no listening line.
  - AC5, AC8, AC9: on port 4500 `/health` returned 200, `POST /health` 405, and `/nope` 404.
  - AC6: the `postgres:16` container `postgres-db` accepted connections, and `curl -i localhost:4500/health` returned `HTTP/1.1 200 OK` with `{"status":"ok","database":"ok"}`.
  - AC10: SIGINT logged `shutdown complete` and exited 0.
  - AC14, AC15: `.env` `PORT=4600` with the shell unset listened on 4600; no `.env` with exported variables listened on 4600; shell `PORT=4700` overrode `.env` and listened on 4700.
  - AC16: `go.mod` lists `github.com/joho/godotenv v1.5.1` as a direct requirement.
  - AC7 (partly): with the database down, `/health` returned 503 while the process was alive and the process exited 1 at the first check (about 28 s after start, observed on an earlier build). The 5-minute path is verified by reading only. Owner decision (2026-10-17): accepted without verification.
  - Malformed `DATABASE_URL` exits 1 with the parse error. A silent database gets a 503 in 2.00 s.
  - Review: PASS WITH NON-BLOCKING OBSERVATIONS (final pass, 2026-10-17).
- **Decisions:**
  - The process exits with status 1 when a scheduled database health check fails (owner decision). The first check runs 30 s after start, then every 5 minutes.
  - No Docker healthcheck or restart policy here; both are in GO-002.
  - `/health` bounds the ping with a 2-second timeout derived from `r.Context()`, so a client that disconnects cancels the ping.
  - The health handler depends on `*pgxpool.Pool` directly; the `Pinger` interface and fake ping were dropped.
  - `ConnectDB(dbUrl)` takes no context and returns an error; `pgxpool.New` only parses the URL and does not dial.
  - A failed response write in the handler is logged, not fatal, so one client disconnect cannot stop the process.
  - `.env` is loaded with `godotenv.Load()` from the working directory (run from `go-backend/`); it does not override existing variables.
  - Shutdown drains for up to 10 s, then deferred calls close the pool.
  - The Compose file is separate from the root one and holds only `postgres:16`.
- **Debt and follow-ups:**
  - AC7's 5-minute exit path has not been run.
  - The background check goroutine sleeps 30 s without listening for shutdown. It ends with the process, so it is harmless today.
  - `GET /health` logs an info line on every request, which will be noisy under polling.
  - The slow-ping test has not been mutation-checked.
  - The `postgres-db` container was left running after the AC6 check.
  - The parameter is named `dbUrl`; Go style usually writes `dbURL`.
  - Concepts and practice areas from the learning check are in `docs/learning/LEARNING-LOG.md` (GO-001 entry) and `docs/learning/QUIZ-PROGRESS.md`.
