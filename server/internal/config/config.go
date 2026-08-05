package config

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
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
type GRPCConfig struct {
	MainPort       int
	EnrollmentPort int
}

type SettingsConfig struct {
	AgentEnrollmentKeyLength int    `yaml:"enrollment_key_length"`
	Address                  string `yaml:"address"`
}

func (c *PostgresConfig) ConnString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}

type Config struct {
	Postgres       *PostgresConfig
	HTTP           *HTTPConfig
	GRPC           *GRPCConfig
	SettingsConfig *SettingsConfig
}

func MustLoadConfig(path string) *Config {
	if err := godotenv.Load(); err != nil {
		log.Println(err.Error())
	}
	cfg := &Config{
		SettingsConfig: &SettingsConfig{},
	}
	util.MustReadYaml(path, cfg.SettingsConfig)
	if cfg.SettingsConfig.AgentEnrollmentKeyLength < 10 {
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

	serverPort := util.MustGetPortEnvVar("HTTP_SERVER_PORT")
	cfg.HTTP = &HTTPConfig{
		ServerPort: serverPort,
	}

	grpcPort := util.MustGetPortEnvVar("GRPC_SERVER_PORT")
	grpcEnrollmentPort := util.MustGetPortEnvVar("GRPC_ENROLLMENT_PORT")
	cfg.GRPC = &GRPCConfig{
		MainPort:       grpcPort,
		EnrollmentPort: grpcEnrollmentPort,
	}
	return cfg
}
