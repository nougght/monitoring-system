package config

import (
	"fmt"
	"log"
	"time"

	"github.com/nougght/monitoring-system/shared/go/util"
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

func MustLoadConfig(path string) *Config {
	cfg, err := LoadConfig(path)
	if err != nil {
		log.Panic(err)
	}

	return cfg
}

func LoadConfig(path string) (*Config, error) {
	cfg := new(Config)
	err := util.ReadYaml(path, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to read yaml config")
	}

	if cfg.NetInterval < time.Second {
		return nil, fmt.Errorf("net interval can't be less than 1 second")
	}
	return cfg, nil

}
