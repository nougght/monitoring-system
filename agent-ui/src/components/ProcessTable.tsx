import { convertBytesToMB } from "../util/units";

interface ProcessData {
    pid: number;
    name: string;
    parentPid: number | null;
    cpuPercent: number | null;
    memoryUsed: number | null;
}
const ProcessRow = (process: ProcessData) => {


    return (
        <div style={{ display: "flex", flexDirection: "row", justifyContent: "space-between" }}>
            <div>{process.pid}</div>
            <div>{process.name}</div>
            <div>{process.parentPid ?? "---"}</div>
            <div>{process.cpuPercent ?? "NO DATA"}</div>
            <div>{process.memoryUsed == null ? "NO DATA" : convertBytesToMB(process.memoryUsed).toFixed(2)} MB</div>
        </div>
    )
}


interface ProcessTableProps {
    processes: ProcessData[];
}
export const ProcessTable = ({ processes }: ProcessTableProps) => {
    return (
        <div>
            {processes.sort((a, b) => (b?.cpuPercent ?? 0) - (a?.cpuPercent ?? 0)).map((proc) =>
                <ProcessRow key={proc.pid} {...proc} />)}
        </div>
    )
}