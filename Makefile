.PHONY: gen gen-upd dev db test
include apps/server/.env
export

gen:
	cd apps/proto && buf generate --include-imports
	cd apps/server && sqlc generate && go mod tidy

gen-upd:
	cd apps/proto && buf dep update && buf generate
	cd apps/server && go mod tidy

dev:
	make dev-server & make dev-web

dev-server:
	cd apps/server && go run cmd/server/main.go

dev-web:
	cd apps/web && bun run dev

test:
	cd apps/server && go run cmd/test/main.go

db:
	docker compose up -d postgres

migrate:
	cd apps/server && goose -dir postgres/migrations postgres "${DATABASE_URL}" up

migrate-down:
	cd apps/server && goose -dir postgres/migrations postgres "${DATABASE_URL}" down
