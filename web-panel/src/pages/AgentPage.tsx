
import TabBar from '../../../shared/ui/src/components/TabBar';
import Specifications from '../../../shared/ui/src/components/Specifications';
import { useSpecs } from "../hooks/useGetSpecs";
import { useAgent } from "../hooks/useGetAgent";
import { useEffect, useState } from 'react';
import type { Agent } from '../domain/agent';
import { useLocation, useParams } from 'react-router-dom';

interface Tab {
    text: string;
    content: React.ReactNode;
}

export const AgentPage = () => {
    const { state } = useLocation() as { state: Agent | null }

    const [agent, setAgent] = useState<Agent | null>(state)

    const { id } = useParams()
    const [activeTab, setActiveTab] = useState(0)

    const [warning, setWarning] = useState<string | null>()
    const {
        data: specs,
        isPending: isSpecsPending,
        isError: _isSpecsError,
        error: _specsError,
        isFetching: _isSpecsFetching,
    } = useSpecs(id ?? "");

    const {
        data: agentResp,
        isPending: isAgentPending,
        isError: _isAgentError,
        error: _agentError,
        isFetching: _isAgentFetching,
    } = state == null && id != undefined ? useAgent(id) : {}

    useEffect(() => {
        if (specs?.error != null) {
            console.error(specs?.error)
            setWarning(`ошибка:${specs?.error.status} ${specs?.error.message}`)
        }
    }, [specs]);


    useEffect(() => {
        if (agentResp?.error != null) {
            console.error(agentResp?.error)
            setWarning(`ошибка:${agentResp?.error.status} ${agentResp?.error.message}`)
        }
        setAgent(agentResp?.agent!)
    }, [agentResp]);


    const tabs: Tab[] = [
        //TODO: full overview page 
        {
            text: "Обзор",
            content: isAgentPending ? <div>Загрузка...</div> :
                agent != null ?
                <div>
                    <p>{`Имя хоста: ${specs?.specs?.host?.hostName ?? "NO DATA"}`}</p>
                    <p>{`Идентификатор агента: ${agent.id}`}</p>
                    <img src={`http://127.0.0.1:8091/api/v1/agents/${agent.id}/frames`} height="200px"/>
                </div> :
                agentResp?.error?.status == 404 && <div>Агент не найден</div>
        },
        {
            text: "Характеристики",
            content: isSpecsPending ? <div>Загрузка...</div> :
                specs?.specs != null ?
                    <div>
                        <Specifications specs={specs.specs} />
                    </div> :
                    specs?.error?.status == 404 && <div>Характеристики не найдены</div>
        }
    ]
    return (
        <div className="agent-page">
            <TabBar tabs={tabs.map((tab) => tab.text)} onSwitch={setActiveTab} activeTab={activeTab} />
            <div>
                {tabs[activeTab].content}
            </div>
            {
                warning != null &&
                <div>
                    <p>{warning}</p>
                </div>
            }
        </div>
    )
}