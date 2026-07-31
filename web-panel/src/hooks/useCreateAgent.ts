import {
  useQuery,
  useMutation,
  useQueryClient,
  QueryClient,
  QueryClientProvider,
} from '@tanstack/react-query'
import type { CreateAgentBody } from '../api/models'
import { postAgents } from '../api/client/monitoringServerAPI'


export const useCreateAgents = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (dto: CreateAgentBody) => postAgents(dto),

    onSuccess: (resp) => {
      // Инвалидируем кэш — React Query сам перезапросит список
      queryClient.invalidateQueries({ queryKey: ['post-agents'] })
    },

    onError: (error) => {
      console.error('failed to post agent', error)
    }
  })
}