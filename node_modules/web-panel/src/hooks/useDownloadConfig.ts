import { useMutation, useQueryClient } from "@tanstack/react-query";
import type { AgentConfigBody } from "../api/models";

interface SetupConfigVariables {
    agentID: string
    dto: AgentConfigBody
}
interface Error {
    status: number
    message: string
}

async function downloadConfig(vars: SetupConfigVariables): Promise<Error | void> {
    const res = await fetch(`http://127.0.0.1:8091/api/v1/agents/${vars.agentID}/setupconfig`, {
        method: 'POST',
        body: JSON.stringify(vars.dto)
    });
    if (res.status != 200) {
        return {
            status: res.status,
            message: await res.text()
        }
    }
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
        mutationFn: async (vars: SetupConfigVariables) => {
            var error = await downloadConfig(vars)
            if (error) {
                throw new Error(`status: ${error.status} message: ${error.message}`)
            }
        },

        onSuccess: (err) => {
            queryClient.invalidateQueries({ queryKey: ['post-agents-config'] })
        },

        onError: (error) => {
            console.error('failed to post agent config', error.message)
        }
    })
}