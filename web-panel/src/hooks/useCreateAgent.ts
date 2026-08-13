import {
    useQuery,
    useMutation,
    useQueryClient,
    QueryClient,
    QueryClientProvider,
} from '@tanstack/react-query'
import type { AgentConfigBody, CreateAgentBody } from '../api/models'
import { postAgents, postAgentsAgentIDSetupconfig,  } from '../api/client/monitoringServerAPI'
import type {postAgentsResponse }from '../api/client/monitoringServerAPI'

export const useCreateAgents = () => {
    const queryClient = useQueryClient()

    return useMutation({
        mutationFn: async (dto: CreateAgentBody) => { 
            var resp = await postAgents(dto)
            if (resp.status != 200) {
                throw new Error(`error: ${resp.data}`)
            }
            return resp
        },

        onSuccess: (resp) => {
            queryClient.invalidateQueries({ queryKey: ['post-agents'] })
        },

        onError: (error) => {
            console.error('failed to post agent', error)
        }
    })
}

