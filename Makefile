start-worker: build run-worker

build:
	@echo ">> Building crypto-watcher..."
	@go build --race -o ./bin/crypto-watcher ./cmd
	@echo ">> Finished"

run-worker:
	@./bin/crypto-watcher worker

wire:
	@cd internal/app/init_module && go run github.com/google/wire/cmd/wire

start-db:
	@echo "Starting the PostgreSQL container..."
	@docker-compose up -d db

deploy-db:
	@cd migration && sqitch deploy 

revert-db:
	@cd migration && sqitch revert

verify-db:
	@cd migration && sqitch verify