package config

import (
	"agent/internal/utils"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	FocusedWindowInterval  time.Duration `yaml:"focused_window_collector_interval"`
	CpuPercentInterval     time.Duration `yaml:"cpu_percent_collector_interval"`
	MemoryInterval         time.Duration `yaml:"memory_collector_interval"`
	DiskInterval           time.Duration `yaml:"disk_collector_interval"`
	NetInterval            time.Duration `yaml:"net_io_collector_interval"`
	ProcessInterval        time.Duration `yaml:"process_collector_interval"`
	MetricsSendingInterval time.Duration `yaml:"metrics_sending_interval"`
}

func LoadConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", path, err)
	}
	defer utils.CloseWithLog(f)
	data, err := io.ReadAll(f)
	if err != nil {
		log.Panicf("failed to read file %s: %s", path, err.Error())
	}

	cfg := new(Config)
	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		log.Panicf("failed to decode config: %s", err.Error())
	}
	if cfg.NetInterval < time.Second {
		log.Panicf("net interval can't be less than 1 second")
	}
	return cfg, nil
}
