# Task 04: Create the connection pool

> Retro-written on 2026-10-21 from the GO-001 closure evidence. The redo checklist is a draft for you to confirm or edit.

- **Item:** GO-001
- **Status:** PASSED
- **Covers:** AC4

## Goal
`ConnectDB` builds a `pgxpool` pool from `DATABASE_URL`. A malformed URL stops startup. A well-formed URL to an unreachable database does not, because the pool does not dial at creation.

## Concepts
`pgxpool.New` (it parses the URL and does not connect), why reachability is checked with a ping, `defer pool.Close()`, and `Close` blocking until connections are returned.

## Check
A malformed `DATABASE_URL` exits 1 with the parse error, and a well-formed URL to a closed port starts and listens.

## Result
2026-10-17: `DATABASE_URL="not a url"` exited 1 with `cannot parse`. A well-formed URL to a closed port logged `listening on port 4500`.

## Redo checklist
1. Run `go get github.com/jackc/pgx/v5` — `pgxpool` lives in this module.
2. Write `ConnectDB(dbUrl string) (*pgxpool.Pool, error)` in `internal/database` that calls `pgxpool.New(context.Background(), dbUrl)` and returns the error — no `log.Fatal` in library code.
3. In `main`, handle the error and then `defer pool.Close()` — the pool lives for the life of the process.
4. Remember `New` does not dial the database; use `Ping` to check reachability.
