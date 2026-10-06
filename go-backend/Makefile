# -----------------------------------------------------------------------------
# Configuration
# -----------------------------------------------------------------------------

export GOOSE_DRIVER := postgres
export GOOSE_DBSTRING := postgres://myapp:myapp_password@localhost:5432/myapp?sslmode=disable
export GOOSE_MIGRATION_DIR := ./internal/database/migrations
export GOOSE_TABLE := goose_db_version

# -----------------------------------------------------------------------------
# Default target
# -----------------------------------------------------------------------------

.PHONY: all
all: up

# -----------------------------------------------------------------------------
# Docker / PostgreSQL
# -----------------------------------------------------------------------------

.PHONY: postgres
postgres:
	@echo "Starting PostgreSQL..."
	docker compose up -d postgres

	@echo "Waiting for PostgreSQL to become healthy..."
	@until docker compose exec -T postgres pg_isready -U myapp -d myapp >/dev/null 2>&1; do \
		sleep 1; \
	done

	@echo "PostgreSQL is ready."

.PHONY: down
down:
	@echo "Stopping PostgreSQL..."
	docker compose down

.PHONY: restart
restart: down postgres

.PHONY: clean
clean:
	@echo "Stopping PostgreSQL and removing volumes..."
	docker compose down -v

# -----------------------------------------------------------------------------
# Database migrations
# -----------------------------------------------------------------------------

.PHONY: up
up: postgres
	@echo "Running database migrations..."
	goose -dir $(GOOSE_MIGRATION_DIR) up


.PHONY: down-migration
down-migration:
	@echo "Rolling back the last migration..."
	goose -dir $(GOOSE_MIGRATION_DIR) down

.PHONY: status
status:
	@goose -dir $(GOOSE_MIGRATION_DIR) status

# -----------------------------------------------------------------------------
# Development
# -----------------------------------------------------------------------------

.PHONY: logs
logs:
	docker compose logs -f postgres