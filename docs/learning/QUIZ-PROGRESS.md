# Go Quiz Progress

Written by the `go-quizzer` agent after each quiz session. Topics are numbered as in [QUIZ-TOPICS.md](QUIZ-TOPICS.md). The learning log ([LEARNING-LOG.md](LEARNING-LOG.md)) is still appended only at work-item closure.

Status: **Solid** (answered correctly without prompting), **Practice** (needs another attempt), **Not yet** (not assessed).

| # | Topic | Status | Last assessed | Notes |
|---|---|---|---|---|
| 1 | Context: deadlines and cancellation | Solid | 2026-10-17 | Practice the parent-child relationship of contexts |
| 2 | `select` and `ctx.Done()` | Solid | 2026-10-17 | Toy example only |
| 3 | Errors: return, log, or exit | Practice | 2026-10-17 | Solid on return versus log; exit codes as a signal to the caller were muddled at first |
| 4 | Defer: order and when it runs | Practice | 2026-10-17 | Solid on reverse order and `log.Fatal` skipping defers; "skipped" versus "not yet registered" needs work |
| 5 | Graceful shutdown | Practice | 2026-10-17 | Ordering solid; `Shutdown` past its deadline needs work |
| 6 | Goroutines and background work | Not yet | | Confused "starts" with "first action" in the startup trace |
| 7 | Connection pools: creating versus connecting | Practice | 2026-10-17 | Answers were guesses until the source was read |
| 8 | Environment configuration | Not yet | | AC14 and AC15 evidence exists; no quiz yet |
| 9 | Testing | Not yet | | |
| 10 | Routing with `net/http` | Not yet | | AC8 and AC9 evidence exists; no quiz yet |
| 11 | Interfaces (deferred) | Not yet | | |
| 12 | Modules, packages and `internal/` | Not yet | | |

## Session log

### 2026-10-17 — GO-001 learning check (go-mentor)
- Source: `go-mentor` teaching session, not the quiz agent.
- Scenarios covered: deadline and cancellation, error handling, defer order, graceful shutdown, pool semantics.
- Learner habits seen: guessing library behaviour before checking. Asked to read `go doc` or the source before answering.
- Next session: topics 5 and 6 first.
