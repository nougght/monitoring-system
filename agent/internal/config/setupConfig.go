package config

import (
	"agent/internal/model"
	"fmt"

	"github.com/nougght/monitoring-system/shared/go/util"
)

type SetupConfig struct {
	path              string `yaml:"-"`
	EnrollmentKey     string `yaml:"enrollment_key"`
	EnrollmentAddress string `yaml:"enrollment_address"`
	ServerAddress     string `yaml:"server_address"`
	CaPath            string `yaml:"ca_path"`
	KeyPath           string `yaml:"key_path"`
	CertPath          string `yaml:"cert_path"`
	// TODO: add agent mode
}

func LoadSetupConfig(path string) (cfg *SetupConfig, err error) {
	if err := util.ReadYaml(path, &cfg); err != nil {
		return nil, fmt.Errorf("failed to read yaml setup config: %w", err)
	}
	cfg.path = path
	if cfg.CaPath == "" {
		cfg.CaPath = model.DefaultCAPath
	}
	if cfg.CertPath == "" {
		cfg.CertPath = model.DefaultAgentCertPath
	}
	if cfg.KeyPath == "" {
		cfg.KeyPath = model.DefaultAgentKeyPath
	}
	return cfg, nil
}

func (c *SetupConfig) SetEnrollmentKeyUsed() error {
	c.EnrollmentKey = model.EnrollmentKeyUsed
	err := util.SaveYaml(c.path, c)
	if err != nil {
		return fmt.Errorf("failed to save updated yaml config: %w", err)
	}
	return nil
}
