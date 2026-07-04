import type { CpuSpecsArchitecture } from "../api/models/cpuSpecsArchitecture";
import type { CpuSpecsAvailability } from "../api/models/cpuSpecsAvailability";

export interface CpuSpecs {
    architecture?: CpuSpecsArchitecture;
    availability?: CpuSpecsAvailability;
    currentClockSpeed?: number;
    dataWidth?: number;
    l2CacheSize?: number;
    l3CacheSize?: number;
    manufacturer?: string;
    maxClockSpeed?: number;
    modelName?: string;
    numberOfCores?: number;
    numberOfEnabledCore?: number;
    numberOfLogicalProcessors?: number;
    processorId?: string;
    socketDesignation?: string;
    stepping?: string;
    virtualizationFirmwareEnabled?: boolean;
}

export interface HostSpecs {
    hostName?: string;
    os?: string;
    osType?: string;
    osVersion?: string;
    osKernelVersion?: string;
    osArch?: string;
}

export interface DiskSpecs {
    device?: string;
    fsType?: string;
    total?: number;
}
export interface Specs {
    cpu?: CpuSpecs;
    host?: HostSpecs;
    disk?: DiskSpecs[];
}