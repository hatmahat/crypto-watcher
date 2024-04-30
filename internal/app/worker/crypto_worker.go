package worker

import (
	"context"
	"crypto-watcher-backend/internal/config"
	"crypto-watcher-backend/internal/constant/worker_const"
	"crypto-watcher-backend/internal/service"
	"crypto-watcher-backend/pkg/logger"

	"github.com/sirupsen/logrus"
)

type (
	CryptoWorker struct {
		*WatcherWorker
		cryptoService service.CryptoService
	}

	CryptoWorkerParam struct {
		*WatcherWorker
		CryptoService service.CryptoService
	}
)

func NewCryptoWorker(param CryptoWorkerParam) *CryptoWorker {
	return &CryptoWorker{
		WatcherWorker: param.WatcherWorker,
		cryptoService: param.CryptoService,
	}
}

func (cw *CryptoWorker) GenerateWorkerParameters() []JobParameter {
	return []JobParameter{
		{
			Name:     worker_const.CryptoWatcher,
			TimeSpec: cw.config.SchedulerCryptoFetch,
			Handler:  cw.CryptoWatcher,
		},
	}
}

func (cw *CryptoWorker) CryptoWatcher(ctx context.Context) {
	const funcName = "[internal][app][worker]CryptoWatcher"
	if config.DebugMode {
		logger.LogInfoWithCustomTime(">>Running<< Crypto Watcher Service")
	}
	err := cw.cryptoService.CryptoWatcher(ctx)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Errorf("%s: Error", funcName)
	}
	if config.DebugMode {
		logger.LogInfoWithCustomTime("<<Stopped>> Crypto Watcher Service ")
	}
}
