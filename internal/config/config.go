package config

import (
	"log"

	"github.com/spf13/viper"
)

type (
	// Config holds all the configuration for the application.
	Config struct {
		ENV string
		ServerConfig
		WorkerConfig
		SchedulerConfig
		APIURLsConfig
		APIKeysConfig
		CredentialsConfig
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

	APIURLsConfig struct {
		APIURLCoinGecko string
	}

	APIKeysConfig struct {
		APIKeyWhatsApp string
	}

	CredentialsConfig struct {
		CredWhatsAppPhoneNumberId string
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

	env := getStringConfigOrDefault("ENV", "staging")

	serverConfig := ServerConfig{
		APIPort:       getIntConfigOrDefault("SERVER_API_PORT", 8080),
		GlobalTimeout: getIntConfigOrDefault("SERVER_GLOBAL_TIMEOUT", 30000),
		APILogLevel:   getStringConfigOrDefault("SERVER_API_LOG_LEVEL", "info"),
	}

	workerConfig := WorkerConfig{
		APIPort:       getIntConfigOrDefault("WORKER_API_PORT", 8081),
		GlobalTimeout: getIntConfigOrDefault("WORKER_GLOBAL_TIMEOUT", 30000),
		APILogLevel:   getStringConfigOrDefault("WORKER_API_LOG_LEVEL", "info"),
	}

	schedulerConfig := SchedulerConfig{
		SchedulerBitCoinFetch: getStringOrPanic("SCHEDULER_BITCOIN_FETCH"),
	}

	APIURLs := APIURLsConfig{
		APIURLCoinGecko: getStringOrPanic("API_URL_COIN_GECKO"),
	}

	APIKeys := APIKeysConfig{
		APIKeyWhatsApp: getStringOrPanic("API_KEY_WHATSAPP"),
	}

	credsConfig := CredentialsConfig{
		CredWhatsAppPhoneNumberId: getStringOrPanic("CRED_WHATSAPP_PHONENUMBER_ID"),
	}

	return &Config{
		ENV:               env,
		ServerConfig:      serverConfig,
		WorkerConfig:      workerConfig,
		SchedulerConfig:   schedulerConfig,
		APIURLsConfig:     APIURLs,
		APIKeysConfig:     APIKeys,
		CredentialsConfig: credsConfig,
	}, nil
}

func getIntConfigOrDefault(key string, defaultValue int) int {
	viper.SetDefault(key, defaultValue)
	return viper.GetInt(key)
}

func getStringConfigOrDefault(key, defaultValue string) string {
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
