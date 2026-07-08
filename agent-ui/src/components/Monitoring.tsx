import { EMPTY_FOCUSED_WINDOW, type Metrics } from "../domain/metrics";
import type { Specs } from "../domain/specs";
import { getGradientColor } from "../util/gradientColor";
import { convertBytesToGB } from "../util/units";

interface MonitoringProps {
    specs?: Specs
    metrics?: Metrics
}

const Monitoring = ({ specs, metrics }: MonitoringProps) => {
    return (
        <section>
            <h2>Active window</h2>
            <p>{metrics?.focusedWindow == null ? "No data" :
                metrics?.focusedWindow == EMPTY_FOCUSED_WINDOW ? "No active window" :
                    metrics?.focusedWindow}
            </p>
            <h2>CPU usage</h2>
            <p style={{
                color: metrics?.cpuPercent ? getGradientColor(["#4cd485", "#e0cb51", "#d44c4c"],
                    Math.round(metrics?.cpuPercent)) : "black"
            }}>
                {metrics?.cpuPercent == null ? "No data" :
                    metrics?.cpuPercent.toFixed(2) + "%"}
            </p>
            <h2>Memory usage</h2>
            <p>{metrics?.memoryUsed == null ? "No data" :
                convertBytesToGB(metrics?.memoryUsed).toFixed(2)} / {convertBytesToGB(specs?.memory?.total).toFixed(2)} GB <span
                    style={{
                        color: metrics?.memoryUsed != null && specs?.memory?.total != null ?
                            getGradientColor(["#4cd485", "#e0cb51", "#d44c4c"], Math.round((convertBytesToGB(metrics?.memoryUsed) /
                                convertBytesToGB(specs?.memory?.total)) * 100)) : "black"
                    }}>
                    ({Math.round((convertBytesToGB(metrics?.memoryUsed) / convertBytesToGB(specs?.memory?.total)) * 100)}%)
                </span>
            </p>
            <h2>Disk usage</h2>
            <p>
                {specs?.disk?.map((disk) => {
                    return (
                        <div key={disk.device}>
                            <p>
                                {disk.device}: {convertBytesToGB(metrics?.diskUsage?.[disk.device] ?? 0).toFixed(2)} /
                                {convertBytesToGB(disk.total).toFixed(2)} GB <span style={{ color: getGradientColor(["#4cd485", "#e0cb51", "#d44c4c"], Math.round((metrics?.diskUsage?.[disk.device] ?? 0) / disk.total * 100)) }}>
                                    ({Math.round((metrics?.diskUsage?.[disk.device] ?? 0) / disk.total * 100)}%)
                                </span>
                            </p>
                        </div>
                    )
                })}
            </p>
            <h2>Net usage</h2>
            <p>
                Up: {metrics?.uploadMbps.toFixed(2)} | Down: {metrics?.downloadMbps.toFixed(2)} Mbit/s
            </p>
        </section>
    )
}

export default Monitoring;