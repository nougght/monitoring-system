import {
  useQuery,
  useMutation,
  useQueryClient,
  QueryClient,
  QueryClientProvider,
} from '@tanstack/react-query'
import type { AgentConfigBody, CreateAgentBody } from '../api/models'
import { postAgents, postAgentsAgentIDSetupconfig } from '../api/client/monitoringServerAPI'


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


interface SetupConfigVariables {
  agentID: string
  dto: AgentConfigBody
}
async function downloadConfig(vars: SetupConfigVariables):Promise<void> {
  const res = await fetch(`http://127.0.0.1:8091/api/v1/agents/${vars.agentID}/config`, {
    method: 'POST',
  });
  const blob = await res.blob();
  const url = URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = 'config.yaml';
  a.click();
  URL.revokeObjectURL(url);
}

export const useDownloadConfig = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (vars: SetupConfigVariables) => downloadConfig(vars),

    onSuccess: () => {
      // Инвалидируем кэш — React Query сам перезапросит список
      queryClient.invalidateQueries({ queryKey: ['post-agents-config'] })
    },

    onError: (error) => {
      console.error('failed to post agent config', error)
    }
  })
}