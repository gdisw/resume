package env

import (
	"errors"
	"os"
	"strings"

	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
)

var ErrUndefined = errors.New("env: undefined environment")

type environment uint8

const (
	Undefined environment = iota
	Local
	Development
	Staging
	Production
)

func (e environment) String() string {
	switch e {
	case Local:
		return "local"
	case Development:
		return "development"
	case Staging:
		return "staging"
	case Production:
		return "production"
	default:
		return "undefined"
	}
}

func Parse(s string) (environment, error) {
	switch strings.ToLower(s) {
	case Development.String():
		return Development, nil
	case Staging.String():
		return Staging, nil
	case Production.String():
		return Production, nil
	case Local.String():
		fallthrough
	default:
		return Local, nil
	}
}

var currentEnvironment environment

func Load(opts ...loadOption) error {
	var cfg loadCfg
	for _, opt := range opts {
		opt(&cfg)
	}

	if cfg.ForceLoad() {
		if err := godotenv.Load(cfg.defaultEnvFilePath()); err != nil {
			return err
		}
	}

	env, err := Parse(os.Getenv("ENVIRONMENT"))
	if err != nil {
		return err
	}

	currentEnvironment = env

	envFile := cfg.envFilePath(env)
	if _, err := os.Stat(envFile); err != nil {
		return nil
	}

	return godotenv.Overload(envFile)
}

func Environment() environment {
	return currentEnvironment
}

func IsLocal() bool {
	return currentEnvironment == Local
}
