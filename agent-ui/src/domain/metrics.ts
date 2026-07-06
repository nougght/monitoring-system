export interface Metrics {
    focusedWindow?: string
    cpuPercent?: number
    memoryUsed?: number
    timestamp: Date
}

export const EMPTY_FOCUSED_WINDOW = "EMPTY_FOCUSED_WINDOW"