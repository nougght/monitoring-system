
import { useQuery } from '@tanstack/react-query';
import { getAllAgents as getAllAgents } from '../api/client/monitoringServerAPI';
import { convertAgentFromDTO } from '../api/mapper';


export function useAgents() {
    return useQuery({
        queryKey: ['agents'],
        queryFn: async () => (await getAllAgents()).map(convertAgentFromDTO),
        refetchOnWindowFocus: true,
        staleTime: 60000,
    });
}