package worker

import (
	"context"
	"crypto-watcher-backend/internal/config"
	"fmt"

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
		Config *config.Config
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
		// TODO: Add new worker here
	}

	c.registerJobs()

	return &c
}

func (c *WatcherWorker) registerJobs() {
	fmt.Println("Registering jobs...")
	defer fmt.Println("Jobs are successfully registered")

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
	fmt.Println("Starting worker...")
	c.cronPool.Start()
}

// Stop stops all currently running jobs
func (c *WatcherWorker) Stop() {
	fmt.Println("Stopping worker...")
	defer fmt.Println("Cron is gracefully stopped")
	c.cronPool.Stop()
}
