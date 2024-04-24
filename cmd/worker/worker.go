package worker

import (
	"crypto-watcher-backend/internal/app"
	"crypto-watcher-backend/internal/config"
	"net/http"

	"github.com/sirupsen/logrus"
)

func Start(cfg *config.Config) {
	logrus.Info("Starting crypto-watcher worker...")
	httpServer := http.Server{}
	app.RunWorker(&httpServer, cfg)
}
