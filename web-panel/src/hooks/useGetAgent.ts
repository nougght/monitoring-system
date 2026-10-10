
import { useQuery } from '@tanstack/react-query';
import { getAgentByID } from '../api/client/monitoringServerAPI';
import { convertAgentFromDTO } from '../api/mapper';


export function useAgent(id: string) {
    return useQuery({
        queryKey: [`agent-${id}`],
        queryFn: async () => convertAgentFromDTO(await getAgentByID(id)),
        staleTime: 60000,
    });
}