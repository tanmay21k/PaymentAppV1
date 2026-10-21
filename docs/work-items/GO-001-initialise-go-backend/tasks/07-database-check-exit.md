# Task 07: Exit when the database check fails

> Retro-written on 2026-10-21 from the GO-001 closure evidence. The redo checklist is a draft for you to confirm or edit.

- **Item:** GO-001
- **Status:** PASSED
- **Covers:** AC7

## Goal
A background check pings the database 30 seconds after start and then every 5 minutes. A failed check logs `database health check failed` and exits the process with status 1 (owner decision).

## Concepts
Goroutines, `time.Sleep` and `time.NewTicker`, a `context.WithTimeout` per check, `os.Exit` skipping deferred calls, and why a restart policy (GO-002) is the recovery path.

## Check
With an unreachable database the service answers `/health` with 503 and exits 1 about 30 seconds after start.

## Result
2026-10-17: against an unreachable database `/health` returned 503, then the process logged `database health check failed` and exited about 28 seconds after start. The 5-minute path was not run (accepted).

## Redo checklist
1. Start a goroutine from `main` — it runs beside the server.
2. In it, `time.Sleep(30 * time.Second)`, then create `time.NewTicker(5 * time.Minute)` with `defer ticker.Stop()` — the first check waits for the database to settle.
3. On each tick, build a 5-second `context.WithTimeout` and ping the pool — the same ping the handler uses.
4. On error, log it and call `os.Exit(1)`; on success, log and wait for the next tick — the exit is deliberate and skips deferred calls.
5. Do not add a restart policy here — that is GO-002.
