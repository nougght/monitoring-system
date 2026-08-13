
import { useQuery } from '@tanstack/react-query';
import { getAgents as getAgentsDTO } from '../api/client/monitoringServerAPI';
import { convertAgentFromDTO } from '../api/mapper';
import type { Agent } from '../domain/agent';


interface Error {
    status: number
    message: string
}

interface UseAgentsResult {
    agents: Agent[] | undefined
    error: Error | null
}

const getAgents = async (): Promise<UseAgentsResult> => {
    const resp = await getAgentsDTO();
    return resp.status == 200 ?
    {
        agents: resp.data.map(convertAgentFromDTO),
        error: null
    } : 
    {
        agents: undefined,
        error: {
            status: resp.status,
            message: JSON.stringify(resp.data)
        }
    }
}


export function useAgents() {
    return useQuery({
        queryKey: ['agents'],
        queryFn: () => getAgents(),
        refetchOnWindowFocus: true,
        staleTime: 60000,
        retry: 2,
    });
}