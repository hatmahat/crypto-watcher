package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type (
	// Config holds all the configuration for the application.
	Config struct {
		ENV string
		DB  map[string]*Database
		ServerConfig
		WorkerConfig
		SchedulerConfig
		CoinConfig
		CoinGeckoConfig
		CurrencyConfig
		WhatsAppConfig
	}

	ServerConfig struct {
		APIPort       int
		GlobalTimeout int
		APILogLevel   string
	}

	WorkerConfig struct {
		APIPort       int
		GlobalTimeout int
		APILogLevel   string
	}

	SchedulerConfig struct {
		SchedulerBitCoinFetch string
	}

	CoinConfig struct {
		CoinAPIHost string
		CoinAPIKey  string
	}

	CoinGeckoConfig struct {
		CoinGeckoAPIHost string
	}

	CurrencyConfig struct {
		CurrencyAPIHost string
		CurrencyAPIKey  string
	}

	WhatsAppConfig struct {
		WhatsAppAPIHost         string
		WhatsAppAPIKey          string
		WhatsAppPhoneNumberId   string
		WhatsAppTestPhoneNumber string
	}
)

// LoadConfig loads the application configuration from a given path and name.
func LoadConfig(configPath, fileName string) (*Config, error) {
	viper.AddConfigPath(configPath)
	viper.SetConfigName(fileName)
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Error reading config file: %s\n", err)
		log.Println("Falling back to environment variables only.")
	}

	env := getStringOrDefault("ENV", "staging")

	dbMasterConnection := DatabaseConnectionConfig{
		Host:     getStringOrPanic("DATABASE_MASTER_HOST"),
		Port:     getStringOrPanic("DATABASE_MASTER_PORT"),
		User:     getStringOrPanic("DATABASE_MASTER_USER"),
		Password: getStringOrPanic("DATABASE_MASTER_PASSWORD"),
		DBName:   getStringOrPanic("DATABASE_MASTER_DBNAME"),
		SSLMode:  getStringOrPanic("DATABASE_MASTER_SSLMODE"),
	}

	dbMasterUrl := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbMasterConnection.Host, dbMasterConnection.Port, dbMasterConnection.User,
		dbMasterConnection.Password, dbMasterConnection.DBName, dbMasterConnection.SSLMode,
	)

	dbMaster := &DatabaseConfig{
		Driver:      getStringOrPanic("DATABASE_MASTER_DRIVER"),
		Url:         dbMasterUrl,
		MaxOpen:     getIntOrDefault("DATABASE_MASTER_MAX_OPEN_CONNECTION", 10),
		MaxIdle:     getIntOrDefault("DATABASE_MASTER_MAX_IDLE", 5),
		MaxLifeTime: getIntOrDefault("DATABASE_MASTER_CONNECTION_MAX_LIFE_TIME", 0),
	}

	db := map[string]*Database{
		CryptoWatcherDB: {
			Master: dbMaster,
			Slave:  nil, // TODO: add slave DB config here
		},
	}

	serverConfig := ServerConfig{
		APIPort:       getIntOrDefault("SERVER_API_PORT", 9000),
		GlobalTimeout: getIntOrDefault("SERVER_GLOBAL_TIMEOUT", 30000),
		APILogLevel:   getStringOrDefault("SERVER_API_LOG_LEVEL", "info"),
	}

	workerConfig := WorkerConfig{
		APIPort:       getIntOrDefault("WORKER_API_PORT", 8081),
		GlobalTimeout: getIntOrDefault("WORKER_GLOBAL_TIMEOUT", 30000),
		APILogLevel:   getStringOrDefault("WORKER_API_LOG_LEVEL", "info"),
	}

	schedulerConfig := SchedulerConfig{
		SchedulerBitCoinFetch: getStringOrPanic("SCHEDULER_BITCOIN_FETCH"),
	}

	coinConfig := CoinConfig{
		CoinAPIHost: getStringOrPanic("COIN_API_HOST"),
		CoinAPIKey:  getStringOrPanic("COIN_API_KEY"),
	}

	coinGeckoConfig := CoinGeckoConfig{
		CoinGeckoAPIHost: getStringOrPanic("COIN_GECKO_API_HOST"),
	}

	whatsAppConfig := WhatsAppConfig{
		WhatsAppAPIHost:         getStringOrPanic("WHATSAPP_API_HOST"),
		WhatsAppAPIKey:          getStringOrPanic("WHATSAPP_API_KEY"),
		WhatsAppPhoneNumberId:   getStringOrPanic("WHATSAPP_PHONENUMBER_ID"),
		WhatsAppTestPhoneNumber: getStringOrDefault("WHATSAPP_TEST_PHONENUMBER", ""),
	}

	currencyConfig := CurrencyConfig{
		CurrencyAPIHost: getStringOrPanic("CURRENCY_API_HOST"),
		CurrencyAPIKey:  getStringOrPanic("CURRENCY_API_KEY"),
	}

	return &Config{
		ENV:             env,
		DB:              db,
		ServerConfig:    serverConfig,
		WorkerConfig:    workerConfig,
		SchedulerConfig: schedulerConfig,
		CoinConfig:      coinConfig,
		CoinGeckoConfig: coinGeckoConfig,
		WhatsAppConfig:  whatsAppConfig,
		CurrencyConfig:  currencyConfig,
	}, nil
}

func getIntOrDefault(key string, defaultValue int) int {
	viper.SetDefault(key, defaultValue)
	return viper.GetInt(key)
}

func getStringOrDefault(key, defaultValue string) string {
	viper.SetDefault(key, defaultValue)
	return viper.GetString(key)
}

func getStringOrPanic(key string) string {
	val := viper.GetString(key)
	if val == "" {
		panic("No value found for key: " + key)
	}
	return val
}
