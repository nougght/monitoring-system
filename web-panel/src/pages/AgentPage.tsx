
import TabBar from '../../../shared/ui/src/components/TabBar';
import Specifications from '../../../shared/ui/src/components/Specifications';
import { useSpecs } from "../hooks/useGetSpecs";
import { useAgent } from "../hooks/useGetAgent";
import { useEffect, useState } from 'react';
import type { Agent } from '../domain/agent';
import { useLocation, useParams } from 'react-router-dom';
import type { Metrics } from '../domain/metrics';
import { convertBytesToGB } from '@monitoring-system/shared/src/util/units';
import { getGradientColor } from '@monitoring-system/shared/src/util/gradientColor';
import type { ActivityUpdate } from '../domain/activity';
import { AgentActivity } from '../components/agentActivity';

interface Tab {
    text: string;
    content: React.ReactNode;
}

export const AgentPage = ({ metricsProp, activity }: { metricsProp?: Metrics, activity?: ActivityUpdate }) => {
    const { state } = useLocation() as { state: Agent | undefined }
    const { id } = useParams()
    const [activeTab, setActiveTab] = useState(0)
    const [warning, setWarning] = useState<string | null>()
    const [metrics, setMetrics] = useState<Metrics | null>()
    const [activityInfo, setActivity] = useState<ActivityUpdate>()
    const {
        data: specs,
        isPending: isSpecsPending,
        isError: _isSpecsError,
        error: _specsError,
        isFetching: _isSpecsFetching,
    } = useSpecs(id ?? "");

    const {
        data: agent,
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
        if (agent?.id != metricsProp?.agentID) {
            console.log(`metrics agent id mismatch: ${metricsProp?.agentID} - ${agent?.id}`)
        } else {
            setMetrics(metricsProp)
        }
    }, [metricsProp, agent?.id])

    useEffect(() => {
        if (agent?.id != activity?.agentID) {
            console.log(`activity agent id mismatch: ${activity?.agentID} - ${agent?.id}`)
        } else {
            setActivity(activity)
        }
    }, [activity, agent?.id])


    const tabs: Tab[] = [
        //TODO: full overview page 
        {
            text: "Обзор",
            content: isAgentPending ? <div>Загрузка...</div> :
                agent != null ?
                    <div>
                        <p>{`Имя хоста: ${specs?.specs?.host?.hostName ?? "NO DATA"}`}</p>
                        <p>{`Идентификатор агента: ${agent.id}`}</p>
                        <section>
                            {/* <h2>Active window</h2>
                            <div>{metrics?.focusedWindow == null ? "No data" :
                                metrics?.focusedWindow == EMPTY_FOCUSED_WINDOW ? "No active window" :
                                    metrics?.focusedWindow}
                            </div> */}
                            <h2>CPU usage</h2>
                            <div style={{
                                color: metrics?.cpuPercent ? getGradientColor(["#4cd485", "#e0cb51", "#d44c4c"],
                                    Math.round(metrics?.cpuPercent)) : "black"
                            }}>
                                {metrics?.cpuPercent == null ? "No data" :
                                    metrics?.cpuPercent.toFixed(2) + "%"}
                            </div>
                            <h2>Memory usage</h2>
                            <div>{metrics?.memoryUsed == null ? "No data" :
                                convertBytesToGB(metrics?.memoryUsed ?? 0).toFixed(2)} / {convertBytesToGB(specs?.specs?.memory?.total ?? 0).toFixed(2)} GB <span
                                    style={{
                                        color: metrics?.memoryUsed != null && specs?.specs?.memory?.total != null ?
                                            getGradientColor(["#4cd485", "#e0cb51", "#d44c4c"], Math.round((convertBytesToGB(metrics?.memoryUsed) /
                                                convertBytesToGB(specs?.specs?.memory?.total)) * 100)) : "black"
                                    }}>
                                    ({Math.round((convertBytesToGB(metrics?.memoryUsed ?? 0) / convertBytesToGB(specs?.specs?.memory?.total ?? 0)) * 100)}%)
                                </span>
                            </div>
                            <h2>Disk usage</h2>
                            <div>
                                {specs?.specs?.disk?.map((disk) => {
                                    return (
                                        <div key={disk.device}>
                                            <p>
                                                {disk.device}: {convertBytesToGB(metrics?.diskUsage?.get(disk.device ?? "") ?? 0).toFixed(2)} /
                                                {convertBytesToGB(disk.total ?? 0).toFixed(2)} GB <span style={{ color: getGradientColor(["#4cd485", "#e0cb51", "#d44c4c"], Math.round((metrics?.diskUsage?.get(disk.device ?? "") ?? 0) / (disk.total ?? 0) * 100)) }}>
                                                    ({Math.round((metrics?.diskUsage?.get(disk.device ?? "") ?? 0) / (disk.total ?? 0) * 100)}%)
                                                </span>
                                            </p>
                                        </div>
                                    )
                                })}
                            </div>
                            <h2>Net usage</h2>
                            <div>
                                Up: {metrics?.uploadMbps?.toFixed(2) ?? "0"} | Down: {metrics?.downloadMbps?.toFixed(2) ?? "0"} Mbit/s
                            </div>
                            <h2>Processes</h2>
                            {/* <ProcessTable
                                processes={metrics?.processList ?? []} /> */}
                        </section>
                    </div> :
                    _isAgentError && <div>Агент не найден</div>
        },
        {
            text: "Активность",
            content: isSpecsPending ? <div>Загрузка...</div> :
                agent != null ?
                    <AgentActivity
                        activity={activityInfo}
                        agentID={agent.id}
                    /> :
                    specs?.error?.status == 404 && <div>Агент не найден</div>
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

