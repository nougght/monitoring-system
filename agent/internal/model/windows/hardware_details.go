//go:build windows

package model

type WinPhysicalMemory struct {
	BankLabel            string
	Capacity             uint64 // объем памяти
	ConfiguredClockSpeed uint32 // настроенная частота памяти
	DeviceLocator        string // расположение памяти
	FormFactor           uint16 // форм-фактор памяти
	HotSwappable         bool   // можно менять память без выключения системы
	Manufacturer         string // производитель
	// MemoryType           uint16  // тип памяти, некорректный в wmi
	Name             string // название памяти
	PartNumber       string // модель памяти
	Removable        bool   // можно ли вынимать память
	Replaceable      bool   // можно ли заменять память
	SerialNumber     string // серийный номер
	Speed            uint32 // поддерживаемя частота памяти
	SMBIOSMemoryType uint16 // тип памяти по smbios, нужно маппить
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
