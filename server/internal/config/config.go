package config

import (
	"fmt"

	"github.com/nougght/monitoring-system/shared/go/util"
)

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func (c *PostgresConfig) ConnString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}

type Config struct {
	Postgres *PostgresConfig
}

func MustLoadConfig() *Config {
	return &Config{
		Postgres: &PostgresConfig{
			Host:     util.MustGetEnvVar("POSTGRES_HOST"),
			Port:     util.MustGetEnvVar("POSTGRES_PORT"),
			User:     util.MustGetEnvVar("POSTGRES_USER"),
			Password: util.MustGetEnvVar("POSTGRES_PASSWORD"),
			DBName:   util.MustGetEnvVar("POSTGRES_DB"),
			SSLMode:  util.MustGetEnvVar("POSTGRES_SSL_MODE"),
		},
	}
}
