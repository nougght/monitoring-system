package agent_server

import (
	"github.com/google/uuid"
	agent_model "github.com/nougght/monitoring-system/server/internal/model/agent"
	metrics_model "github.com/nougght/monitoring-system/server/internal/model/metrics"
	agentv1 "github.com/nougght/monitoring-system/shared/go/proto/gen/agent/v1"
	"github.com/nougght/monitoring-system/shared/go/util"
)

func convertMetricSampleFromProto(p *agentv1.MetricSample) *metrics_model.MetricSample {
	if p == nil {
		return nil
	}

	return &metrics_model.MetricSample{
		Kind:      int32(p.Kind),
		Label:     p.Label,
		Value:     p.Value,
		Timestamp: p.Timestamp.AsTime(),
	}
}

func convertMetricsBatchWithAgentIDFromProto(p *agentv1.MetricsBatch, agetnID uuid.UUID) *metrics_model.MetricsBatch {
	if p == nil {
		return nil
	}

	return &metrics_model.MetricsBatch{
		AgentID: agetnID,
		ID:      p.BatchId,
		Metrics: util.Map(p.Samples, convertMetricSampleFromProto),
	}
}

// specification

func convertPhysicalMemoryInfoFromProto(p *agentv1.PhysicalMemoryInfo) *agent_model.PhysicalMemoryInfo {
	if p == nil {
		return nil
	}

	return &agent_model.PhysicalMemoryInfo{
		DeviceLocator:        p.DeviceLocator,
		MemoryType:           p.MemoryType,
		Capacity:             p.Capacity,
		FormFactor:           p.FormFactor,
		Speed:                p.Speed,
		ConfiguredClockSpeed: p.ConfiguredClockSpeed,
		Manufacturer:         p.Manufacturer,
		ModelName:            p.ModelName,
		SerialNumber:         p.SerialNumber,
		BankLabel:            p.BankLabel,
		HotSwappable:         p.HotSwappable,
		Removable:            p.Removable,
		Replaceable:          p.Replaceable,
	}
}

func convertDiskSpecsFromProto(p *agentv1.DiskSpecs) *agent_model.DiskSpecs {
	if p == nil {
		return nil
	}
	return &agent_model.DiskSpecs{
		Device: p.Device,
		FsType: p.FsType,
		Total:  p.Total,
	}
}

func convertSpecsFromProto(p *agentv1.Specs) *agent_model.Specs {
	if p == nil {
		return nil
	}
	return &agent_model.Specs{
		HostSpecs: agent_model.HostSpecs{
			Hostname:      p.Host.Hostname,
			OSType:        p.Host.OsType,
			OS:            p.Host.Os,
			OSArch:        p.Host.OsArch,
			OSVersion:     p.Host.OsVersion,
			KernelVersion: p.Host.OsKernelVersion,
		},
		CpuSpecs: agent_model.CpuSpecs{
			ModelName:                     p.Cpu.ModelName,
			Architecture:                  p.Cpu.Architecture,
			Availability:                  p.Cpu.Availability,
			CurrentClockSpeed:             p.Cpu.CurrentClockSpeed,
			DataWidth:                     p.Cpu.DataWidth,
			L2CacheSize:                   p.Cpu.L2CacheSize,
			L3CacheSize:                   p.Cpu.L3CacheSize,
			Manufacturer:                  p.Cpu.Manufacturer,
			MaxClockSpeed:                 p.Cpu.MaxClockSpeed,
			NumberOfCores:                 p.Cpu.NumberOfCores,
			NumberOfEnabledCores:          p.Cpu.NumberOfEnabledCores,
			NumberOfLogicalProcessors:     p.Cpu.NumberOfLogicalProcessors,
			ProcessorId:                   p.Cpu.ProcessorId,
			SocketDesignation:             p.Cpu.SocketDesignation,
			Stepping:                      p.Cpu.Stepping,
			VirtualizationFirmwareEnabled: p.Cpu.VirtualizationFirmwareEnabled,
		},
		DiskSpecs: util.Map(p.Disk.Disk, convertDiskSpecsFromProto),
		MemorySpecs: agent_model.MemorySpecs{
			Total:              p.Memory.Total,
			PhysicalMemoryInfo: util.Map(p.Memory.PhysicalMemory, convertPhysicalMemoryInfoFromProto),
		},
	}
}
