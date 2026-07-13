export interface Process {
    pid: number;
    name: string;
    parentPid: number | null;
    cpuPercent: number | null;
    memoryUsed: number | null;
}