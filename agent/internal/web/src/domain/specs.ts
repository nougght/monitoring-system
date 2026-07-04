

export interface CpuSpecs {
    coreCount?: number;
    logicalCoreCount?: number;
    modelName?: string;
}

export interface HostSpecs {
    hostName?: string;
    os?: string;
    osType?: string;
    osVersion?: string;
    osKernelVersion?: string;
    osArch?: string;
}
export interface Specs {
    cpu?: CpuSpecs;
    host?: HostSpecs;
}