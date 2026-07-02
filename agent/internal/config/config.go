package config

import (
	"io"
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	FocusedWindowInterval time.Duration `yaml:"focused_window_interval"`
	CpuPercentInterval    time.Duration `yaml:"cpu_percent_interval"`
}

func LoadConfig(path string) *Config {
	f, err := os.Open(path)
	if err != nil {
		log.Panicf("failed to open file %s: %s", path, err.Error())
	}
	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		log.Panicf("failed to read file %s: %s", path, err.Error())
	}

	cfg := new(Config)
	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		log.Panicf("failed to decode config: %s", err.Error())
	}

	return cfg
}
