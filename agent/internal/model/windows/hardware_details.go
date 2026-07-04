//go:build windows

package model

type WinPhysicalMemory struct {
	BankLabel     string `json:"bankLabel"`
	Capacity      uint64 `json:"capacity"`
	DataWidth     uint16 `json:"dataWidth"`
	Description   string `json:"description"`
	DeviceLocator string `json:"deviceLocator"`
	FormFactor    int    `json:"formFactor"`
	Speed         uint32
	Manufacturer  string `json:"manufacturer"`
	Model         string `json:"model"`
	PartNumber    string `json:"partNumber"`
}

type WinProcessor struct {
	Architecture                  uint16 // архитектура, нужно маппить
	Availability                  uint16 // состояние, нужно маппить
	CurrentClockSpeed             uint32
	DataWidth                     uint16 // разрядность процессора
	L2CacheSize                   uint32
	L3CacheSize                   uint32
	Manufacturer                  string // производитель
	MaxClockSpeed                 uint32
	Name                          string // модель процессора
	NumberOfCores                 uint32
	NumberOfEnabledCore           uint32
	NumberOfLogicalProcessors     uint32
	ProcessorId                   string
	SocketDesignation             string // сокет
	Stepping                      string // степпинг
	VirtualizationFirmwareEnabled bool   // включена ли виртуализация
	// AddressWidth             uint16    // разрядность ОС, есть в gopsutil
	// Family             		uint16 // слишком много значений
	// LoadPercentage      		uint16 // нагрузка, есть в gopsutil
}
