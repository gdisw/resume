package db

import (
	"os"
	"strings"
)

const (
	VerboseLevelMinimum = VerboseLevel(iota)
	VerboseLevelErrors
	VerboseLevelAll
)

type VerboseLevel uint8

type Driver string

const (
	DriverPostgres Driver = "postgres"
	DriverSQLite   Driver = "sqlite"
)

func inferredDriverFromDSN(dsn string) Driver {
	if strings.HasPrefix(dsn, "postgres://") || strings.HasPrefix(dsn, "postgresql://") {
		return DriverPostgres
	}
	return DriverSQLite
}

type connectionConfiguration struct {
	Driver            Driver
	DSN               string
	ConnectionsNumber int
	VerboseLevel      VerboseLevel
	Timeout           int
	ReadTimeout       int
}

func getDefaultConfiguration() connectionConfiguration {
	return connectionConfiguration{
		DSN:               os.Getenv("DATABASE_DSN"),
		ConnectionsNumber: 2,
		VerboseLevel:      VerboseLevelErrors,
		Timeout:           0,
	}
}
