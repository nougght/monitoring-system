import {
    useMutation,
    useQueryClient,
} from '@tanstack/react-query'
import type { CreateAgentBody } from '../api/models'
import { createAgent, } from '../api/client/monitoringServerAPI'

export const useCreateAgents = () => {
    const queryClient = useQueryClient()

    return useMutation({
        mutationFn: async (dto: CreateAgentBody) => { 
            return await createAgent(dto)
        },

        onSuccess: (_resp) => {
            queryClient.invalidateQueries({ queryKey: ['agents'] })
        },

        onError: (error) => {
            console.error('failed to post agent', error)
        }
    })
}

