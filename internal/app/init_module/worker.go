package init_module

import (
	"crypto-watcher-backend/internal/app/worker"
	"fmt"
)

type (
	// WorkerWrapper wraps a Worker module with its gracful handler
	WorkerWrapper struct {
		Cron *worker.WatcherWorker
	}

	WorkerGracefulHandler struct{}
)

// NewWorkerWrapper creates a new Worker wrapper instance.
func NewWorkerWrapper(cron *worker.WatcherWorker) *WorkerWrapper {
	return &WorkerWrapper{
		Cron: cron,
	}
}

func NewWorkerGracefulHandler() *WorkerGracefulHandler {
	return &WorkerGracefulHandler{}
}

func (h *WorkerGracefulHandler) ReleaseResource() {
	fmt.Println("Releasing worker resources...")
	defer fmt.Println("Worker resources are gracefully released")
}
