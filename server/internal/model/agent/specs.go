package agent_model

import (
	"time"

	"github.com/google/uuid"
)

type Specs struct {
	AgentID     uuid.UUID    `json:"-" db:"agent_id"`
	UpdatedAt   time.Time    `json:"-" db:"updated_at"`
	HostSpecs   HostSpecs    `json:"host_specs"`
	CpuSpecs    CpuSpecs     `json:"cpu_specs"`
	DiskSpecs   []*DiskSpecs `json:"disk_specs"`
	MemorySpecs MemorySpecs  `json:"memory_specs"`
}

type HostSpecs struct {
	Hostname      string `json:"hostname"`
	OSType        string `json:"os_type"`
	OS            string `json:"os"`
	OSVersion     string `json:"os_version"`
	OSArch        string `json:"os_arch"`
	KernelVersion string `json:"os_kernel_version"`
}

type CpuSpecs struct {
	ModelName                     string `json:"model_name"`
	Architecture                  string `json:"architecture"`
	Availability                  string `json:"availability"`
	CurrentClockSpeed             uint32 `json:"current_clock_speed"`
	DataWidth                     uint32 `json:"data_width"`
	L2CacheSize                   uint32 `json:"l2_cache_size"`
	L3CacheSize                   uint32 `json:"l3_cache_size"`
	Manufacturer                  string `json:"manufacturer"`
	MaxClockSpeed                 uint32 `json:"max_clock_speed"`
	NumberOfCores                 uint32 `json:"number_of_cores"`
	NumberOfEnabledCores          uint32 `json:"number_of_enabled_cores"`
	NumberOfLogicalProcessors     uint32 `json:"number_of_logical_processors"`
	ProcessorId                   string `json:"processor_id"`
	SocketDesignation             string `json:"socket_designation"`
	Stepping                      string `json:"stepping"`
	VirtualizationFirmwareEnabled bool   `json:"virtualization_firmware_enabled"`
}

type DiskSpecs struct {
	Device string `json:"device"`
	FsType string `json:"fs_type"`
	Total  uint64 `json:"total"`
}

type MemorySpecs struct {
	Total              uint64                `json:"total"`
	PhysicalMemoryInfo []*PhysicalMemoryInfo `json:"physical_memory"`
}

type PhysicalMemoryInfo struct {
	DeviceLocator        string `json:"device_locator"`
	MemoryType           string `json:"memory_type"`
	Capacity             uint64 `json:"capacity"`
	FormFactor           string `json:"form_factor"`
	Speed                uint32 `json:"speed"`
	ConfiguredClockSpeed uint32 `json:"configured_clock_speed"`
	Manufacturer         string `json:"manufacturer"`
	ModelName            string `json:"model_name"`
	SerialNumber         string `json:"serial_number"`
	BankLabel            string `json:"bank_label"`
	HotSwappable         bool   `json:"hot_swappable"`
	Removable            bool   `json:"removable"`
	Replaceable          bool   `json:"replaceable"`
}
