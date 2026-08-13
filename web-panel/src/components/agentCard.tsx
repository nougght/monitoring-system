import type { Agent } from "../domain/agent"


const onlineInd = "🟢"
const offlineInd = "🔴"

export const AgentCard = ({ agent }: { agent: Agent }) => {
    return (
        <div className="agent-card">
            <h2>{agent.name}</h2>
            <p>{agent.description ?? "no description"}</p>
            {
                agent.status != null &&
                <p>{`Status: ${agent.status}`}</p>
            }
            <div>
                <span>{agent.isOnline ? onlineInd : offlineInd}</span>
                <span>{agent.isOnline ? "online" : "offline"}</span>
            </div>

        </div>
    )
}