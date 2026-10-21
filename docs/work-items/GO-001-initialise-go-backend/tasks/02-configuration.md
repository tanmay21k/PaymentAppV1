# Task 02: Read and validate configuration

> Retro-written on 2026-10-21 from the GO-001 closure evidence. The redo checklist is a draft for you to confirm or edit.

- **Item:** GO-001
- **Status:** PASSED
- **Covers:** AC3, AC4, AC5

## Goal
`PORT` (default 4500) and `DATABASE_URL` (required) are read and validated, and a bad value stops startup with a message that names the variable.

## Concepts
`os.Getenv` versus `os.LookupEnv`, `strconv.Atoi`, returning an error instead of exiting from library code, table-driven tests, and `t.Setenv`.

## Check
`go test ./internal/config -v` passes (`TestLoad` with 8 subtests), and running the service with `DATABASE_URL` unset exits 1 with a message containing `DATABASE_URL`, before listening.

## Result
2026-10-17: `TestLoad` passed with 8 subtests. With `DATABASE_URL` unset the service printed `invalid configuration: DATABASE_URL is not set` and exited 1.

## Redo checklist
1. Define a `Config` struct and a `Load() (Config, error)` function in `internal/config` — library code returns errors, and `main` decides what to do with them.
2. Read `PORT` with `os.Getenv`; if empty use the default "4500" — `Getenv` returns an empty string for both unset and empty.
3. Parse `PORT` with `strconv.Atoi` and reject values outside 1 to 65535 — a number that parses can still be an invalid port.
4. Require `DATABASE_URL` and return an error that names it — the message is what the person running the service reads.
5. Write a table-driven test using `t.Setenv` for each case — `t.Setenv` restores the environment after each subtest.
6. In `main`, call `Load()` and `log.Fatalf` on error — `log.Fatalf` exits with status 1.
