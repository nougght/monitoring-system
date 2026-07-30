package config

import (
	"fmt"
	"log"

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

type HTTPConfig struct {
	ServerPort int
}

func (c *PostgresConfig) ConnString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}

type Config struct {
	Postgres                 *PostgresConfig
	Http                     *HTTPConfig
	AgentEnrollmentKeyLength int `yaml:"enrollment_key_length"`
}

func MustLoadConfig(path string) *Config {
	cfg := new(Config)
	util.MustReadYaml(path, cfg)
	if cfg.AgentEnrollmentKeyLength < 10 {
		log.Panicf("agent enrollment key length can't be less than 10")
	}

	cfg.Postgres = &PostgresConfig{
		Host:     util.MustGetEnvVar("POSTGRES_HOST"),
		Port:     util.MustGetEnvVar("POSTGRES_PORT"),
		User:     util.MustGetEnvVar("POSTGRES_USER"),
		Password: util.MustGetEnvVar("POSTGRES_PASSWORD"),
		DBName:   util.MustGetEnvVar("POSTGRES_DB"),
		SSLMode:  util.MustGetEnvVar("POSTGRES_SSL_MODE"),
	}

	serverPort := util.MustGetIntEnvVar("HTTP_SERVER_PORT")
	if serverPort <= 1024 || serverPort > 65535 {
		log.Panicf("http server port must be in range 1025-65535")
	}
	cfg.Http = &HTTPConfig{
		ServerPort: serverPort,
	}
	return cfg
}
