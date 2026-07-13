import type { CpuSpecsArchitecture } from "../api/models/cpuSpecsArchitecture";
import type { CpuSpecsAvailability } from "../api/models/cpuSpecsAvailability";
import type { PhysicalMemoryInfoFormFactor } from "../api/models/physicalMemoryInfoFormFactor";
import type { PhysicalMemoryInfoMemoryType } from "../api/models/physicalMemoryInfoMemoryType";

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

export interface MemorySpecs {
    /** физические плашки памяти */
    physicalMemoryList?: PhysicalMemoryInfo[];
    /** общий объем памяти */
    total?: number;
}
export interface PhysicalMemoryInfo {
    /** подключение памяти */
    bankLabel?: string;
    /** объем */
    capacity?: number;
    /** настроенная частота памяти */
    configuredClockSpeed?: number;
    /** расположение памяти */
    deviceLocator?: string;
    /** форм-фактор памяти */
    formFactor?: PhysicalMemoryInfoFormFactor;
    /** можно менять память без выключения системы */
    hotSwappable?: boolean;
    /** производитель */
    manufacturer?: string;
    /** тип памяти */
    memoryType?: PhysicalMemoryInfoMemoryType;
    /** модель/партия памяти */
    modelName?: string;
    /** можно ли вынимать память */
    removable?: boolean;
    /** можно ли заменять память */
    replaceable?: boolean;
    /** серийный номер */
    serialNumber?: string;
    /** поддерживаемя частота памяти */
    speed?: number;
}

export interface Specs {
    cpu?: CpuSpecs;
    host?: HostSpecs;
    disk?: DiskSpecs[];
    memory?: MemorySpecs;
}