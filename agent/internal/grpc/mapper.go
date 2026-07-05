package grpc_client

import (
	"agent/internal/model"
	"strconv"

	pb "github.com/nougght/monitoring-system/shared/go/proto/gen/agent/v1"
)

func convertSpecsToProto(specs *model.Specs) *pb.Specs {
	return &pb.Specs{
		Host:   convertHostSpecsToProto(&specs.Host),
		Cpu:    convertCpuSpecsToProto(&specs.CPU),
		Disk:   convertDiskSpecsListToProto(&specs.Disk),
		Memory: convertMemorySpecsToProto(&specs.Memory),
	}
}

func convertHostSpecsToProto(host *model.HostSpecs) *pb.HostSpecs {
	return &pb.HostSpecs{
		Hostname:        host.Hostname,
		OsType:          host.OsType,
		Os:              host.Os,
		OsVersion:       host.OsVersion,
		OsKernelVersion: host.OsKernelVersion,
		OsArch:          host.OsArch,
	}
}

func convertCpuSpecsToProto(cpu *model.CpuSpecs) *pb.CpuSpecs {
	return &pb.CpuSpecs{
		ModelName:                     cpu.ModelName,
		Architecture:                  string(cpu.Architecture),
		Availability:                  string(cpu.Availability),
		CurrentClockSpeed:             cpu.CurrentClockSpeed,
		DataWidth:                     uint32(cpu.DataWidth),
		L2CacheSize:                   cpu.L2CacheSize,
		L3CacheSize:                   cpu.L3CacheSize,
		Manufacturer:                  cpu.Manufacturer,
		MaxClockSpeed:                 cpu.MaxClockSpeed,
		NumberOfCores:                 cpu.NumberOfCores,
		NumberOfEnabledCores:          cpu.NumberOfEnabledCore,
		NumberOfLogicalProcessors:     cpu.NumberOfLogicalProcessors,
		ProcessorId:                   cpu.ProcessorId,
		SocketDesignation:             cpu.SocketDesignation,
		Stepping:                      cpu.Stepping,
		VirtualizationFirmwareEnabled: cpu.VirtualizationFirmwareEnabled,
	}
}

func convertDiskSpecsListToProto(disk *model.DiskSpecsList) *pb.DiskSpecsList {
	pbSpecs := pb.DiskSpecsList{
		Disk: make([]*pb.DiskSpecs, len(*disk)),
	}
	for i, d := range *disk {
		pbSpecs.Disk[i] = convertDiskSpecsToProto(&d)
	}
	return &pbSpecs
}

func convertDiskSpecsToProto(disk *model.DiskSpecs) *pb.DiskSpecs {
	return &pb.DiskSpecs{
		Device: disk.Device,
		FsType: disk.FsType,
		Total:  disk.Total,
	}
}

func convertMemorySpecsToProto(memory *model.MemorySpecs) *pb.MemorySpecs {
	pbSpecs := pb.MemorySpecs{
		Total:          strconv.FormatUint(memory.Total, 10),
		PhysicalMemory: make([]*pb.PhysicalMemoryInfo, len(memory.PhysicalMemoryList)),
	}
	for i, m := range memory.PhysicalMemoryList {
		pbSpecs.PhysicalMemory[i] = convertPhysicalMemoryInfoToProto(&m)
	}
	return &pbSpecs
}

func convertPhysicalMemoryInfoToProto(physicalMemory *model.PhysicalMemoryInfo) *pb.PhysicalMemoryInfo {
	return &pb.PhysicalMemoryInfo{
		DeviceLocator:        physicalMemory.DeviceLocator,
		MemoryType:           string(physicalMemory.MemoryType),
		Capacity:             physicalMemory.Capacity,
		FormFactor:           string(physicalMemory.FormFactor),
		Speed:                physicalMemory.Speed,
		ConfiguredClockSpeed: physicalMemory.ConfiguredClockSpeed,
		Manufacturer:         physicalMemory.Manufacturer,
		ModelName:            physicalMemory.ModelName,
		SerialNumber:         physicalMemory.SerialNumber,
		BankLabel:            physicalMemory.BankLabel,
		HotSwappable:         physicalMemory.HotSwappable,
		Removable:            physicalMemory.Removable,
		Replaceable:          physicalMemory.Replaceable,
	}
}
