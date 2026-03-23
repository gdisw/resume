package env

import (
	"fmt"
	"path"
)

type loadCfg struct {
	envPath string
}

func (cfg loadCfg) ForceLoad() bool {
	return cfg.envPath != ""
}

func (cfg loadCfg) defaultEnvFilePath() string {
	if cfg.envPath != "" {
		return path.Join(cfg.envPath, ".env")
	}

	return ".env"
}

func (cfg loadCfg) envFilePath(env environment) string {
	envFileName := fmt.Sprintf(".env.%s", env)

	if cfg.envPath != "" {
		return path.Join(cfg.envPath, envFileName)
	}

	return envFileName
}

type loadOption func(*loadCfg)

func WithEnvPath(path string) loadOption {
	return func(cfg *loadCfg) {
		cfg.envPath = path
	}
}
