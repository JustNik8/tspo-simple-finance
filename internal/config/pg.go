package config

import (
	"errors"
	"fmt"
	"os"
)

type PGConfig interface {
	DSN() string
}

type pgConfig struct {
	dsn string
}

func NewPGConfig() (PGConfig, error) {
	dbUser, found := os.LookupEnv("DB_USER")
	if !found {
		return nil, errors.New("DB_USER not found")
	}

	dbPass, found := os.LookupEnv("DB_PASS")
	if !found {
		return nil, errors.New("DB_PASS not found")
	}

	dbHost, found := os.LookupEnv("DB_HOST")
	if !found {
		return nil, errors.New("DB_HOST not found")
	}

	dbPort, found := os.LookupEnv("DB_PORT")
	if !found {
		return nil, errors.New("DB_PORT not found")
	}

	dbName, found := os.LookupEnv("DB_NAME")
	if !found {
		return nil, errors.New("DB_NAME not found")
	}

	return &pgConfig{
		dsn: fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPass, dbHost, dbPort, dbName),
	}, nil
}

func (cfg *pgConfig) DSN() string {
	return cfg.dsn
}
