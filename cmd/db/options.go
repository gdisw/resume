package db

type connectionOption func(*connectionConfiguration)

func WithConnectionsNumber(number int) connectionOption {
	return func(cfg *connectionConfiguration) {
		cfg.ConnectionsNumber = number
	}
}

func WithTimeout(timeout int) connectionOption {
	return func(cfg *connectionConfiguration) {
		cfg.Timeout = timeout
	}
}

func WithReadTimeout(timeout int) connectionOption {
	return func(cfg *connectionConfiguration) {
		cfg.ReadTimeout = timeout
	}
}

func WithVerboseLevel(level VerboseLevel) connectionOption {
	return func(cfg *connectionConfiguration) {
		cfg.VerboseLevel = level
	}
}

func WithDriver(driver Driver) connectionOption {
	return func(cfg *connectionConfiguration) {
		cfg.Driver = driver
	}
}

func WithDSN(dsn string) connectionOption {
	return func(cfg *connectionConfiguration) {
		cfg.DSN = dsn
	}
}
