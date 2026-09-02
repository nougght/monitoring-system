
import { useQuery } from '@tanstack/react-query';
import { getAgentByID } from '../api/client/monitoringServerAPI';
import { convertAgentFromDTO } from '../api/mapper';
import type { Agent } from '../domain/agent';
import type {Error} from './common'

interface UseAgentResult {
    agent?: Agent
    error: Error | null
}

const getAgent = async (id: string): Promise<UseAgentResult> => {
    const resp = await getAgentByID(id);
    return resp.status == 200 ?
    {
        agent: convertAgentFromDTO(resp.data),
        error: null
    } : 
    {
        error: {
            status: resp.status,
            message: JSON.stringify(resp.data)
        }
    }
}


export function useAgent(id: string) {
    return useQuery({
        queryKey: ['agent'],
        queryFn: () => getAgent(id),
        staleTime: 60000,
        retry: 2,
    });
}