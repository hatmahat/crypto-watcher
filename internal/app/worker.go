package app

import (
	"context"
	"crypto-watcher-backend/internal/app/init_module"
	"crypto-watcher-backend/internal/config"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

func RunWorker(httpServer *http.Server, cfg *config.Config) {
	ctx := context.Background()
	httpClient := http.Client{
		Timeout: time.Duration(cfg.WorkerConfig.GlobalTimeout) * time.Millisecond,
	}

	worker := init_module.NewWorker(ctx, cfg, &httpClient)
	worker.Cron.Start()
	defer worker.Cron.Stop()

	router := mux.NewRouter()
	httpServer.Addr = fmt.Sprintf(":%d", cfg.WorkerConfig.APIPort)
	httpServer.Handler = router

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Print("Server Started")

	<-done
	log.Print("Worker Stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	err := httpServer.Shutdown(ctx)
	if err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Print("Worker Exited")
}
