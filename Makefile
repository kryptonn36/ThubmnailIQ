.PHONY: infra-up infra-down migrate-up migrate-down migrate-status sqlc-generate api worker web build test dev

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

build:
	go build -o bin/api ./cmd/api
	go build -o bin/worker ./cmd/worker

test:
	go test ./... -race -count=1

# Runs infra + migrations, then prints the three commands to run in
# separate terminals (api/worker/web need their own foreground process).
dev: infra-up migrate-up
	@echo ""
	@echo "Infra is up and migrated. Now run each of these in its own terminal:"
	@echo "  make api"
	@echo "  make worker"
	@echo "  make web"
