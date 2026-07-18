package config

import (
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
	cfg := new(Config)
	util.MustReadYaml(path, cfg)
	if cfg.NetInterval < time.Second {
		log.Panicf("net interval can't be less than 1 second")
	}
	return cfg
}
