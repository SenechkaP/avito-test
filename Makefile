include .env
export $(shell sed 's/=.*//' .env)

GOOSE_DRIVER = postgres
GOOSE_DBSTRING = postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB_NAME)
MIGRATIONS_DIR = migrations

.PHONY: migrate-up migrate-down migrate-status

migrate-up:
	@echo "Applying migrations from $(MIGRATIONS_DIR)..."
	GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING="$(GOOSE_DBSTRING)" goose -dir $(MIGRATIONS_DIR) up

migrate-down:
	@echo "Rolling back migrations..."
	GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING="$(GOOSE_DBSTRING)" goose -dir $(MIGRATIONS_DIR) down

migrate-status:
	@echo "Migration status in $(MIGRATIONS_DIR):"
	GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING="$(GOOSE_DBSTRING)" goose -dir $(MIGRATIONS_DIR) status