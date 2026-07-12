//go:build windows

package collector

import (
	"agent/internal/model"
	windows "agent/internal/model/windows"
	"agent/internal/utils"

	"github.com/yusufpapurcu/wmi"
)

func convertWinPhysicalMemoryInfo(m windows.WinPhysicalMemory) model.PhysicalMemoryInfo {
	return model.PhysicalMemoryInfo{
		BankLabel:            m.BankLabel,
		Capacity:             m.Capacity,
		ConfiguredClockSpeed: m.ConfiguredClockSpeed,
		DeviceLocator:        m.DeviceLocator,
		FormFactor:           utils.ConvertWinPhysicalFormFactor(m.FormFactor),
		HotSwappable:         m.HotSwappable,
		Manufacturer:         m.Manufacturer,
		MemoryType:           utils.ConvertWinPhysicalMemoryType(m.SMBIOSMemoryType),
		ModelName:            m.PartNumber,
		Removable:            m.Removable,
		Replaceable:          m.Replaceable,
		SerialNumber:         m.SerialNumber,
	}
}

func GetMemoryDetails() ([]model.PhysicalMemoryInfo, error) {
	var detailsList []windows.WinPhysicalMemory
	err := wmi.Query("SELECT * FROM Win32_PhysicalMemory", &detailsList)
	if err != nil {
		return nil, err
	}
	physicalMemoryList := make([]model.PhysicalMemoryInfo, len(detailsList))
	for i, m := range detailsList {
		physicalMemoryList[i] = convertWinPhysicalMemoryInfo(m)
	}
	return physicalMemoryList, nil
}

func convertWinProcessorInfo(m windows.WinProcessor) model.CpuSpecs {
	return model.CpuSpecs{
		ModelName:                     m.Name,
		Architecture:                  utils.ConvertWinCpuArch(m.Architecture),
		Availability:                  utils.ConvertWinCpuAvailability(m.Availability),
		CurrentClockSpeed:             m.CurrentClockSpeed,
		DataWidth:                     m.DataWidth,
		L2CacheSize:                   m.L2CacheSize,
		L3CacheSize:                   m.L3CacheSize,
		Manufacturer:                  m.Manufacturer,
		MaxClockSpeed:                 m.MaxClockSpeed,
		NumberOfCores:                 m.NumberOfCores,
		NumberOfEnabledCore:           m.NumberOfEnabledCore,
		NumberOfLogicalProcessors:     m.NumberOfLogicalProcessors,
		ProcessorId:                   m.ProcessorId,
		SocketDesignation:             m.SocketDesignation,
		Stepping:                      m.Stepping,
		VirtualizationFirmwareEnabled: m.VirtualizationFirmwareEnabled,
	}
}

func GetProcessorDetails() ([]model.CpuSpecs, error) {
	var detailsList []windows.WinProcessor
	err := wmi.Query("SELECT * FROM Win32_Processor", &detailsList)
	if err != nil {
		return nil, err
	}
	cpuSpecsList := make([]model.CpuSpecs, len(detailsList))
	for i, m := range detailsList {
		cpuSpecsList[i] = convertWinProcessorInfo(m)
	}
	return cpuSpecsList, nil
}
