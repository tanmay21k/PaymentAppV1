# Task 05: Add GET /health

> Retro-written on 2026-10-21 from the GO-001 closure evidence. The redo checklist is a draft for you to confirm or edit.

- **Item:** GO-001
- **Status:** PASSED
- **Covers:** AC3, AC6, AC8, AC9

## Goal
`GET /health` returns 200 `{"status":"ok","database":"ok"}` when the database answers and 503 `{"status":"unavailable","database":"unreachable"}` when it does not or when it takes more than 2 seconds. Other methods return 405 and unknown paths 404.

## Concepts
`http.Handler` and `ServeHTTP`, `ServeMux` method patterns, `context.WithTimeout` derived from `r.Context()`, JSON encoding of exported fields, `httptest`, and logging (not exiting) when a response write fails.

## Check
`go test ./internal/health -v` passes (`TestHealthReturns503WhenDatabaseUnreachable` and `TestHealthReturns503WhenDatabaseIsSlow`), and against a running `postgres:16`, `curl -i localhost:4500/health` returns 200, `POST` returns 405, and `/nope` returns 404.

## Result
2026-10-17: both health tests passed (the slow case in 2.00 s). `curl` returned `HTTP/1.1 200 OK` with `{"status":"ok","database":"ok"}`; `POST /health` returned 405 and `/nope` 404.

## Redo checklist
1. Define a response struct with exported fields and `json` tags — `encoding/json` only sees exported fields.
2. Define a handler type that holds the pool and implement `ServeHTTP(w, r)` — that makes it an `http.Handler`.
3. Inside the handler, create `ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)` and `defer cancel()` — the ping stops after 2 seconds or when the client leaves.
4. Ping the pool with `ctx`; set `Content-Type`, write the status, then encode the JSON — headers must be set before `WriteHeader`.
5. If encoding fails, log the error and return — one failed write must not stop the server.
6. Register the route on `http.NewServeMux()` as `"GET /health"` — the method in the pattern gives 405 for other methods and 404 for unknown paths.
7. Test the unreachable case with a pool pointing at a closed port, and the slow case with a local TCP listener that accepts but never replies — the 200 case is checked with a real database.
