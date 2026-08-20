

import type { Agent as AgentDTO, AgentSpecs } from "./models";
import type { Agent } from "../domain/agent";
import type { Specs } from '../../../shared/ui/src/domain/specs';

export const convertAgentFromDTO = (agentDTO: AgentDTO): Agent => {
    return {
        id: agentDTO.id ?? '',
        name: agentDTO.name ?? '',
        description: agentDTO.description,
        createdAt: agentDTO.createdAt ?? '',
        status: agentDTO.status,
        lastSeenAt: agentDTO.lastSeenAt ? new Date(agentDTO.lastSeenAt) : undefined,
        isOnline: agentDTO.isOnline ?? false
    }
}

export const convertSpecsFromDTO = (specsDTO: AgentSpecs): Specs => {
    return {
        cpu: {
            architecture: specsDTO.cpuSpecs?.architecture,
            availability: specsDTO.cpuSpecs?.availability,
            currentClockSpeed: specsDTO.cpuSpecs?.currentClockSpeed,
            dataWidth: specsDTO.cpuSpecs?.data_width,
            l2CacheSize: specsDTO.cpuSpecs?.l2CacheSize,
            l3CacheSize: specsDTO.cpuSpecs?.l3CacheSize,
            manufacturer: specsDTO.cpuSpecs?.manufacturer,
            maxClockSpeed: specsDTO.cpuSpecs?.maxClockSpeed,
            modelName: specsDTO.cpuSpecs?.modelName,
            numberOfCores: specsDTO.cpuSpecs?.numberOfCores,
            numberOfEnabledCore: specsDTO.cpuSpecs?.numberOfEnabledCores,
            numberOfLogicalProcessors: specsDTO.cpuSpecs?.numberOfLogicalProcessors,
            processorId: specsDTO.cpuSpecs?.processorId,
            socketDesignation: specsDTO.cpuSpecs?.socketDesignation,
            stepping: specsDTO.cpuSpecs?.stepping,
            virtualizationFirmwareEnabled: specsDTO.cpuSpecs?.virtualizationFirmwareEnabled,
        },
        host: {
            hostName: specsDTO.hostSpecs?.hostname,
            os: specsDTO.hostSpecs?.os,
            osType: specsDTO.hostSpecs?.osType,
            osVersion: specsDTO.hostSpecs?.osVersion,
            osKernelVersion: specsDTO.hostSpecs?.kernelVersion,
            osArch: specsDTO.hostSpecs?.osArch,
        },
        disk: specsDTO.diskSpecs?.map(disk => ({
            device: disk.device,
            fsType: disk.fsType,
            total: disk.total,
        })),
        memory: {
            physicalMemoryList: specsDTO.memorySpecs?.physicalMemoryInfo,
            total: specsDTO.memorySpecs?.total,
        },
    }
}