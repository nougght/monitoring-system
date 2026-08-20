package agent_model

import (
	"time"

	"github.com/google/uuid"
)

type Specs struct {
	AgentID     uuid.UUID    `json:"-" db:"agent_id"`
	UpdatedAt   time.Time    `json:"-" db:"updated_at"`
	HostSpecs   HostSpecs    `json:"hostSpecs"`
	CpuSpecs    CpuSpecs     `json:"cpuSpecs"`
	DiskSpecs   []*DiskSpecs `json:"diskSpecs"`
	MemorySpecs MemorySpecs  `json:"memorySpecs"`
}

type HostSpecs struct {
	Hostname      string `json:"hostname"`
	OSType        string `json:"osType"`
	OS            string `json:"os"`
	OSVersion     string `json:"osVersion"`
	OSArch        string `json:"osArch"`
	KernelVersion string `json:"kernelVersion"`
} // @Name HostSpecs

type CpuSpecs struct {
	ModelName                     string `json:"modelName"`
	Architecture                  string `json:"architecture"`
	Availability                  string `json:"availability"`
	CurrentClockSpeed             uint32 `json:"currentClockSpeed"`
	DataWidth                     uint32 `json:"data_width"`
	L2CacheSize                   uint32 `json:"l2CacheSize"`
	L3CacheSize                   uint32 `json:"l3CacheSize"`
	Manufacturer                  string `json:"manufacturer"`
	MaxClockSpeed                 uint32 `json:"maxClockSpeed"`
	NumberOfCores                 uint32 `json:"numberOfCores"`
	NumberOfEnabledCores          uint32 `json:"numberOfEnabledCores"`
	NumberOfLogicalProcessors     uint32 `json:"numberOfLogicalProcessors"`
	ProcessorId                   string `json:"processorId"`
	SocketDesignation             string `json:"socketDesignation"`
	Stepping                      string `json:"stepping"`
	VirtualizationFirmwareEnabled bool   `json:"virtualizationFirmwareEnabled"`
} // @Name CpuSpecs

type DiskSpecs struct {
	Device string `json:"device"`
	FsType string `json:"fsType"`
	Total  uint64 `json:"total"`
} // @Name DiskSpecs

type MemorySpecs struct {
	Total              uint64                `json:"total"`
	PhysicalMemoryInfo []*PhysicalMemoryInfo `json:"physicalMemoryInfo"`
} // @Name MemorySpecs

type PhysicalMemoryInfo struct {
	DeviceLocator        string `json:"deviceLocator"`
	MemoryType           string `json:"memoryType"`
	Capacity             uint64 `json:"capacity"`
	FormFactor           string `json:"form_factor"`
	Speed                uint32 `json:"speed"`
	ConfiguredClockSpeed uint32 `json:"configuredClockSpeed"`
	Manufacturer         string `json:"manufacturer"`
	ModelName            string `json:"modelName"`
	SerialNumber         string `json:"serialNumber"`
	BankLabel            string `json:"bankLabel"`
	HotSwappable         bool   `json:"hotSwappable"`
	Removable            bool   `json:"removable"`
	Replaceable          bool   `json:"replaceable"`
} // @Name PhysicalMemoryInfo
