

import type { Agent as AgentDTO} from "./models";
import type { Agent } from "../domain/agent";

export const convertAgentFromDTO = (agentDTO: AgentDTO): Agent => {
    return {
        id: agentDTO.id ?? '',
        name: agentDTO.name ?? '',
        description: agentDTO.description,
        createdAt: agentDTO.createdAt ?? '',
        lastSeenAt: agentDTO.lastSeenAt ? new Date(agentDTO.lastSeenAt) : undefined,
        isOnline: agentDTO.isOnline ?? false
    }
}