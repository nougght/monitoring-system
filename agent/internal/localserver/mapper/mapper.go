package mapper

import (
	"agent/internal/localserver/dto"
	"agent/internal/model"
)

func ConvertProcessToDto(process *model.Process) *dto.Process {
	return &dto.Process{
		PID:        process.Pid,
		Name:       process.Name,
		ParentPid:  process.ParentPid,
		CPUPercent: process.CPUPercent,
		MemoryUsed: process.MemoryUsed,
	}
}

func ConvertProcessListToDto(processList []model.Process) []dto.Process {
	processListDto := make([]dto.Process, len(processList))
	for i, process := range processList {
		processListDto[i] = *ConvertProcessToDto(&process)
	}
	return processListDto
}

func ConvertMetricsToDto(metrics *model.Metrics) *dto.Metrics {
	focusedWindow := metrics.FocusedWindow.Value()
	cpuPercent := metrics.CpuPercent.Value()
	memoryUsed := metrics.MemoryUsage.Value()
	diskUsage := metrics.DiskUsage.Value()
	uploadMbps := metrics.NetworkUsage.UploadMbps()
	downloadMbps := metrics.NetworkUsage.DownloadMbps()
	timestamp := metrics.Timestamp
	return &dto.Metrics{
		FocusedWindow: &focusedWindow,
		CpuPercent:    &cpuPercent,
		MemoryUsed:    &memoryUsed,
		DiskUsage:     &diskUsage,
		UploadMbps:    &uploadMbps,
		DownloadMbps:  &downloadMbps,
		ProcessList:   ConvertProcessListToDto(metrics.Process.ProcessList()),
		Timestamp:     timestamp,
	}
}
