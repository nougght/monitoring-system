package mapper

import (
	"agent/internal/localserver/dto"
	"agent/internal/model"
)

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
		Timestamp:     timestamp,
	}
}
