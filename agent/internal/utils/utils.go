package utils

import "agent/internal/model"

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
