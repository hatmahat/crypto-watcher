package worker

import (
	"context"
	"crypto-watcher-backend/internal/config"
	"crypto-watcher-backend/internal/service"

	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

type (
	WatcherWorker struct {
		cronPool *cron.Cron
		config   *config.Config
		cronJob  []worker
	}

	WatcherWorkerParam struct {
		Config        *config.Config
		CryptoService service.CryptoService
	}

	JobParameter struct {
		Name     string
		TimeSpec string
		Handler  func(ctx context.Context)
	}

	worker interface {
		GenerateWorkerParameters() []JobParameter
	}
)

func NewWatcherWorker(param WatcherWorkerParam) *WatcherWorker {
	cronPool := cron.New(
		cron.WithParser(
			cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor),
		),
	)

	c := WatcherWorker{
		cronPool: cronPool,
		config:   param.Config,
	}

	c.cronJob = []worker{
		NewCryptoWorker(CryptoWorkerParam{
			WatcherWorker: &c,
			CryptoService: param.CryptoService,
		}),
	}

	c.registerJobs()

	return &c
}

func (c *WatcherWorker) registerJobs() {
	logrus.Info("Registering jobs...")
	defer logrus.Info("Jobs are successfully registered")

	for i := range c.cronJob {
		for _, cronJob := range c.cronJob[i].GenerateWorkerParameters() {
			_, err := c.cronPool.AddFunc(cronJob.TimeSpec, c.GenerateHandler(cronJob.Handler))
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"jobName":  cronJob.Name,
					"timeSpec": cronJob.TimeSpec,
					"err":      err.Error(),
				}).Errorf("Error Registering Job: %s", cronJob.Name)
			}
		}
	}
}

func (c *WatcherWorker) GenerateHandler(f func(ctx context.Context)) func() {
	return func() {
		f(context.Background())
	}
}

// Start starts all the registered jobs
func (c *WatcherWorker) Start() {
	logrus.Info("Starting worker...")
	c.cronPool.Start()
}

// Stop stops all currently running jobs
func (c *WatcherWorker) Stop() {
	logrus.Info("Stopping worker...")
	defer logrus.Info("Cron is gracefully stopped")
	c.cronPool.Stop()
}
