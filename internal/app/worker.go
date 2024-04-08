package app

import (
	"crypto-watcher-backend/internal/config"
	"fmt"
	"net/http"
)

func RunWorker(httpServer *http.Server, cfg *config.Config) {
	// ctx := context.Background()
	// httpClient := http.Client{
	// 	Timeout: time.Duration(cfg.WorkerConfig.GlobalTimeout) * time.Millisecond,
	// }

	// TODO: Run Worker
	fmt.Println("RUN WORKER")
}
