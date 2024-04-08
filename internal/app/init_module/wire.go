//go:build wireinject
// +build wireinject

package init_module

import (
	"context"
	"crypto-watcher-backend/internal/app/worker"
	"crypto-watcher-backend/internal/config"
	"net/http"

	"github.com/google/wire"
)

var (
	cfgSet = wire.NewSet(
		wire.FieldsOf(new(*config.Config), "ServerConfig"),
		wire.FieldsOf(new(*config.Config), "WorkerConfig"),
		wire.FieldsOf(new(*config.Config), "SchedulerConfig"),
		wire.FieldsOf(new(*config.Config), "CoinGeckoConfig"),
		wire.FieldsOf(new(*config.Config), "WhatsAppConfig"),
	)

	appSet = wire.NewSet(
		wire.Struct(new(worker.WatcherWorkerParam), "*"),
		worker.NewWatcherWorker,
		NewWorkerWrapper,
		NewWorkerGracefulHandler,
	)

	allSet = wire.NewSet(
		appSet,
		cfgSet,
	)
)

func NewWorker(ctx context.Context, cfg *config.Config, httpClient *http.Client) *WorkerWrapper {
	wire.Build(allSet)
	return &WorkerWrapper{}
}
