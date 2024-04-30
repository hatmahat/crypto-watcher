//go:build wireinject
// +build wireinject

package init_module

import (
	"context"
	"crypto-watcher-backend/internal/app/worker"
	"crypto-watcher-backend/internal/config"
	"crypto-watcher-backend/internal/repository"
	"crypto-watcher-backend/internal/service"
	"net/http"

	"github.com/google/wire"
)

var (
	cfgSet = wire.NewSet(
		wire.FieldsOf(new(*config.Config), "ServerConfig"),
		wire.FieldsOf(new(*config.Config), "DB"),
		wire.FieldsOf(new(*config.Config), "WorkerConfig"),
		wire.FieldsOf(new(*config.Config), "SchedulerConfig"),
		wire.FieldsOf(new(*config.Config), "CoinGeckoConfig"),
		wire.FieldsOf(new(*config.Config), "WhatsAppConfig"),
	)

	dependencySet = wire.NewSet(
		InitializeDB,
		NewCoin,
		NewCoinGecko,
		NewCurrency,
		NewCurrencyConverter,
		NewWaMessaging,
		NewTelegramBot,
	)

	repoSet = wire.NewSet(
		repository.NewCurrencyRateRepo,
		wire.Struct(new(repository.CurrencyRateRepoParam), "*"),
		repository.NewAssetPriceRepo,
		wire.Struct(new(repository.AssetPriceRepoParam), "*"),
		repository.NewUserRepo,
		wire.Struct(new(repository.UserRepoParam), "*"),
		repository.NewNotificationRepo,
		wire.Struct(new(repository.NotificationRepoParam), "*"),
		repository.NewUserPreferenceRepo,
		wire.Struct(new(repository.UserPreferenceRepoParam), "*"),
	)

	serviceSet = wire.NewSet(
		service.NewCryptoService,
		wire.Struct(new(service.CryptoServiceParam), "*"),
	)

	appSet = wire.NewSet(
		wire.Struct(new(worker.WatcherWorkerParam), "*"),
		worker.NewWatcherWorker,
		NewWorkerWrapper,
		NewWorkerGracefulHandler,
	)

	allSet = wire.NewSet(
		cfgSet,
		dependencySet,
		repoSet,
		serviceSet,
		appSet,
	)
)

func NewWorker(ctx context.Context, cfg *config.Config, httpClient *http.Client) *WorkerWrapper {
	wire.Build(allSet)
	return &WorkerWrapper{}
}
