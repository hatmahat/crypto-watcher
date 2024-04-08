package worker

import (
	"crypto-watcher-backend/internal/app"
	"crypto-watcher-backend/internal/config"
	"fmt"
)

func Start(cfg *config.Config) {
	fmt.Println("Starting crypto-watcher worker...")
	app.RunWorker(cfg)
}
