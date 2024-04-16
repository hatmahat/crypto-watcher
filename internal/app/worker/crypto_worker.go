package worker

import (
	"context"
	"crypto-watcher-backend/internal/constant/worker_const"
	"crypto-watcher-backend/internal/service"
	"crypto-watcher-backend/pkg/logger"
	"fmt"

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
			Name:     worker_const.BitcoinWatcher,
			TimeSpec: cw.config.SchedulerBitCoinFetch,
			Handler:  cw.BitcoinWatcher,
		},
	}
}

func (cw *CryptoWorker) BitcoinWatcher(ctx context.Context) {
	const funcName = "[internal][app][worker]BitcoinWatcher"
	logger.LogWithCustomTime(fmt.Sprintf("%s: Running", funcName))
	err := cw.cryptoService.BitcoinWatcher(ctx)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Errorf("%s: Error", funcName)
	}
}
