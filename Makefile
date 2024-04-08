start-worker: build run-worker

build:
	@echo ">> Building crypto-watcher..."
	@go build --race -o crypto-watcher ./cmd
	@echo ">> Finished"

run-worker:
	@./crypto-watcher worker

wire:
	@cd internal/app/init_module && go run github.com/google/wire/cmd/wire