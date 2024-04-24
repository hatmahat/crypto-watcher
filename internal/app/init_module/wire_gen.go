// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package init_module

import (
	"context"
	"crypto-watcher-backend/internal/app/worker"
	"crypto-watcher-backend/internal/config"
	"crypto-watcher-backend/internal/repository"
	"crypto-watcher-backend/internal/service"
	"github.com/google/wire"
	"net/http"
)

// Injectors from wire.go:

func NewWorker(ctx context.Context, cfg *config.Config, httpClient *http.Client) *WorkerWrapper {
	coinGecko := NewCoinGecko(cfg)
	coin := NewCoin(cfg)
	currencyConverter := NewCurrencyConverter(cfg)
	waMessaging := NewWaMessaging(cfg)
	telegramBot := NewTelegramBot(cfg)
	v := cfg.DB
	v2 := InitializeDB(v)
	currencyRateRepoParam := repository.CurrencyRateRepoParam{
		DB: v2,
	}
	currencyRateRepo := repository.NewCurrencyRateRepo(currencyRateRepoParam)
	assetPriceRepoParam := repository.AssetPriceRepoParam{
		DB: v2,
	}
	assetPriceRepo := repository.NewAssetPriceRepo(assetPriceRepoParam)
	cryptoServiceParam := service.CryptoServiceParam{
		CoinGecko:         coinGecko,
		Coin:              coin,
		CurrencyConverter: currencyConverter,
		WaMessaging:       waMessaging,
		TelegramBot:       telegramBot,
		Cfg:               cfg,
		CurrencyRateRepo:  currencyRateRepo,
		AssetPriceRepo:    assetPriceRepo,
	}
	cryptoService := service.NewCryptoService(cryptoServiceParam)
	watcherWorkerParam := worker.WatcherWorkerParam{
		Config:        cfg,
		CryptoService: cryptoService,
	}
	watcherWorker := worker.NewWatcherWorker(watcherWorkerParam)
	workerWrapper := NewWorkerWrapper(watcherWorker)
	return workerWrapper
}

// wire.go:

var (
	cfgSet = wire.NewSet(wire.FieldsOf(new(*config.Config), "ServerConfig"), wire.FieldsOf(new(*config.Config), "DB"), wire.FieldsOf(new(*config.Config), "WorkerConfig"), wire.FieldsOf(new(*config.Config), "SchedulerConfig"), wire.FieldsOf(new(*config.Config), "CoinGeckoConfig"), wire.FieldsOf(new(*config.Config), "WhatsAppConfig"))

	dependencySet = wire.NewSet(
		InitializeDB,
		NewCoin,
		NewCoinGecko,
		NewCurrency,
		NewCurrencyConverter,
		NewWaMessaging,
		NewTelegramBot,
	)

	repoSet = wire.NewSet(repository.NewCurrencyRateRepo, wire.Struct(new(repository.CurrencyRateRepoParam), "*"), repository.NewAssetPriceRepo, wire.Struct(new(repository.AssetPriceRepoParam), "*"))

	serviceSet = wire.NewSet(service.NewCryptoService, wire.Struct(new(service.CryptoServiceParam), "*"))

	appSet = wire.NewSet(wire.Struct(new(worker.WatcherWorkerParam), "*"), worker.NewWatcherWorker, NewWorkerWrapper,
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
