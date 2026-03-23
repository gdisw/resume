package db

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/driver/sqliteshim"
	"github.com/uptrace/bun/extra/bundebug"
	"github.com/uptrace/bun/schema"
)

func OpenConnection(opts ...connectionOption) (*bun.DB, error) {
	cfg := getDefaultConfiguration()

	for _, opt := range opts {
		opt(&cfg)
	}

	if cfg.Driver == "" {
		cfg.Driver = inferredDriverFromDSN(cfg.DSN)
	}

	var sqldb *sql.DB
	var dialect schema.Dialect

	switch cfg.Driver {
	case DriverSQLite:
		dsn := cfg.DSN
		if cfg.Timeout > 0 {
			sep := "?"
			if strings.Contains(dsn, "?") {
				sep = "&"
			}
			// Map cfg.Timeout (seconds) to busy_timeout (milliseconds)
			// modernc.org/sqlite uses _pragma=busy_timeout(ms)
			dsn = fmt.Sprintf("%s%s_pragma=busy_timeout(%d)", dsn, sep, cfg.Timeout*1000)
		}

		var err error
		sqldb, err = sql.Open(sqliteshim.ShimName, dsn)
		if err != nil {
			return nil, err
		}

		dialect = sqlitedialect.New()

	case DriverPostgres:
		dOpts := []pgdriver.Option{
			pgdriver.WithDSN(cfg.DSN),
		}

		if cfg.Timeout > 0 {
			dOpts = append(dOpts, pgdriver.WithTimeout(
				time.Duration(cfg.Timeout)*time.Second,
			))
		}

		if cfg.ReadTimeout > 0 {
			dOpts = append(dOpts, pgdriver.WithReadTimeout(
				time.Duration(cfg.ReadTimeout)*time.Second,
			))
		}

		sqldb = sql.OpenDB(pgdriver.NewConnector(dOpts...))
		dialect = pgdialect.New()
	default:
		return nil, fmt.Errorf("unsupported driver: %s", cfg.Driver)
	}

	sqldb.SetMaxOpenConns(cfg.ConnectionsNumber)
	sqldb.SetMaxIdleConns(cfg.ConnectionsNumber)
	db := bun.NewDB(sqldb, dialect)

	switch cfg.VerboseLevel {
	case VerboseLevelErrors:
		db.AddQueryHook(bundebug.NewQueryHook())
	case VerboseLevelAll:
		db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	}

	return db, nil
}
