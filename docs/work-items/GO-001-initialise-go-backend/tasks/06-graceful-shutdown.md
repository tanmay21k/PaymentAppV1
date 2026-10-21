# Task 06: Shut down gracefully

> Retro-written on 2026-10-21 from the GO-001 closure evidence. The redo checklist is a draft for you to confirm or edit.

- **Item:** GO-001
- **Status:** PASSED
- **Covers:** AC10

## Goal
SIGINT or SIGTERM stops new requests, lets in-flight requests finish for up to 10 seconds, closes the pool, and exits 0.

## Concepts
Running `http.Server` in a goroutine, `signal.NotifyContext`, `Server.Shutdown` with a deadline (it stops waiting at the deadline and does not kill requests), and deferred calls running only when `main` returns.

## Check
Ctrl+C logs `shutting down` and `shutdown complete`, and `echo $?` prints 0.

## Result
2026-10-17: SIGINT logged `shutdown complete` and the process exited 0.

## Redo checklist
1. Create an `http.Server` with `Addr` set to `":" + port` and the mux as `Handler` — this is the object you shut down later.
2. Run `server.ListenAndServe()` in a goroutine and ignore `http.ErrServerClosed` — a clean shutdown makes `ListenAndServe` return that error.
3. Create `sigCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)` and `defer stop()` — the context is cancelled when a signal arrives.
4. Block on `<-sigCtx.Done()` — `main` waits here while the server runs.
5. Create a 10-second `shutdownCtx` (with `defer cancel`) and call `server.Shutdown(shutdownCtx)`; log an error but continue — `Shutdown` waits up to the deadline.
6. Let `main` return so deferred calls such as `pool.Close()` run — `os.Exit` and `log.Fatal` would skip them.
