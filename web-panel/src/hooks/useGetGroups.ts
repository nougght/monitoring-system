import {useQuery} from "@tanstack/react-query";
import {getAllAgentGroups} from "../api/client/monitoringServerAPI.ts";
import {convertGroupFromDTO} from "../api/mapper.ts";


export function useGroups() {
    return useQuery({
        queryKey: ['groups'],
        queryFn: async () => (await getAllAgentGroups()).map(convertGroupFromDTO),
        refetchOnWindowFocus: true,
        staleTime: 60000,
    });
}