package config

const CryptoWatcherDB = "crypto_watcher_db"

type (
	Database struct {
		Master *DatabaseConfig
		Slave  *DatabaseConfig
	}

	DatabaseConfig struct {
		Driver      string
		Url         string
		MaxIdle     int
		MaxOpen     int
		MaxLifeTime int
	}

	DatabaseConnectionConfig struct {
		Host     string
		Port     string
		User     string
		Password string
		DBName   string
		SSLMode  string
	}
)
