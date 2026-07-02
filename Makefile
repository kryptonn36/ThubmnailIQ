.PHONY: infra-up infra-down migrate-up migrate-down migrate-status sqlc-generate api worker web admin-web admin-seed build test dev

-include .env
export

infra-up:
	docker compose up -d postgres redis minio cv-service

infra-down:
	docker compose down

migrate-up:
	goose -dir db/migrations postgres "$(DATABASE_URL)" up

migrate-down:
	goose -dir db/migrations postgres "$(DATABASE_URL)" down

migrate-status:
	goose -dir db/migrations postgres "$(DATABASE_URL)" status

sqlc-generate:
	cd db && sqlc generate

api: infra-up
	docker compose up --build --no-deps api

worker: infra-up
	docker compose up --build --no-deps worker

web:
	cd web && npm run dev

admin-web:
	cd admin-web && npm run dev

# One-time bootstrap for the first admin account — set ADMIN_SEED_EMAIL /
# ADMIN_SEED_PASSWORD in .env first. Safe to re-run (no-op if it exists).
admin-seed:
	go run ./cmd/admin-seed

build:
	go build -o bin/api ./cmd/api
	go build -o bin/worker ./cmd/worker
	go build -o bin/admin-seed ./cmd/admin-seed

test:
	go test ./... -race -count=1

# Runs infra + migrations, then prints the commands to run in separate
# terminals (each needs its own foreground process).
dev: infra-up migrate-up
	@echo ""
	@echo "Infra is up and migrated. Now run each of these in its own terminal:"
	@echo "  make api"
	@echo "  make worker"
	@echo "  make web"
	@echo "  make admin-web"
	@echo ""
	@echo "First time only, to bootstrap an admin account:"
	@echo "  make admin-seed"
