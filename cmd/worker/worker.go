package worker

import (
	"crypto-watcher-backend/internal/app"
	"crypto-watcher-backend/internal/config"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

func Start(cfg *config.Config) {
	logrus.Info("Starting crypto-watcher worker...")
	httpServer := http.Server{
		ReadTimeout:  time.Duration(cfg.WorkerConfig.GlobalTimeout) * time.Microsecond,
		WriteTimeout: time.Duration(cfg.WorkerConfig.GlobalTimeout) * time.Microsecond,
		IdleTimeout:  time.Duration(cfg.WorkerConfig.GlobalTimeout) * time.Microsecond,
	}
	app.RunWorker(&httpServer, cfg)
}
