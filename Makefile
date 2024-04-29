start-watcher: build run-watcher

build:
	@echo ">> Building crypto-watcher..."
	@go build --race -o ./bin/crypto-watcher ./cmd
	@echo ">> Finished"

run-watcher:
	@./bin/crypto-watcher watcher

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