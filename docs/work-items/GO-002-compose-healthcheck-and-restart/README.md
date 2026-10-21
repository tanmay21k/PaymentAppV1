# GO-002: Compose healthcheck and restart policy

- **Status:** TODO
- **Created:** 2026-10-17
- **Type:** feature
- **Depends on:** GO-001 (`go-backend/docker-compose.yml` and the Go service entrypoint)
- **Revised:** 2026-10-21 — owner decision: status set to TODO (compose changes are in the root `docker-compose.yml` and will be merged to main later)

## Problem
GO-001 creates `go-backend/docker-compose.yml` with only a `postgres` service and no healthcheck. The Go process exits about 30 seconds after start when its database check fails, and no restart policy brings it back [Confirmed: GO-001 Contract and behavior]. Nothing in Compose tells the Go service to wait for Postgres to be ready, or restarts it after it exits.

## Scope
- A Compose healthcheck on the `postgres` service that reports whether the database accepts connections. [Assumption: a `pg_isready`-style check; the exact command is chosen in the first task.]
- A Go service entry in `go-backend/docker-compose.yml` that depends on the healthy `postgres` service (`depends_on` with `condition: service_healthy`).
- A restart policy on the Go service so that an exit from a failed database check is recovered. [Assumption: `restart: unless-stopped`, matching the existing `postgres` entry.]
- A Go Dockerfile, if the Go service runs in Compose (open question 1).

## Non-goals
- Changing the root `docker-compose.yml` or the JS stack.
- Changing the Go health handler or the exit-on-failure behaviour from GO-001.
- Any API endpoint from `openapi.yaml`.

## Contract and behavior
- **Schemas and responses:** none.
- **OpenAPI:** none. Compose and container behaviour only.
- **JS reference:** `docker-compose.yml:7`, `:25`, `:37` — the existing healthcheck and `depends_on` in the root Compose file, used as a style reference only.
- **Behavior per case:**

| Case | Result | Source |
| --- | --- | --- |
| `docker compose ... ps` once Postgres accepts connections | `postgres` shows `healthy` | [Assumption] |
| Go service started with `postgres` not yet ready | It waits for the healthy gate, then starts | [Assumption] |
| Postgres stopped after startup | The Go process exits 1, Compose restarts it, and it keeps exiting until Postgres returns | [Assumption] |

- **Discrepancies (spec vs JS):** none.
- **Design note:** a healthcheck on the Go container is not proposed. It would only poll `/health`, and the process already exits on database failure. The restart policy is the recovery path.

## Acceptance criteria
- [ ] **AC1:** `docker compose -f go-backend/docker-compose.yml up -d postgres` followed by `docker compose -f go-backend/docker-compose.yml ps` shows `postgres` as `healthy`. Verify with: command output.
- [ ] **AC2:** With the Go service in Compose, it starts only after `postgres` is healthy. Verify with: `docker compose ... logs` ordering or `ps` output.
- [ ] **AC3:** If the Go process exits with status 1 after a failed database check, Compose restarts it according to the policy. Verify with: `docker compose ... ps` and the logs after stopping Postgres.

## Open questions
1. Does the Go service run in Compose in this item (which needs a Dockerfile), or is that a later item?
2. Restart policy value: `unless-stopped` (matches `postgres`) or `on-failure`?
3. Healthcheck command for `postgres`: which approach?

Risk: restart loops while Postgres is down, each restart taking about 30 seconds to fail.

## Tasks
Tasks are added one at a time with `/next-task GO-002`. Each task is a file in `tasks/`; a task is ticked when `/check-task` passes it.

- [ ] (none yet)

## Closure
- **Evidence:**
- **Decisions:**
- **Debt and follow-ups:**
