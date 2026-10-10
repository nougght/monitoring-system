import {useQuery} from "@tanstack/react-query";
import {convertOverviewFromDTO} from "../api/mapper";
import {getFleetOverview} from "../api/client/monitoringServerAPI";


export function useOverview() {
    return useQuery({
        queryKey: ['overview'],
        queryFn: async () => convertOverviewFromDTO(await getFleetOverview()),

        refetchOnWindowFocus: true,
        staleTime: 60000,
    });
}