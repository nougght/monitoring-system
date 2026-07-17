package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
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
	Postgres PostgresConfig
}

func mustGetEnvVar(key string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	panic(fmt.Sprintf(`env variable "%s" not found`, key))
}

//nolint:unused
func getOptionalEnvVar(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

//nolint:unused
func getOptionalInt(key string, defaultValue int) int {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		log.Panicf(`env variable "%s" is not a valid integer: %s`, key, err)
	}
	return intValue
}

func LoadConfig() *Config {
	return &Config{
		Postgres: PostgresConfig{
			Host:     mustGetEnvVar("POSTGRES_HOST"),
			Port:     mustGetEnvVar("POSTGRES_PORT"),
			User:     mustGetEnvVar("POSTGRES_USER"),
			Password: mustGetEnvVar("POSTGRES_PASSWORD"),
			DBName:   mustGetEnvVar("POSTGRES_DB"),
			SSLMode:  mustGetEnvVar("POSTGRES_SSL_MODE"),
		},
	}
}
