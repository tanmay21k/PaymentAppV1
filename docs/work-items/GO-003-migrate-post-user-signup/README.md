# GO-003: Migrate POST /api/v1/user/signup to Go

- **Status:** DOING
- **Created:** 2026-10-21
- **Type:** migration
- **Depends on:** GO-001 (Go service and PostgreSQL pool in `go-backend/internal/database/db.go`)

## Problem
The Go backend only serves `/health`. The first endpoint in `openapi.yaml` (`POST /api/v1/user/signup`) still exists only in `backend/`. The frontend cannot create accounts against the Go backend until it is migrated. This is the first endpoint chosen because it is first in the spec; the owner did not name one. [Assumption: the owner confirms this endpoint in open question 6]

## Scope
- `POST /api/v1/user/signup` in `go-backend/`, matching the contract below.
- PostgreSQL tables for `users` and `accounts` that the signup writes to [Assumption: created by SQL in `go-backend/`; the method (plain SQL file or migration tool) is open question 7].
- Password hashing with bcrypt, cost 10, as the JS code does [Confirmed: `backend/routes/user.js:43`].
- JWT signing with the same `JWT_SECRET` and the same 7-day expiry [Confirmed: `backend/routes/user.js:58-62`, `openapi.yaml` AuthResponse].

## Non-goals
- `POST /api/v1/user/signin` and every other endpoint.
- Migrating existing MongoDB data.
- Changing `backend/` or `frontend/`.
- Adding a third-party validation or JWT library without owner agreement (see open questions).

## Contract and behavior

- **OpenAPI:** `POST /api/v1/user/signup`. Request schema: `SignupRequest`. Responses: 200 `AuthResponse`, 400 `MessageResponse`, 500 `ServerError`.
- **JS reference:** `backend/routes/user.js:15-20` (zod schema), `:22-72` (handler).
- **Behavior per case:**

| Case | Status | Response body | Source |
| --- | --- | --- | --- |
| Valid body, new username | 200 | `{"msg":"User created successfully","token":"<JWT>"}` | [Confirmed] `user.js:64-67`, spec 200 |
| Body fails validation (missing field, bad email, password under 6 chars) | 400 | `{"msg":"Invalid inputs"}` | [Confirmed] `user.js:27-30`; spec 400 |
| Username already exists | 400 | `{"msg":"User already exists"}` | [Confirmed] `user.js:37-41`; spec 400 describes this case but gives no message text |
| Database or other unexpected error | 500 | `{"msg":"Server error"}` | [Confirmed] `user.js:68-71`; spec 500 |
| Success token | — | JWT signed HS256 with `JWT_SECRET`, payload `userId`, expires in 7 days | [Confirmed] `user.js:58-62`; [Assumption] Go payload key and id type, see open question 4 |
| Account starting balance | — | Random value from 1.00 to 10001.00, two decimals | [Confirmed] `user.js:55`; the spec does not define it |

- **Discrepancies (spec vs JS):**
  1. Duplicate-username message: the spec describes the case but does not fix the text. JS returns `"User already exists"`. [Owner decision needed: open question 1]
  2. Signup writes the user and the account in two separate steps with no transaction (`user.js:44-56`). If the account insert fails, a user with no account is left behind. The spec does not say what should happen. [Owner decision needed: open question 2]
  3. The JS code compares usernames exactly as sent, with no lowercasing or trimming (`user.js:33-35`). The spec does not say whether `User@x.com` and `user@x.com` are the same user. [Owner decision needed: open question 3]

## Acceptance criteria
- [ ] **AC1:** A valid signup returns HTTP 200 and JSON with `msg` equal to `User created successfully` and a non-empty `token`. Verify with: `curl -i -X POST localhost:4500/api/v1/user/signup -H "Content-Type: application/json" -d '{"firstname":"A","lastname":"B","username":"a@example.com","password":"secret1"}'`, output captured.
- [ ] **AC2:** After AC1, the `users` table has one row for `a@example.com`, and the `accounts` table has one row for that user with a balance from 1.00 to 10001.00 and two decimals. Verify with: `psql` query output.
- [ ] **AC3:** The stored password is a bcrypt hash (prefix `$2a$10$` or `$2b$10$`) and `bcrypt.CompareHashAndPassword` accepts `secret1` for it. Verify with: `psql` output and a Go test named `TestSignupStoresBcryptHash`.
- [ ] **AC4:** The token decodes with `JWT_SECRET`, uses HS256, and expires 7 days after `iat`. Verify with: a Go test named `TestSignupTokenClaims`, output captured.
- [ ] **AC5:** Each validation failure in the table (missing `firstname`, `username` not an email, `password` of 5 characters) returns 400 and `{"msg":"Invalid inputs"}`, and writes no rows. Verify with: `curl` for each case plus a `psql` row count before and after.
- [ ] **AC6:** A second signup with the same `username` returns 400 with the message decided in open question 1, and writes no rows. Verify with: `curl` twice and a `psql` row count.
- [ ] **AC7:** With PostgreSQL stopped, signup returns 500 and `{"msg":"Server error"}`. Verify with: `docker compose ... stop postgres`, then `curl`, output captured.
- [ ] **AC8:** If the account insert fails, the user insert is rolled back (per open question 2). Verify with: a Go test that forces the account insert to fail, and a `psql` row count afterwards.
- [ ] **AC9:** `go test ./...` in `go-backend/` passes. Verify with: command output.
- [ ] **AC10:** No other route is added. Verify with: `grep` of the route registration in `go-backend/` showing only `/health` and `/api/v1/user/signup`.

## Open questions
1. **Duplicate-username message:** keep `"User already exists"` (matches JS, what the frontend currently receives), or use a new wording from the spec description?
2. **Transaction:** Go wraps the user and account inserts in one database transaction, so a failure leaves neither row (recommended), or mirror the JS two-step behavior?
3. **Username case:** exact match as in JS (recommended for parity), or lowercase before storing and checking?
4. **Token payload:** keep the `userId` key and send the Postgres id? The Postgres id type is not chosen yet (integer or UUID). Tokens from the JS backend use a Mongo ObjectId, so a JS-issued token will not map to a Go user. Is that acceptable during the migration?
5. **Validation library:** use a Go validation library, or hand-written checks with the standard library (recommended, no new dependency)? The rule is no new dependency without agreement.
6. **Endpoint choice:** confirm `POST /api/v1/user/signup` as the first endpoint, or name another.
7. **Schema tool:** plain SQL file run by hand or in Compose, or a migration tool (adds a dependency)?

## Tasks
Tasks are added one at a time with `/next-task GO-003`. Each task is a file in `tasks/`; a task is ticked when `/check-task` passes it.

- [ ] (none yet)

## Closure
- **Evidence:**
- **Decisions:**
- **Debt and follow-ups:**
