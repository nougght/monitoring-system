
import { useQuery } from '@tanstack/react-query';
import { getAgentSpecs, getAllAgents as getAllAgents } from '../api/client/monitoringServerAPI';
import type {Error} from './common'
import type { AgentSpecs } from '../api/models';
import type {Specs} from '../../../shared/ui/src/domain/specs';
import { convertSpecsFromDTO } from '../api/mapper';

interface getSpecsResult {
    specs?: Specs
    error: Error | null
}

const getSpecs = async (id: string): Promise<getSpecsResult> => {
    const resp = await getAgentSpecs(id);
    return resp.status == 200 ?
    {
        specs: convertSpecsFromDTO(resp.data),
        error: null
    } : 
    {
        error: {
            status: resp.status,
            message: JSON.stringify(resp.data)
        }
    }
}


export function useSpecs(id: string) {
    return useQuery({
        queryKey: ['specs'],
        queryFn: () => getSpecs(id),
        retry: 2,
    });
}