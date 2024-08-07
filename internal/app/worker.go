package app

import (
	"context"
	"crypto-watcher-backend/internal/app/init_module"
	"crypto-watcher-backend/internal/config"
	"crypto-watcher-backend/pkg/logger"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/labstack/echo"
)

func RunWorker(httpServer *http.Server, cfg *config.Config) {
	ctx := context.Background()
	httpClient := http.Client{
		Timeout: time.Duration(cfg.WorkerConfig.GlobalTimeout) * time.Millisecond,
	}

	worker := init_module.NewWorker(ctx, cfg, &httpClient)
	worker.Cron.Start()
	defer worker.Cron.Stop()

	echoServer := echo.New()
	echoServer.Use(logger.MiddlewareLogging)

	httpServer.Addr = fmt.Sprintf(":%d", cfg.WorkerConfig.APIPort)
	httpServer.Handler = echoServer

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("listen: %s\n", err)
		}
	}()
	logger.LogInfoWithCustomTime("Server Started")

	<-done
	logger.LogInfoWithCustomTime("Worker Stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	err := httpServer.Shutdown(ctx)
	if err != nil {
		logrus.Error(fmt.Sprintf("Server Shutdown Failed: %+v", err))
	}
	logger.LogInfoWithCustomTime("Worker Exited")
}
