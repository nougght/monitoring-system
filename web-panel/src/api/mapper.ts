

import type { Agent as AgentDTO, AgentSpecs, AgentWithCPUUsageDTO, AgentWithMemoryUsageDTO } from "./models";
import type { Agent } from "../domain/agent";
import type { Specs } from '../../../shared/ui/src/domain/specs';
import type { AgentsOverview as OverviewDTO } from "./models";
import type { AgentWithCPUUsage, AgentWithMemoryUsage, Overview } from "../domain/overview";
import type { AgentsSummary as AgentsSummaryDTO } from "./models";
import type { AgentsSummary } from "../domain/overview";
import type { TopNAgents as TopNAgentsDTO } from "./models";
import type { TopNAgents } from "../domain/overview";
import type {AgentGroup as AgentGroupDTO} from "./models"
import type {Group} from "../domain/group.ts";


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

export const convertGroupFromDTO = (dto: AgentGroupDTO): Group => {
    return {
        id: dto.id!,
        name: dto.name!,
        description: dto.description,
        createdAt: new Date(dto.createdAt!),
        deletedAt: dto.deletedAt ? new Date(dto.deletedAt) : undefined,
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

const convertAgentsSummaryFromDTO = (agentsSummaryDTO: AgentsSummaryDTO): AgentsSummary => {
    return {
        onlineAgents: agentsSummaryDTO.onlineAgents ?? 0,
        totalAgents: agentsSummaryDTO.totalAgents ?? 0,
        averageCPUUsage: agentsSummaryDTO.averageCPUUsage ?? 0,
        averageMemoryUsage: agentsSummaryDTO.averageMemoryUsage ?? 0,
        cpuUsageDistribution: agentsSummaryDTO.cpuUsageDistribution ? 
        {
            high: agentsSummaryDTO.cpuUsageDistribution.high ? { percent: agentsSummaryDTO.cpuUsageDistribution.high.percent ?? 0, count: agentsSummaryDTO.cpuUsageDistribution.high.count ?? 0 } : { percent: 0, count: 0 },
            medium: agentsSummaryDTO.cpuUsageDistribution.medium ? { percent: agentsSummaryDTO.cpuUsageDistribution.medium.percent ?? 0, count: agentsSummaryDTO.cpuUsageDistribution.medium.count ?? 0 } : { percent: 0, count: 0 },
            low: agentsSummaryDTO.cpuUsageDistribution.low ? { percent: agentsSummaryDTO.cpuUsageDistribution.low.percent ?? 0, count: agentsSummaryDTO.cpuUsageDistribution.low.count ?? 0 } : { percent: 0, count: 0 },
        } : {
            high: { percent: 0, count: 0 },
            medium: { percent: 0, count: 0 },
            low: { percent: 0, count: 0 },
        },
        memoryUsageDistribution: agentsSummaryDTO.memoryUsageDistribution ?
        {
            high: agentsSummaryDTO.memoryUsageDistribution.high ? { percent: agentsSummaryDTO.memoryUsageDistribution.high.percent ?? 0, count: agentsSummaryDTO.memoryUsageDistribution.high.count ?? 0 } : { percent: 0, count: 0 },
            medium: agentsSummaryDTO.memoryUsageDistribution.medium ? { percent: agentsSummaryDTO.memoryUsageDistribution.medium.percent ?? 0, count: agentsSummaryDTO.memoryUsageDistribution.medium.count ?? 0 } : { percent: 0, count: 0 },
            low: agentsSummaryDTO.memoryUsageDistribution.low ? { percent: agentsSummaryDTO.memoryUsageDistribution.low.percent ?? 0, count: agentsSummaryDTO.memoryUsageDistribution.low.count ?? 0 } : { percent: 0, count: 0 },
        } : {
            high: { percent: 0, count: 0 },
            medium: { percent: 0, count: 0 },
            low: { percent: 0, count: 0 },
        },
        diskUsageDistribution: agentsSummaryDTO.diskUsageDistribution ?
        {
            high: agentsSummaryDTO.diskUsageDistribution.high ? { percent: agentsSummaryDTO.diskUsageDistribution.high.percent ?? 0, count: agentsSummaryDTO.diskUsageDistribution.high.count ?? 0 } : { percent: 0, count: 0 },
            medium: agentsSummaryDTO.diskUsageDistribution.medium ? { percent: agentsSummaryDTO.diskUsageDistribution.medium.percent ?? 0, count: agentsSummaryDTO.diskUsageDistribution.medium.count ?? 0 } : { percent: 0, count: 0 },
            low: agentsSummaryDTO.diskUsageDistribution.low ? { percent: agentsSummaryDTO.diskUsageDistribution.low.percent ?? 0, count: agentsSummaryDTO.diskUsageDistribution.low.count ?? 0 } : { percent: 0, count: 0 },
        } : {
            high: { percent: 0, count: 0 },
            medium: { percent: 0, count: 0 },
            low: { percent: 0, count: 0 },
        },
    }
}

const convertTopNAgentsFromDTO = (topNAgentsDTO: TopNAgentsDTO): TopNAgents => {
    return {
        cpuUsage: topNAgentsDTO.cpuUsage ? topNAgentsDTO.cpuUsage.map(cpuUsage => convertAgentWithCPUUsageFromDTO(cpuUsage)) : [],
        memoryUsage: topNAgentsDTO.memoryUsage ? topNAgentsDTO.memoryUsage.map(memoryUsage => convertAgentWithMemoryUsageFromDTO(memoryUsage)) : [],
    }
}

const convertAgentWithCPUUsageFromDTO = (dto: AgentWithCPUUsageDTO): AgentWithCPUUsage => {
    return {
        id: dto.id ?? "",
        name: dto.name ?? "",
        cpuUsage: dto.cpuUsage ?? 0,
    }
}

const convertAgentWithMemoryUsageFromDTO = (dto: AgentWithMemoryUsageDTO): AgentWithMemoryUsage => {
    return {
        id: dto.id ?? "",
        name: dto.name ?? "",
        memoryUsage: dto.memoryUsage ?? 0,
    }
}
export const convertOverviewFromDTO = (overviewDTO: OverviewDTO): Overview => {
    return {
        summary: overviewDTO.summary ? convertAgentsSummaryFromDTO(overviewDTO.summary) : {
            averageCPUUsage: 0,
            averageMemoryUsage: 0,
            cpuUsageDistribution: {
                high: { percent: 0, count: 0 },
                medium: { percent: 0, count: 0 },
                low: { percent: 0, count: 0 },
            },
            memoryUsageDistribution: {
                high: { percent: 0, count: 0 },
                medium: { percent: 0, count: 0 },
                low: { percent: 0, count: 0 },
            },
            diskUsageDistribution: {
                high: { percent: 0, count: 0 },
                medium: { percent: 0, count: 0 },
                low: { percent: 0, count: 0 },
            },
            onlineAgents: 0,
            totalAgents: 0,
        },
        topN: overviewDTO.topN ? convertTopNAgentsFromDTO(overviewDTO.topN) : {
            cpuUsage: [],
            memoryUsage: [],
        },
    }
}