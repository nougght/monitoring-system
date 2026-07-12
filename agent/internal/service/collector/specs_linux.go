//go:build linux

package collector

import (
	"agent/internal/model"
	linux "agent/internal/model/linux"
)

func convertLinuxPhysicalMemoryInfo(m linux.LinuxPhysicalMemory) model.PhysicalMemoryInfo {
	return model.PhysicalMemoryInfo{}
}

func convertLinuxProcessorInfo(m linux.LinuxProcessor) model.CpuSpecs {
	return model.CpuSpecs{}
}

func GetMemoryDetails() ([]model.PhysicalMemoryInfo, error) {
	var detailsList []linux.LinuxPhysicalMemory
	// TODO: implement
	physicalMemoryList := make([]model.PhysicalMemoryInfo, len(detailsList))
	for i, m := range detailsList {
		physicalMemoryList[i] = convertLinuxPhysicalMemoryInfo(m)
	}
	return physicalMemoryList, nil
}

func GetProcessorDetails() ([]model.CpuSpecs, error) {
	var detailsList []linux.LinuxProcessor
	// TODO: implement
	cpuSpecsList := make([]model.CpuSpecs, len(detailsList))
	for i, m := range detailsList {
		cpuSpecsList[i] = convertLinuxProcessorInfo(m)
	}
	return cpuSpecsList, nil
}
