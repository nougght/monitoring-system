
interface ProcessData{
    pid: number;
    name: string;
    parentPid: number | null;
    cpuPercent: number | null;
    memoryUsed: number | null;
}
const ProcessRow = ({ pid, name, parentPid, cpuPercent, memoryUsed }: ProcessData) => {


    return(
        <div style={{ display: "flex", flexDirection: "row", justifyContent: "space-between" }}>
            <div>{pid}</div>
            <div>{name}</div>
            <div>{parentPid ?? "---"}</div>
            <div>{cpuPercent ?? "NO DATA"}</div>
            <div>{memoryUsed ?? "NO DATA"}</div>
        </div>
    )
}


interface ProcessTableProps {
    processes: ProcessData[];
}
export const ProcessTable = ({ processes }: ProcessTableProps) => {
    return(
        <div>
            {processes.sort((a, b) => (b?.cpuPercent ?? 0) - (a?.cpuPercent ?? 0)).map((process) => <ProcessRow key={process.pid} {...process} />)}
        </div>
    )
}