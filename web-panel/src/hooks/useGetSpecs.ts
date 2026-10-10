import {useQuery} from '@tanstack/react-query';
import {getAgentSpecs} from '../api/client/monitoringServerAPI';
import {convertSpecsFromDTO} from '../api/mapper';


export function useSpecs(id: string) {
    return useQuery({
        queryKey: ['specs'],
        queryFn: async () => convertSpecsFromDTO(await getAgentSpecs(id)),
        retry: 2,
    });
}