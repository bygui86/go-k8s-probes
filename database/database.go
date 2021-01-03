package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ExpansiveWorlds/instrumentedsql"
	instrumentedsqlopentracing "github.com/ExpansiveWorlds/instrumentedsql/opentracing"
	"github.com/lib/pq"

	"github.com/bygui86/go-k8s-probes/logging"
)

const (
	// no tracing
	dbConnectionStringFormat = "host=%s port=%d user=%s password=%s dbname=%s sslmode=%s"
	dbDriverName             = "postgres"

	// with tracing
	instrumentedDbConnectionStringFormat = "postgres://%s:%s@%s:%d/%s?sslmode=%s"
	instrumentedDbDriverName             = "instrumented-" + dbDriverName
)

func New() (*sql.DB, error) {
	logging.Log.Info("Create new DB interface")

	cfg := loadConfig()

	db, dbErr := sql.Open(
		dbDriverName,
		fmt.Sprintf(dbConnectionStringFormat,
			cfg.dbHost, cfg.dbPort,
			cfg.dbUsername, cfg.dbPassword, cfg.dbName,
			cfg.dbSslMode,
		),
	)
	if dbErr != nil {
		return nil, dbErr
	}

	_, tableErr := db.Exec(createTableQuery)
	if tableErr != nil {
		return nil, tableErr
	}

	return db, nil
}

func NewWithWrappedTracing() (*sql.DB, error) {
	logging.Log.Info("Create new DB interface instrumented for tracing")

	cfg := loadConfig()

	// Get a database driver.Connector for a fixed configuration.
	connector, connErr := pq.NewConnector(fmt.Sprintf(instrumentedDbConnectionStringFormat,
		cfg.dbUsername, cfg.dbPassword,
		cfg.dbHost, cfg.dbPort,
		cfg.dbName, cfg.dbSslMode,
	))
	if connErr != nil {
		return nil, connErr
	}

	sql.Register(
		instrumentedDbDriverName,
		instrumentedsql.WrapDriver(
			connector.Driver(),
			instrumentedsql.WithTracer(instrumentedsqlopentracing.NewTracer()),
			instrumentedsql.WithLogger(
				instrumentedsql.LoggerFunc(func(ctx context.Context, msg string, keyvals ...interface{}) {
					logging.SugaredLog.Infof("%s %v", msg, keyvals)
				})),
		),
	)
	db, dbErr := sql.Open(
		instrumentedDbDriverName,
		fmt.Sprintf(dbConnectionStringFormat,
			cfg.dbHost, cfg.dbPort,
			cfg.dbUsername, cfg.dbPassword, cfg.dbName,
			cfg.dbSslMode,
		),
	)
	if dbErr != nil {
		return nil, dbErr
	}

	_, tableErr := db.Exec(createTableQuery)
	if tableErr != nil {
		return nil, tableErr
	}

	return db, nil
}
