# GO-004: Set up sqlc and create the users and accounts tables

- **Status:** TODO
- **Created:** 2026-10-21
- **Type:** feature
- **Depends on:** GO-001 (Go service and PostgreSQL pool in `go-backend/internal/database/db.go`)

## Problem
The Go backend has a database pool but no schema and no database-access code: GO-001 deferred migrations. The signup migration (GO-003) needs `users` and `accounts` tables and a safe, typed way to insert into them. This item adds the migration tooling, the two tables, and sqlc-generated insert functions, without touching any endpoint.

## Scope
- Install sqlc and a migration tool as local command-line tools. [Assumption: sqlc local only, goose as the migration tool; see open questions 1 and 2]
- Folder layout under `go-backend/`: `db/migration/` for migration files, `db/queries/` for SQL queries, and sqlc-generated Go code in `db/` as package `db`. [Confirmed: owner, 2026-10-21; base folder `go-backend/` is an Assumption]
- One migration (with an up and a down) that creates the `users` and `accounts` tables.
- A sqlc config file, and two queries that insert a user and insert an account, each returning the new row.
- Run `sqlc generate` so the generated code compiles with the existing `pgx/v5` pool.
- A Go test, using the generated code against the `postgres-db` container, that proves both inserts work and that a duplicate username is rejected.

## Non-goals
- The signup endpoint, `main.go` routes, password hashing, and JWT (GO-003).
- Any other table (for example transactions), or reading and listing queries beyond the two inserts.
- Migrating MongoDB data.
- Changes to `docker-compose.yml`, `backend/`, or `frontend/`.
- CI wiring, a `Makefile`, or auto-running migrations at service start.

## Contract and behavior
- **Schemas and responses:** none.
- **OpenAPI:** none. No HTTP change. The `balance` response in `openapi.yaml` is `type: number` with decimal values (example `2450.75`), which matters for the balance column type (open question 4).
- **JS reference:** `backend/config/userschema.js:3-30` (user fields, `username` unique, `paymentPin` and `profilePicture` default `""`); `backend/config/bankschema.js:4-14` (account has a required `userId` and a required numeric `balance`); `backend/routes/user.js:53-56` (signup creates one account per user).
- **Behavior per case:**
  - Applying the migration up creates `users` and `accounts`. [Assumption: all columns `NOT NULL` unless stated]
  - `users` holds: an id, `firstname`, `lastname`, `username` (unique), `password_hash` (the bcrypt output, GO-003 produces it), and a created-at timestamp. [Assumption]
  - `accounts` holds: an id, a `user_id` that references `users`, and a `balance`. A user has at most one account. [Assumption: `user_id` is unique, matching one account per signup]
  - Applying the migration down removes both tables in an order the foreign key allows. [Assumption: reversible]
  - Inserting a second user with an existing `username` fails with a unique-violation error. [Confirmed: `username` is `unique` in `userschema.js:12-16`]
  - The generated functions `CreateUser` and `CreateAccount` take typed arguments and return the inserted row. [Assumption: query names]
- **Discrepancies (spec vs JS):** none for this item. JS fields `paymentPin` and `profilePicture` are not in `SignupRequest`; see open question 5.

## Acceptance criteria
- [ ] **AC1:** `sqlc version` and the migration tool's version command each print a version. Verify with: command output.
- [ ] **AC2:** Applying the migration up exits 0, and in `psql`, `\d users` and `\d accounts` show the columns, the unique `username`, and the foreign key from `accounts.user_id` to `users`. Verify with: migration command output and `psql` output.
- [ ] **AC3:** Applying the migration down removes both tables, and applying up again recreates them. Verify with: `psql` `\dt` output after each step.
- [ ] **AC4:** `sqlc generate` exits 0 and writes Go files in `go-backend/db/` in package `db`, including `CreateUser` and `CreateAccount`. Verify with: command output and `ls go-backend/db`.
- [ ] **AC5:** `go build ./...` and `go vet ./...` in `go-backend/` exit 0 with the generated code present. Verify with: command output.
- [ ] **AC6:** A Go test named `TestCreateUserAndAccount` inserts a user and then an account for that user and reads the values back; a second `CreateUser` with the same `username` returns an error. Verify with: `go test ./db/ -run TestCreateUserAndAccount -v` output and a `psql` row count.
- [ ] **AC7:** `go test ./...` in `go-backend/` passes. Verify with: command output.
- [ ] **AC8:** This item does not change `go-backend/cmd/api/main.go` or add a route. Verify with: `git diff --stat` showing no change to `cmd/api/main.go` from this item.

## Open questions
1. **Migration tool:** goose (the choice noted in GO-001) or another? Recommend goose. [Owner decision needed]
2. **Install sqlc how:** as a local command-line tool only (recommended, nothing added to `go.mod`), or as a tracked tool in `go.mod`? [Owner decision needed]
3. **Id type:** a 64-bit integer identity (recommended, simplest with sqlc and pgx) or UUID? This also answers GO-003 open question 4. [Owner decision needed]
4. **Balance type:** `NUMERIC(12,2)` (recommended: keeps parity with the two-decimal values JS stores, and the spec returns a decimal `number`) or integer cents (simpler Go type, but the API would need a conversion)? [Owner decision needed]
5. **JS fields `paymentPin` and `profilePicture`:** leave them out of `users` now, since signup does not use them (recommended), or add them now with empty defaults? [Owner decision needed]
6. **Overlap with GO-003:** GO-003's scope and its open question 7 also cover the tables and the schema tool. If this item goes ahead, table creation moves here and GO-003 depends on GO-004. Edit GO-003 now, or after approval? [Owner decision needed]

## Tasks
Tasks are added one at a time with `/next-task GO-004`. Each task is a file in `tasks/`; a task is ticked when `/check-task` passes it.

- [ ] (none yet)

## Draft sequence (from the earlier plan; not tasks yet)
Titles only. The coach rewrites each one as a task with `/next-task GO-004`, after asking you questions.

1. Install sqlc.
2. Install the migration tool.
3. Create the folders.
4. Write the migration.
5. Apply the migration up.
6. Prove the migration is reversible.
7. Write the queries.
8. Write the sqlc config.
9. Generate and read the code.
10. Prove it works against the database.

## Closure
- **Evidence:**
- **Decisions:**
- **Debt and follow-ups:**
