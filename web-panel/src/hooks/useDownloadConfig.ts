import {useMutation, } from "@tanstack/react-query";
import type { AgentConfigBody } from "../api/models";
import {downloadAgentFiles} from "../api/client/monitoringServerAPI.ts";

interface SetupConfigVariables {
    agentID: string
    dto: AgentConfigBody
}
// async function downloadConfig(vars: SetupConfigVariables): Promise<Error | void> {
//     const res = await fetch(`/api/v1/agents/${vars.agentID}/setupconfig`, {
//         method: 'POST',
//         body: JSON.stringify(vars.dto)
//     });
//     if (res.status != 200) {
//         return {
//             status: res.status,
//             message: await res.text()
//         }
//     }
//     const blob = await res.blob();
//     const url = URL.createObjectURL(blob);
//     const a = document.createElement('a');
//     a.href = url;
//     a.download = 'agent.zip';
//     a.click();
//     URL.revokeObjectURL(url);
// }

export const useDownloadConfig = () => {
    // const queryClient = useQueryClient()
    return useMutation({
        mutationFn: async (vars: SetupConfigVariables) => {
            const blob = await downloadAgentFiles(vars.agentID, vars.dto)
            const url = URL.createObjectURL(blob);
            const a = document.createElement('a');
            a.href = url;
            a.download = 'agent.zip';
            a.click();
            URL.revokeObjectURL(url);
        },

        onSuccess: () => {
            // queryClient.invalidateQueries({ queryKey: ['post-agents-config'] })
        },

        onError: (error) => {
            console.error('failed to post agent config', error.message)
        }
    })
}