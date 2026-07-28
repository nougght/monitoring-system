import type { Process } from "./process"

export interface Metrics {
    focusedWindow?: string
    cpuPercent?: number
    memoryUsed?: number
    diskUsage?: { [deviceName: string]: number }
    uploadMbps?: number
    downloadMbps?: number
    processList?: Process[]
    timestamp: Date
}

export const EMPTY_FOCUSED_WINDOW = "EMPTY_FOCUSED_WINDOW"