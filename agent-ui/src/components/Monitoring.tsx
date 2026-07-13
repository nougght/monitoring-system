import { EMPTY_FOCUSED_WINDOW, type Metrics } from "../domain/metrics";
import type { Specs } from "../domain/specs";
import { getGradientColor } from "../util/gradientColor";
import { convertBytesToGB } from "../util/units";


import { ProcessTable } from "./ProcessTable";

interface MonitoringProps {
    specs?: Specs
    metrics?: Metrics
}

const Monitoring = ({ specs, metrics }: MonitoringProps) => {
    return (
        <section>
            <h2>Active window</h2>
            <div>{metrics?.focusedWindow == null ? "No data" :
                metrics?.focusedWindow == EMPTY_FOCUSED_WINDOW ? "No active window" :
                    metrics?.focusedWindow}
            </div>
            <h2>CPU usage</h2>
            <div style={{
                color: metrics?.cpuPercent ? getGradientColor(["#4cd485", "#e0cb51", "#d44c4c"],
                    Math.round(metrics?.cpuPercent)) : "black"
            }}>
                {metrics?.cpuPercent == null ? "No data" :
                    metrics?.cpuPercent.toFixed(2) + "%"}
            </div>
            <h2>Memory usage</h2>
            <div>{metrics?.memoryUsed == null ? "No data" :
                convertBytesToGB(metrics?.memoryUsed ?? 0).toFixed(2)} / {convertBytesToGB(specs?.memory?.total ?? 0).toFixed(2)} GB <span
                    style={{
                        color: metrics?.memoryUsed != null && specs?.memory?.total != null ?
                            getGradientColor(["#4cd485", "#e0cb51", "#d44c4c"], Math.round((convertBytesToGB(metrics?.memoryUsed) /
                                convertBytesToGB(specs?.memory?.total)) * 100)) : "black"
                    }}>
                    ({Math.round((convertBytesToGB(metrics?.memoryUsed ?? 0) / convertBytesToGB(specs?.memory?.total ?? 0)) * 100)}%)
                </span>
            </div>
            <h2>Disk usage</h2>
            <div>
                {specs?.disk?.map((disk) => {
                    return (
                        <div key={disk.device}>
                            <p>
                                {disk.device}: {convertBytesToGB(metrics?.diskUsage?.[disk.device ?? ""] ?? 0).toFixed(2)} /
                                {convertBytesToGB(disk.total ?? 0).toFixed(2)} GB <span style={{ color: getGradientColor(["#4cd485", "#e0cb51", "#d44c4c"], Math.round((metrics?.diskUsage?.[disk.device ?? ""] ?? 0) / (disk.total ?? 0) * 100)) }}>
                                    ({Math.round((metrics?.diskUsage?.[disk.device ?? ""] ?? 0) / (disk.total ?? 0) * 100)}%)
                                </span>
                            </p>
                        </div>
                    )
                })}
            </div>
            <h2>Net usage</h2>
            <div>
                Up: {metrics?.uploadMbps?.toFixed(2) ?? "0"} | Down: {metrics?.downloadMbps?.toFixed(2) ?? "0"} Mbit/s
            </div>
            <h2>Processes</h2>
            <ProcessTable
                processes={metrics?.processList ?? []}/>
        </section>
    )
}

export default Monitoring;