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
			Name:     worker_const.BitcoinPriceWatcher,
			TimeSpec: cw.config.SchedulerBitCoinFetch,
			Handler:  cw.BitcoinPriceWatcher,
		},
	}
}

func (cw *CryptoWorker) BitcoinPriceWatcher(ctx context.Context) {
	const funcName = "[internal][app][worker]BitcoinPriceWatcher"
	logger.LogWithCustomTime(fmt.Sprintf("%s: Running", funcName))
	err := cw.cryptoService.BitcoinPriceWatcher(ctx)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Errorf("%s: Error", funcName)
	}
}
