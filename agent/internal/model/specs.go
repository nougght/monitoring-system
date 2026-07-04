package model

//
type Specs struct {
	Host HostSpecs     `json:"host"`
	CPU  CpuSpecs      `json:"cpu"`
	Disk DiskSpecsList `json:"disk"`
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
	ModelName        string `json:"modelName"`
	CoreCount        int    `json:"coreCount"`
	LogicalCoreCount int    `json:"logicalCoreCount"`
} // @name CpuSpecs

type MemorySpecs struct {
	Total uint64 `json:"total"`
} // @name MemorySpecs

type DiskSpecs struct {
	Device string `json:"device"`
	FsType string `json:"fsType"`
	Total  uint64 `json:"total"`
} // @name DiskSpecs

type DiskSpecsList []DiskSpecs // @name DiskSpecsList
