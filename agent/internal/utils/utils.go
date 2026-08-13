package utils

import (
	"agent/internal/model"
)

// https://learn.microsoft.com/en-us/windows/win32/cimwin32prov/win32-processor
func ConvertWinCpuArch(arch uint16) model.CpuArchitecture {
	switch arch {
	case 0:
		return model.CpuArchitectureX86
	case 1:
		return model.CpuArchitectureMIPS
	case 2:
		return model.CpuArchitectureAlpha
	case 3:
		return model.CpuArchitecturePowerPC
	case 5:
		return model.CpuArchitectureARM
	case 6:
		return model.CpuArchitectureIA64
	case 9:
		return model.CpuArchitectureX64
	case 12:
		return model.CpuArchitectureARM64
	default:
		return model.CpuArchitectureOther
	}
}

// https://learn.microsoft.com/en-us/windows/win32/cimwin32prov/win32-processor
func ConvertWinCpuAvailability(availability uint16) model.CpuAvailability {
	switch availability {
	case 1:
		return model.CpuAvailabilityOther
	case 2:
		return model.CpuAvailabilityUnknown
	case 3:
		return model.CpuAvailabilityRunning
	case 4:
		return model.CpuAvailabilityWarning
	case 5:
		return model.CpuAvailabilityInTest
	case 6:
		return model.CpuAvailabilityNotApplicable
	case 7:
		return model.CpuAvailabilityPowerOff
	case 8:
		return model.CpuAvailabilityOffLine
	case 9:
		return model.CpuAvailabilityOffDuty
	case 10:
		return model.CpuAvailabilityDegraded
	case 11:
		return model.CpuAvailabilityNotInstalled
	case 12:
		return model.CpuAvailabilityInstallError
	case 13:
		return model.CpuAvailabilityPowerSaveUnknown
	case 14:
		return model.CpuAvailabilityPowerSaveLowPowerMode
	case 15:
		return model.CpuAvailabilityPowerSaveStandby
	case 16:
		return model.CpuAvailabilityPowerCycle
	case 17:
		return model.CpuAvailabilityPowerSaveWarning
	case 18:
		return model.CpuAvailabilityPaused
	case 19:
		return model.CpuAvailabilityNotReady
	case 20:
		return model.CpuAvailabilityNotConfigured
	case 21:
		return model.CpuAvailabilityQuiesced
	default:
		return model.CpuAvailabilityOther
	}
}

func ConvertWinPhysicalMemoryType(memoryType uint16) model.PhysicalMemoryType {
	switch memoryType {
	case 20:
		return model.PhysicalMemoryTypeDDR2
	case 21:
		return model.PhysicalMemoryTypeDDR3
	case 24:
		return model.PhysicalMemoryTypeDDR4
	case 27:
		return model.PhysicalMemoryTypeLPDDR2
	case 29:
		return model.PhysicalMemoryTypeLPDDR3
	case 30:
		return model.PhysicalMemoryTypeLPDDR4
	case 34:
		return model.PhysicalMemoryTypeDDR5
	case 35:
		return model.PhysicalMemoryTypeLPDDR5
	default:
		return model.PhysicalMemoryTypeUnknown
	}
}

func ConvertWinPhysicalFormFactor(memoryType uint16) model.PhysicalMemoryFormFactor {
	switch memoryType {
	case 0:
		return model.PhysicalMemoryFormFactorUnknown
	case 1:
		return model.PhysicalMemoryFormFactorOther
	case 2:
		return model.PhysicalMemoryFormFactorSIP
	case 3:
		return model.PhysicalMemoryFormFactorDIP
	case 4:
		return model.PhysicalMemoryFormFactorZIP
	case 5:
		return model.PhysicalMemoryFormFactorSOJ
	case 6:
		return model.PhysicalMemoryFormFactorProprietary
	case 7:
		return model.PhysicalMemoryFormFactorSIMM
	case 8:
		return model.PhysicalMemoryFormFactorDIMM
	case 9:
		return model.PhysicalMemoryFormFactorTSOP
	case 10:
		return model.PhysicalMemoryFormFactorPGA
	case 11:
		return model.PhysicalMemoryFormFactorRIMM
	case 12:
		return model.PhysicalMemoryFormFactorSODIMM
	case 13:
		return model.PhysicalMemoryFormFactorSRIMM
	case 14:
		return model.PhysicalMemoryFormFactorSMD
	case 15:
		return model.PhysicalMemoryFormFactorSSMP
	case 16:
		return model.PhysicalMemoryFormFactorQFP
	case 17:
		return model.PhysicalMemoryFormFactorTQFP
	case 18:
		return model.PhysicalMemoryFormFactorSOIC
	case 19:
		return model.PhysicalMemoryFormFactorLCC
	case 20:
		return model.PhysicalMemoryFormFactorPLCC
	case 21:
		return model.PhysicalMemoryFormFactorBGA
	case 22:
		return model.PhysicalMemoryFormFactorFPBGA
	case 23:
		return model.PhysicalMemoryFormFactorLGA
	default:
		return model.PhysicalMemoryFormFactorUnknown
	}
}
