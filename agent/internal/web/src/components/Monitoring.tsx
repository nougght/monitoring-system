import { EMPTY_FOCUSED_WINDOW } from "../domain/metrics";
import { getGradientColor } from "../util/gradientColor";

interface MonitoringProps {
    focusedWindow?: string
    cpuPercent?: number
}

const Monitoring = ({ focusedWindow, cpuPercent }: MonitoringProps) => {
    return (
        <section>
            <h2>Active window</h2>
            <p>{focusedWindow == null ? "No data" : 
                focusedWindow == EMPTY_FOCUSED_WINDOW ? "No active window" : 
                focusedWindow}
            </p>
            <h2>CPU usage</h2>
            <p style={{ color: cpuPercent ? getGradientColor(["#00f000", "#ff0000"], Math.round(cpuPercent)) : "black" }}>{cpuPercent == null ? "No data" : Math.round(cpuPercent * 100) / 100 + "%"}</p>
        </section>
    )
}

export default Monitoring;