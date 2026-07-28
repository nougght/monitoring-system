import type { Metrics as DtoMetrics } from "./models"
import type { Metrics } from "../domain/metrics"

export const metricsFromDto = (dto: DtoMetrics): Metrics => {
    return {
        cpuPercent: dto.cpuPercent,
        focusedWindow: dto.focusedWindow,
        memoryUsed: dto.memoryUsed,
        diskUsage: dto.diskUsage,
        uploadMbps: dto.uploadMbps,
        downloadMbps: dto.downloadMbps,
        timestamp: new Date(dto.timestamp ?? new Date().toISOString()),
    }
}