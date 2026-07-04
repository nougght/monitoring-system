//go:build windows

package collector

import (
	model "agent/internal/model/windows"

	"github.com/yusufpapurcu/wmi"
)

func GetMemoryDetails() ([]model.WinPhysicalMemory, error) {
	var detailsList []model.WinPhysicalMemory
	err := wmi.Query("SELECT * FROM Win32_PhysicalMemory", &detailsList)
	if err != nil {
		return nil, err
	}
	return detailsList, nil
}

func GetProcessorDetails() ([]model.WinProcessor, error) {
	var detailsList []model.WinProcessor
	err := wmi.Query("SELECT * FROM Win32_Processor", &detailsList)
	if err != nil {
		return nil, err
	}
	return detailsList, nil
}
