package model

type Specs struct {
	Host HostSpecs `json:"host"`
	CPU  CpuSpecs  `json:"cpu"`
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
