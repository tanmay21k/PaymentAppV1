# Task 08: Run Postgres in Compose

> Retro-written on 2026-10-21 from the GO-001 closure evidence. The redo checklist is a draft for you to confirm or edit.

- **Item:** GO-001
- **Status:** PASSED
- **Covers:** AC6

## Goal
`go-backend/docker-compose.yml` runs a `postgres:16` service with a named volume, so the service has a database to connect to.

## Concepts
Compose services and images, environment variables for the Postgres image, published ports, and named volumes.

## Check
`docker compose -f go-backend/docker-compose.yml up -d postgres` starts the container, `pg_isready` reports it accepts connections, and `/health` then returns 200.

## Result
2026-10-17: the `postgres-db` container (image `postgres:16`, project `go-backend`) accepted connections, and `/health` returned 200 with `{"status":"ok","database":"ok"}`.

## Redo checklist
1. Create `go-backend/docker-compose.yml` with one service, `postgres`, using the image `postgres:16` — keep it separate from the root Compose file.
2. Set `POSTGRES_USER`, `POSTGRES_PASSWORD` and `POSTGRES_DB` — the image creates the user and database on first start.
3. Publish port 5432 — the service connects through localhost.
4. Add a named volume mapped to `/var/lib/postgresql/data` — data survives container restarts.
5. Run `docker compose -f go-backend/docker-compose.yml up -d postgres` and wait for `pg_isready` — the container needs a moment to accept connections.
6. Put the matching `DATABASE_URL` in `go-backend/.env` — the service reads it at startup. A Compose healthcheck is GO-002.
