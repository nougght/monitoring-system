import {useMutation, useQueryClient} from "@tanstack/react-query";
import type {CreateAgentGroupBody} from "../api/models";
import {createAgentGroup} from "../api/client/monitoringServerAPI.ts";

export const useCreateGroup = () => {
    const queryClient = useQueryClient()

    return useMutation({
        mutationFn: async (dto: CreateAgentGroupBody) => {
            return await createAgentGroup(dto)
        },

        onSuccess: (_resp) => {
            queryClient.invalidateQueries({queryKey: ['groups']})
        },

    })
}

