# Go Quiz Topics

Topics to quiz on, drawn from GO-001 (`go-backend/`). Each topic gives the code to look at and an example question. The quiz agent (`go-quizzer`) uses this list and records progress in [QUIZ-PROGRESS.md](QUIZ-PROGRESS.md).

Status key: **Solid** = answered correctly in the GO-001 learning check. **Practice** = needs another attempt. **Not yet** = not covered.

---

## 1. Context: deadlines and cancellation
- **Code:** [handler.go:36-37](../../go-backend/internal/health/handler.go#L36-L37) (`WithTimeout` on `r.Context()`), [main.go:81-82](../../go-backend/cmd/api/main.go#L81-L82) (shutdown deadline).
- **Ask:** What does a timeout context actually stop? Does it kill work, or tell the caller to stop waiting? Why derive from `r.Context()` rather than `context.Background()`?
- **Status:** Solid (deadline stops waiting; `r.Context()` cancels when the client leaves). Practice: the exact parent-child relationship.

## 2. `select` and `ctx.Done()`
- **Code:** no direct example in GO-001. Use the toy example from the learning check: a `select` between `time.After` and `ctx.Done()`.
- **Ask:** Which case of the `select` runs when the deadline fires first? What happens to the slow work that is still running?
- **Status:** Solid (the program printed "stopped waiting" after 2 seconds, not 5).

## 3. Errors: return, log, or exit
- **Code:** [handler.go](../../go-backend/internal/health/handler.go) (the `Encode` error branch, logged and returned), [config.go:26-28](../../go-backend/internal/config/config.go#L26-L28) (validation error returned to `main`), [main.go:22-23](../../go-backend/cmd/api/main.go#L22-L23) (`log.Fatalf` for a bad `.env`).
- **Ask:** Why should a failed write to one client not stop the process? What is the difference between returning an error and calling `log.Fatal`? What do exit codes 0 and 1 tell the caller?
- **Status:** Solid for returning and logging. Practice: exit codes as a signal to the caller (the first answer mixed up "expected or unexpected").

## 4. Defer: order and when it runs
- **Code:** [main.go:39](../../go-backend/cmd/api/main.go#L39) (`defer pool.Close()`), [main.go:67](../../go-backend/cmd/api/main.go#L67) (`defer stop()`), [main.go:82](../../go-backend/cmd/api/main.go#L82) (`defer cancelShutdown()`).
- **Ask:** In what order do deferred calls run? Which of them would be skipped by `log.Fatal`? Why is a `defer` registered after a blocking wait not yet registered while the wait is running?
- **Status:** Solid (reverse order; `log.Fatal` skips defers). Practice: the difference between "skipped" and "not yet registered".

## 5. Graceful shutdown
- **Code:** [main.go:66-86](../../go-backend/cmd/api/main.go#L66-L86): `signal.NotifyContext`, the wait on `sigCtx.Done()`, `server.Shutdown` with a 10-second context.
- **Ask:** What does `Shutdown` wait for, and when does it give up? If it gives up, are in-flight requests stopped? What happens to the process when `main` returns?
- **Status:** Solid for the ordering. **Practice:** what `Shutdown` does when a request outlasts its deadline (the answer needs the code, not a guess).

## 6. Goroutines and background work
- **Code:** [main.go:42-58](../../go-backend/cmd/api/main.go#L42-L58) (the health-check goroutine, `time.Sleep` then `NewTicker`); [main.go:70-76](../../go-backend/cmd/api/main.go#L70-L76) (server in a goroutine).
- **Ask:** When does the goroutine start: immediately or after its sleep? What ends a goroutine when the process exits? Why does the sleep ignore shutdown, and is that a problem?
- **Status:** Not yet assessed. Practice: confusing "starts" with "first action".

## 7. Connection pools: creating versus connecting
- **Code:** [db.go:13-19](../../go-backend/internal/database/db.go#L13-L19) (`pgxpool.New`), [main.go:35-39](../../go-backend/cmd/api/main.go#L35-L39).
- **Ask:** Does `pgxpool.New` connect to the database? What does `pool.Close()` wait for? Why was the 5-second context removed?
- **Status:** Practice: the answer was a guess until the source was read. Check `go doc` and the source before predicting library behaviour.

## 8. Environment configuration
- **Code:** [config.go:20-36](../../go-backend/internal/config/config.go#L20-L36) (`os.Getenv`, `strconv.Atoi`, range check), [main.go:22](../../go-backend/cmd/api/main.go#L22) (`godotenv.Load`).
- **Ask:** What is the difference between `os.Getenv` and `os.LookupEnv`? Does `godotenv.Load` override a variable already set in the shell? Why is a missing `.env` not an error?
- **Status:** Not yet assessed in the learning check. The AC14/AC15 evidence covers precedence.

## 9. Testing
- **Code:** [config_test.go](../../go-backend/internal/config/config_test.go) (table-driven, `t.Setenv`), [handler_test.go:17-43](../../go-backend/internal/health/handler_test.go#L17-L43) (`httptest.NewRecorder`), [handler_test.go:45-86](../../go-backend/internal/health/handler_test.go#L45-L86) (a real TCP listener that never replies).
- **Ask:** Why is `t.Setenv` safer than `os.Setenv` in a test? What does `httptest.NewRecorder` let you check? Why does the slow-ping test use a real listener instead of a fake?
- **Status:** Not yet assessed in the learning check.

## 10. Routing with `net/http`
- **Code:** [main.go:61-63](../../go-backend/cmd/api/main.go#L61-L63) (`http.NewServeMux`, `mux.Handle("GET /health", ...)`).
- **Ask:** What does a method in the pattern do for `POST /health`? What status does an unknown path get? Why use a new `ServeMux` rather than the default one?
- **Status:** Not yet assessed. Evidence: AC8 (405) and AC9 (404) passed.

## 11. Interfaces (deferred)
- **Code:** none in GO-001. The `Pinger` interface was deprecated; see the GO-001 closure record.
- **Ask:** When is an interface worth introducing, and when does a concrete type do the job? What is "implicit satisfaction"?
- **Status:** Not yet assessed. Revisit with the first item that needs a second implementation.

## 12. Modules, packages and `internal/`
- **Code:** [go.mod](../../go-backend/go.mod), the `internal/` directories, `cmd/api/`.
- **Ask:** What does `internal/` prevent? Why does `cmd/api` contain `main` and the other packages not? What does `go mod tidy` change?
- **Status:** Not yet assessed.

---

## Suggested order for the next quiz
1. Topic 5 (shutdown with a deadline) and topic 6 (goroutine timing), the two "Practice" areas with the most consequence.
2. Topic 7 (pool semantics), by reading the source first.
3. Topics 8 to 10, which have not been assessed.
