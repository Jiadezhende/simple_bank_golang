# =========================================================
# Load environment variables
# =========================================================
ifneq (,$(wildcard .env))
	include .env
	export
endif

# =========================================================
# Derived variables
# =========================================================
DB_URL := postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

MIGRATE        := migrate
MIGRATION_DIR := db/migration

# =========================================================
# Docker: Postgres
# =========================================================
postgres:
	@docker ps -a --format '{{.Names}}' | grep -q "^$(POSTGRES_CONTAINER)$$" || \
	docker run --name $(POSTGRES_CONTAINER) \
		-p $(DB_PORT):5432 \
		-e POSTGRES_USER=$(DB_USER) \
		-e POSTGRES_PASSWORD=$(DB_PASSWORD) \
		-d $(POSTGRES_IMAGE)

postgres-stop:
	@docker stop $(POSTGRES_CONTAINER)

postgres-rm:
	@docker rm -f $(POSTGRES_CONTAINER)

# =========================================================
# Database lifecycle (idempotent)
# =========================================================
createdb:
	@docker exec -i $(POSTGRES_CONTAINER) \
		psql -U $(DB_USER) postgres \
		-c "CREATE DATABASE $(DB_NAME);" \
	|| echo "database $(DB_NAME) already exists"

dropdb:
	@docker exec -i $(POSTGRES_CONTAINER) \
		psql -U $(DB_USER) postgres \
		-c "DROP DATABASE $(DB_NAME);" \
	|| echo "database $(DB_NAME) does not exist"

# =========================================================
# Migrations
# =========================================================
migrate-create:
	@test $(name) || (echo "❌ usage: make migrate-create name=xxx" && exit 1)
	@$(MIGRATE) create -ext sql -dir $(MIGRATION_DIR) -seq $(name)

migrate-up:
	@$(MIGRATE) -path $(MIGRATION_DIR) -database "$(DB_URL)" up

migrate-down:
	@$(MIGRATE) -path $(MIGRATION_DIR) -database "$(DB_URL)" down 1

migrate-force:
	@test $(version) || (echo "❌ usage: make migrate-force version=xxx" && exit 1)
	@$(MIGRATE) -path $(MIGRATION_DIR) -database "$(DB_URL)" force $(version)

# =========================================================
# SQLC: Generate codes for SQL in Go
# =========================================================
sqlc:
	sqlc generate

# =========================================================
# Phony targets
# =========================================================
.PHONY: \
	postgres postgres-stop postgres-rm \
	createdb dropdb \
	migrate-create migrate-up migrate-down migrate-force \
	sqlc
