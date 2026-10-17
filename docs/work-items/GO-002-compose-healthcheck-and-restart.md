# GO-002: Compose healthcheck and restart policy

- **Status:** DOING
- **Created:** 2026-10-17
- **Approved by owner:** no
- **Depends on:** GO-001 (`go-backend/docker-compose.yml` and the Go service entrypoint)

## 1. Title
Add a Docker Compose healthcheck for `postgres`, and a restart policy for the Go service once it runs in Compose.

## 2. Problem statement
GO-001 creates `go-backend/docker-compose.yml` with only a `postgres` service and no healthcheck. The owner decided (2026-10-17) to drop the healthcheck from GO-001 and record it as future scope [Confirmed: owner]. The Go process exits about 30 seconds after start when the database check fails, and no restart policy brings it back [Confirmed: GO-001 section 8]. Nothing in Compose yet tells the Go service to wait for Postgres to be ready, or restarts it after it exits.

## 3. Business and technical context
- GO-001 section 8: the Go process exits when a scheduled database health check fails. It does not recover on its own.
- GO-001 section 14: the Go service has no Dockerfile and no Compose entry in GO-001.
- Root `docker-compose.yml` (JS stack) already uses a healthcheck and `depends_on` [Confirmed: `docker-compose.yml:7`, `:25`, `:37`]. This item should follow the same style if it fits.

## 4. Scope
- Compose healthcheck on the `postgres` service that reports whether the database accepts connections [Assumption: `pg_isready`-style check; exact command to be chosen at approval].
- A Go service entry in `go-backend/docker-compose.yml` that depends on the healthy `postgres` service (`depends_on` with `condition: service_healthy`).
- A restart policy on the Go service so that a process exit from a failed database check is recovered [Assumption: `restart: unless-stopped`, matching the existing `postgres` entry].
- A Go Dockerfile, if the owner decides the Go service runs in Compose [Open question 1].

## 5. Non-goals
- Changing the root `docker-compose.yml` or the JS stack.
- Changing the Go health handler or the exit-on-failure behaviour in GO-001 section 8.
- Any API endpoint from `openapi.yaml`.

## 6. Relevant OpenAPI operations and schemas
- None. Compose and container behaviour only.

## 7. Relevant JavaScript implementation references
- `docker-compose.yml:7`, `:25`, `:37` — the existing healthcheck and `depends_on` in the root Compose file. Used as the style reference only.

## 8. Expected behavior
- **[Assumption]** `docker compose -f go-backend/docker-compose.yml ps` shows `postgres` as `healthy` once it accepts connections.
- **[Assumption]** The Go service starts only after `postgres` is healthy.
- **[Assumption]** If the Go process exits, Compose restarts it, subject to the restart policy.

## 9. Proposed design and architectural considerations
- Healthcheck on `postgres` before adding the Go service, so `depends_on` has a real signal to wait for.
- Restart policy is the recovery path for the exit-on-failure design. A healthcheck on the Go container is not proposed: it would only poll `/health`, and the process already exits on database failure.

## 10. Acceptance criteria
- [ ] **AC1:** `docker compose -f go-backend/docker-compose.yml up -d postgres` followed by `docker compose -f go-backend/docker-compose.yml ps` shows `postgres` as `healthy`. Verify: command output.
- [ ] **AC2:** With the Go service in Compose, it starts only after `postgres` is healthy. Verify: `docker compose ... logs` ordering or `ps` output.
- [ ] **AC3:** If the Go process exits with status 1 after a failed database check, Compose restarts it according to the policy. Verify: `docker compose ... ps` and logs after stopping Postgres.

## 11. Edge cases and error handling
| Case | Expected status/response | Source |
| --- | --- | --- |
| Postgres not ready on first `up` | Go service waits (healthy gate), then starts | [Assumption] section 8 |
| Postgres stopped after startup | Go process exits; Compose restarts it; repeated exits until Postgres returns | [Assumption] restart policy |

## 12. Testing requirements
- Manual, with the commands in section 10. Capture the output as evidence.

## 13. Relevant Go concepts to learn
- None specific. Compose and container configuration only.

## 14. Dependencies and risks
**Open questions for the owner:**
1. Does the Go service run in Compose in this item, which needs a Dockerfile, or is that a later item?
2. Restart policy value: `unless-stopped` (matches `postgres`) or `on-failure`?
3. Healthcheck command for `postgres`: confirm the approach before implementation.

**Risks:**
- Restart loops while Postgres is down, with each restart taking about 30 seconds to fail.

## 15. Definition of Done
- [ ] All acceptance criteria verified with evidence
- [ ] Tests pass (output observed)
- [ ] Review outcome PASS or PASS WITH NON-BLOCKING OBSERVATIONS; findings addressed or explicitly accepted
- [ ] Documentation updated where needed
- [ ] Owner can explain key design decisions
- [ ] Owner authorized closure
