package env

import (
	"fmt"
	"os"
)

const dsnEnvName = "PG_DSN"

type pgConfig struct {
	dsn string
}

func NewPgConfig() (*pgConfig, error) {
	dsn := os.Getenv(dsnEnvName)
	if len(dsn) == 0 {
		return nil, fmt.Errorf("environment variable %s is not set", dsnEnvName)
	}

	return &pgConfig{dsn: dsn}, nil
}

func (cfg *pgConfig) DSN() string {
	return cfg.dsn
}
