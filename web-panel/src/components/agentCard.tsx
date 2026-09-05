import type { Agent } from "../domain/agent"


const onlineInd = "🟢"
const offlineInd = "🔴"

export const AgentCard = ({ agent, onClick }: { agent: Agent, onClick: (id: string) => void }) => {
    return (
        <div className="agent-card" onClick={() => onClick(agent.id)}>
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