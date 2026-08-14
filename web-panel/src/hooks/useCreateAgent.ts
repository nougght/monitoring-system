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
            var resp = await createAgent(dto)
            if (resp.status != 200) {
                throw new Error(`error: ${resp.data}`)
            }
            return resp
        },

        onSuccess: (_resp) => {
            queryClient.invalidateQueries({ queryKey: ['post-agents'] })
        },

        onError: (error) => {
            console.error('failed to post agent', error)
        }
    })
}

