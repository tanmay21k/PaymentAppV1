# Go Learning Log

One entry per closed work item, appended by `/close-work-item`. Keep entries short. Suggested follow-ups are proposals, not assignments.

## Entry template

### GO-NNN — <title> (YYYY-MM-DD)
- **Concepts encountered:**
- **Understood well:**
- **Needs more practice:**
- **Mistakes and lessons:**
- **Design decisions and rationale:**
- **Follow-up exercises (optional, only if owner agrees):**

---

### GO-001 — Initialise the Go backend (2026-10-17)
- **Concepts encountered:** modules and `cmd/` and `internal/` layout; environment configuration with `os.Getenv` and `strconv`; `.env` loading with `godotenv.Load` (shell variables win); `pgxpool` creation and `Ping`; `net/http` `ServeMux` method patterns; `httptest` tests and a real TCP listener for a slow dependency; context deadlines and cancellation (`WithTimeout`, `r.Context()`); graceful shutdown with `signal.NotifyContext` and `Server.Shutdown`; deferred calls; `log.Fatal` versus returning errors.
- **Understood well:** context deadlines stop the caller waiting and do not kill the work by themselves; `r.Context()` cancels when the client leaves; returning or logging errors keeps the server up, while `log.Fatal` exits and skips deferred calls; deferred calls run in reverse order and are skipped by `os.Exit`; the order of graceful shutdown (drain, return, deferred cleanup).
- **Needs more practice:** what `Shutdown` does when a request outlasts its deadline; `pool.Close()` blocks until connections are returned; checking `go doc` or the library source before predicting behaviour; "skipped" versus "not yet registered" for a defer; exit codes as a signal to the caller; the difference between `os.Getenv` and `os.LookupEnv`.
- **Mistakes and lessons:** several answers were guesses that the code contradicted: a timeout that terminated the process (it only stops the wait); a 15-second `/health` request, which the 2-second timeout makes impossible; a 503 from a request the process never finished. Each was corrected by reading the code or the library source. Lesson: read the code or `go doc` first, then predict.
- **Design decisions and rationale:** the process exits when a scheduled database health check fails (owner decision, first check 30 s after start, then every 5 minutes); no Docker healthcheck or restart policy in GO-001 (moved to GO-002); `/health` uses a 2-second ping timeout from `r.Context()`; no `Pinger` interface (fake ping deprecated); `ConnectDB` takes no context and returns an error; handler write errors are logged, not fatal.
- **Follow-up exercises (optional, only if owner agrees):** quiz topics 5 (shutdown with a deadline) and 6 (goroutine timing) in `docs/learning/QUIZ-TOPICS.md`, then topic 7 (pool semantics) by reading the source first.

---
