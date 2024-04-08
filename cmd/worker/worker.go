package worker

import (
	"crypto-watcher-backend/internal/app"
	"crypto-watcher-backend/internal/config"
	"fmt"
	"net/http"
)

func Start(cfg *config.Config) {
	fmt.Println("Starting crypto-watcher worker...")
	httpServer := http.Server{}
	app.RunWorker(&httpServer, cfg)
}
