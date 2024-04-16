package init_module

import (
	"context"
	"crypto-watcher-backend/internal/config"
	"crypto-watcher-backend/pkg/database"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

func InitializeDB(dbConfig map[string]*config.Database) map[string]*database.Replication {
	const funcName = "[internal][app][init_module]InitializeDB"

	ctx := context.Background()
	replicaCollection := make(map[string]*database.Replication)
	for dbName, cfg := range dbConfig {
		if cfg == nil {
			logrus.Infof("%s: empty config for [%s]", funcName, dbName)
			continue
		}

		var (
			replica *database.Replication
			err     error
		)

		replica = &database.Replication{
			DriverName: cfg.Master.Driver,
		}

		replica.Master, err = Connect(ctx, cfg.Master.Driver, cfg.Master.Url, &database.ConnectionOptions{
			MaxIdleConnections:    cfg.Master.MaxIdle,
			MaxOpenConnections:    cfg.Master.MaxOpen,
			ConnectionMaxLifetime: time.Duration(cfg.Master.MaxLifeTime) * time.Second,
			Retry:                 1,
		})
		if err != nil {
			logrus.Fatalf("%s: Failed to connect to master DB [%s]", funcName, err)
		}

		if cfg.Slave != nil {
			replica.Slave, err = Connect(ctx, cfg.Master.Driver, cfg.Slave.Url, &database.ConnectionOptions{
				MaxIdleConnections:    cfg.Slave.MaxIdle,
				MaxOpenConnections:    cfg.Slave.MaxOpen,
				ConnectionMaxLifetime: time.Duration(cfg.Slave.MaxLifeTime) * time.Second,
				Retry:                 1,
			})
			if err != nil {
				logrus.Fatalf("%s: Failed to connect to slave DB [%s]", funcName, err)
			}
		}
		replicaCollection[dbName] = replica
	}
	logrus.Info("All database initiated")
	return replicaCollection
}

func connectWithRetry(ctx context.Context, driver, dataSource string, retry int) (*sqlx.DB, error) {
	var (
		db  *sqlx.DB
		err error
	)

	for c := 0; c < retry; c++ {
		db, err = sqlx.ConnectContext(ctx, driver, dataSource)
		if err != nil {
			time.Sleep(time.Second * 3)
		} else {
			break
		}
	}

	return db, err
}

func Connect(ctx context.Context, driver, dataSource string, conOpts *database.ConnectionOptions) (*sqlx.DB, error) {
	opts := conOpts
	if opts == nil {
		opts = &database.ConnectionOptions{}
	}

	db, err := connectWithRetry(ctx, driver, dataSource, opts.Retry)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(opts.MaxOpenConnections)
	db.SetMaxIdleConns(opts.MaxIdleConnections)
	db.SetConnMaxIdleTime(opts.ConnectionMaxLifetime)
	if e := db.Ping(); e != nil {
		logrus.Infof("%s: Error when ping [%s]", dataSource, err)
		panic(e)
	}

	return db, nil
}
