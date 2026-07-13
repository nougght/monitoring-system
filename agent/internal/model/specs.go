package model

type CpuArchitecture string

const (
	CpuArchitectureX86     = "X86"
	CpuArchitectureX64     = "X64"
	CpuArchitectureMIPS    = "MIPS"
	CpuArchitectureAlpha   = "Alpha"
	CpuArchitecturePowerPC = "PowerPC"
	CpuArchitectureARM     = "ARM"
	CpuArchitectureIA64    = "IA64"
	CpuArchitectureItanium = "Itanium"
	CpuArchitectureARM64   = "ARM64"
	CpuArchitectureOther   = "Other"
)

type CpuAvailability string

const (
	CpuAvailabilityOther                 = "Other"
	CpuAvailabilityUnknown               = "Unknown"
	CpuAvailabilityRunning               = "Running"
	CpuAvailabilityWarning               = "Warning"
	CpuAvailabilityInTest                = "In Test"
	CpuAvailabilityNotApplicable         = "Not Applicable"
	CpuAvailabilityPowerOff              = "Power Off"
	CpuAvailabilityOffLine               = "Off Line"
	CpuAvailabilityOffDuty               = "Off Duty"
	CpuAvailabilityDegraded              = "Degraded"
	CpuAvailabilityNotInstalled          = "Not Installed"
	CpuAvailabilityInstallError          = "Install Error"
	CpuAvailabilityPowerSaveUnknown      = "Power Save - Unknown"
	CpuAvailabilityPowerSaveLowPowerMode = "Power Save - Low Power Mode"
	CpuAvailabilityPowerSaveStandby      = "Power Save - Standby"
	CpuAvailabilityPowerCycle            = "Power Cycle"
	CpuAvailabilityPowerSaveWarning      = "Power Save - Warning"
	CpuAvailabilityPaused                = "Paused"
	CpuAvailabilityNotReady              = "Not Ready"
	CpuAvailabilityNotConfigured         = "Not Configured"
	CpuAvailabilityQuiesced              = "Quiesced"
)

type PhysicalMemoryType string

const (
	PhysicalMemoryTypeUnknown = "Unknown"
	PhysicalMemoryTypeDDR2    = "DDR2"
	PhysicalMemoryTypeDDR3    = "DDR3"
	PhysicalMemoryTypeDDR4    = "DDR4"
	PhysicalMemoryTypeLPDDR2  = "LPDDR2"
	PhysicalMemoryTypeLPDDR3  = "LPDDR3"
	PhysicalMemoryTypeLPDDR4  = "LPDDR4"
	PhysicalMemoryTypeDDR5    = "DDR5"
	PhysicalMemoryTypeLPDDR5  = "LPDDR5"
	PhysicalMemoryTypeOther   = "Other"
)

type PhysicalMemoryFormFactor string

const (
	PhysicalMemoryFormFactorUnknown = "Unknown"

	PhysicalMemoryFormFactorOther       = "Other"
	PhysicalMemoryFormFactorSIP         = "SIP"
	PhysicalMemoryFormFactorDIP         = "DIP"
	PhysicalMemoryFormFactorZIP         = "ZIP"
	PhysicalMemoryFormFactorSOJ         = "SOJ"
	PhysicalMemoryFormFactorProprietary = "Proprietary"
	PhysicalMemoryFormFactorSIMM        = "SIMM"
	PhysicalMemoryFormFactorDIMM        = "DIMM"
	PhysicalMemoryFormFactorTSOP        = "TSOP"
	PhysicalMemoryFormFactorPGA         = "PGA"
	PhysicalMemoryFormFactorRIMM        = "RIMM"
	PhysicalMemoryFormFactorSODIMM      = "SODIMM"
	PhysicalMemoryFormFactorSRIMM       = "SRIMM"
	PhysicalMemoryFormFactorSMD         = "SMD"
	PhysicalMemoryFormFactorSSMP        = "SSMP"
	PhysicalMemoryFormFactorQFP         = "QFP"
	PhysicalMemoryFormFactorTQFP        = "TQFP"
	PhysicalMemoryFormFactorSOIC        = "SOIC"
	PhysicalMemoryFormFactorLCC         = "LCC"
	PhysicalMemoryFormFactorPLCC        = "PLCC"
	PhysicalMemoryFormFactorBGA         = "BGA"
	PhysicalMemoryFormFactorFPBGA       = "FPBGA"
	PhysicalMemoryFormFactorLGA         = "LGA"
)

type Specs struct {
	Host   HostSpecs     `json:"host"`
	CPU    CpuSpecs      `json:"cpu"`
	Disk   DiskSpecsList `json:"disk"`
	Memory MemorySpecs   `json:"memory"`
} // @name Specs

type HostSpecs struct {
	Hostname        string `json:"hostName"`        // имя хоста
	OsType          string `json:"osType"`          // семейство операционной системы
	Os              string `json:"os"`              // операционная система
	OsVersion       string `json:"osVersion"`       // версия операционной системы
	OsKernelVersion string `json:"osKernelVersion"` // версия ядра операционной системы
	OsArch          string `json:"osArch"`          // архитектура операционной системы
} // @name HostSpecs

type CpuSpecs struct {
	ModelName string `json:"modelName"` // модель процессора
	// архитектура https://learn.microsoft.com/en-us/windows/win32/cimwin32prov/win32-processor
	Architecture CpuArchitecture `json:"architecture" enums:"x86,x64,mips,alpha,powerpc,arm,ia64,itanium,arm64,other"`
	// состояние https://learn.microsoft.com/en-us/windows/win32/cimwin32prov/win32-processor
	Availability                  CpuAvailability `json:"availability" enums:"other,unknown,running,warning,inTest,notApplicable,powerOff,offLine,offDuty,degraded,notInstalled,installError,powerSaveUnknown,powerSaveLowPowerMode,powerSaveStandby,powerCycle,powerSaveWarning,paused,notReady,notConfigured,quiesced"`
	CurrentClockSpeed             uint32          `json:"currentClockSpeed"`             // текущая частота процессора
	DataWidth                     uint16          `json:"dataWidth"`                     // разрядность процессора
	L2CacheSize                   uint32          `json:"l2CacheSize"`                   // размер L2 кэша
	L3CacheSize                   uint32          `json:"l3CacheSize"`                   // размер L3 кэша
	Manufacturer                  string          `json:"manufacturer"`                  // производитель
	MaxClockSpeed                 uint32          `json:"maxClockSpeed"`                 // максимальная частота процессора(низкая точность)
	NumberOfCores                 uint32          `json:"numberOfCores"`                 // количество ядер
	NumberOfEnabledCore           uint32          `json:"numberOfEnabledCore"`           // количество доступных ядер
	NumberOfLogicalProcessors     uint32          `json:"numberOfLogicalProcessors"`     // количество логических ядер(потоков)
	ProcessorId                   string          `json:"processorId"`                   // идентификатор процессора
	SocketDesignation             string          `json:"socketDesignation"`             // сокет
	Stepping                      string          `json:"stepping"`                      // степпинг
	VirtualizationFirmwareEnabled bool            `json:"virtualizationFirmwareEnabled"` // включена ли виртуализация
} // @name CpuSpecs

type MemorySpecs struct {
	Total              uint64               `json:"total"`              // общий объем памяти
	PhysicalMemoryList []PhysicalMemoryInfo `json:"physicalMemoryList"` // физические плашки памяти
} // @name MemorySpecs

// информация о физической плашке памяти
type PhysicalMemoryInfo struct {
	DeviceLocator string `json:"deviceLocator"` // расположение памяти
	// тип памяти
	MemoryType PhysicalMemoryType `json:"memoryType" enums:"unknown,ddr2,ddr3,ddr4,lpddr2,lpddr3,lpddr4,ddr5,lpddr5,other"`
	Capacity   uint64             `json:"capacity"` // объем
	// форм-фактор памяти
	FormFactor           PhysicalMemoryFormFactor `json:"formFactor" enums:"unknown,other,sip,dip,zip,soj,proprietary,simm,dimm,tsop,pga,rimm,sodimm,srimm,smd,ssmp,qfp,tqfp,soic,lcc,plcc,bga,fpgba,lga"`
	Speed                uint32                   `json:"speed"`                // поддерживаемя частота памяти
	ConfiguredClockSpeed uint32                   `json:"configuredClockSpeed"` // настроенная частота памяти
	Manufacturer         string                   `json:"manufacturer"`         // производитель
	ModelName            string                   `json:"modelName"`            // модель/партия памяти
	SerialNumber         string                   `json:"serialNumber"`         // серийный номер
	BankLabel            string                   `json:"bankLabel"`            // подключение памяти
	HotSwappable         bool                     `json:"hotSwappable"`         // можно менять память без выключения системы
	Removable            bool                     `json:"removable"`            // можно ли вынимать память
	Replaceable          bool                     `json:"replaceable"`          // можно ли заменять память
} // @name PhysicalMemoryInfo

type DiskSpecs struct {
	Device string `json:"device"`
	FsType string `json:"fsType"`
	Total  uint64 `json:"total"`
} // @name DiskSpecs

type DiskSpecsList []DiskSpecs // @name DiskSpecsList
