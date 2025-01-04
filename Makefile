start-watcher: build run-watcher

start-nohup-watcher: build run-nohup-watcher

build:
	@echo ">> Building crypto-watcher..."
	@go build --race -o ./bin/crypto-watcher ./cmd
	@echo ">> Finished"

run-watcher:
	@./bin/crypto-watcher watcher

run-nohup-watcher:
	@nohup ./bin/crypto-watcher watcher > app.log 2>&1 &
	@echo "Watcher started in the background. Logs: app.log"

stop-nohup-watcher:
	@pkill -f './bin/crypto-watcher watcher' || echo "No watcher process found"

watch-log:
	@tail -f app.log

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