import { Link } from "react-router-dom"
import { useAgents } from "../hooks/useGetAgents";
import { useEffect, useState } from "react";
import { AgentCard } from "../components/agentCard";



export const AgentsPage = () => {
    const [warning, setWarning] = useState<string | null>()
    const {
        data: agents,
        isPending: isAgentsPending,
        isError: _isAgentsError,
        error: _agentsError,
        isFetching: _isAgentsFetching,
    } = useAgents();


    useEffect(() => {
        if (agents?.error != null) {
            setWarning(`ошибка:${agents?.error.status} ${agents?.error.message}`)
        }
    }, [agents]);

    if (isAgentsPending) {
        return <div>Загрузка...</div>;
    }
    return (
        <div className='agentsPage'>
            <h1>Агенты</h1>
            <main>
                <div className="agent-cards-container">
                    {agents?.agents != null && agents?.agents.length > 0 &&
                        agents?.agents?.filter((a) => a.status != null).map((agent) => 
                            // {agent.status != null &&
                            <div key={agent.id}>
                                <AgentCard agent={agent} />
                            </div>
                            // }
                        )
                    }
                </div>
                <div className='bottomArea'>
                    <button className='addAgentButton'>
                        <Link to="./new">Добавить</Link>
                    </button>
                </div>
            </main>
            {
                warning != null &&
                <div>
                    <p>{warning}</p>
                </div>
            }
        </div>
    )
}