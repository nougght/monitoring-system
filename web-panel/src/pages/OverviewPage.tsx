import {useEffect, useState} from "react"
import {Link} from "react-router-dom";
import {useOverview} from "../hooks/useOverview";
import SimplePieChart from "../components/pieChart";
import {CountDistToPieList, type Overview} from "../domain/overview";
import {AgentCount} from "../components/agentCount";
import {LoadBar} from "../components/loadBar";
import {AvgCPU} from "../components/avgCPU";
import {AvgMem} from "../components/avgMem";
import styles from "./overviewPage.module.css"
import commonStyles from "../common.module.css"

export const OverviewPage = ({overviewProp}: { overviewProp: Overview | undefined }) => {
    const [warning, setWarning] = useState<string | null>()
    const [info, _setInfo] = useState<string | null>()


    const {
        data: overview,
        isPending: isOverviewLoading,
        isError: _isOverviewError,
        error: _overviewError,
        isFetching: _isOverviewFetching,
    } = useOverview();

    useEffect(() => {
        if (overviewProp != undefined) {
            setWarning(null)
            overview!.summary = overviewProp.summary
            overview!.topN = overviewProp.topN

        }
    }, [overviewProp]);

    useEffect(() => {
        if (_isOverviewError) {
            setWarning(`ошибка:${_overviewError.message}`)
        }
        // if (overview?.overview != null) {
        //     setInfo("Данные успешно загружены")

    }, [overview]);

    if (isOverviewLoading) {
        return <div>Загрузка...</div>;
    }

    return (

        <div>
            <div>
                <h1>Обзор</h1>
                {overview?.summary != null && (

                    <div>
                        <div className={styles.topCards}>
                            <AgentCount
                                online={overview?.summary?.onlineAgents}
                                count={overview?.summary?.totalAgents}
                            />
                            <AvgCPU
                                value={overview?.summary?.averageCPUUsage}
                            />
                            <AvgMem
                                value={overview?.summary?.averageMemoryUsage}
                            />
                            {/* <label>Среднее использование CPU </label>
                            <span>{overview?.overview?.summary?.averageCPUUsage.toFixed(2)}%</span>
                            <br />
                            <label>Среднее использование памяти </label>
                            <span>{overview?.overview?.summary?.averageMemoryUsage.toFixed(2)}%</span> */}
                        </div>
                        <div>
                            {/* <label>Распределение использования CPU</label>
                            <br />
                            <label>Высокое({overview?.overview?.summary?.cpuUsageDistribution?.high?.percent}%) </label>
                            <span>{overview?.overview?.summary?.cpuUsageDistribution?.high?.count}</span>
                            <br />
                            <label>Среднее({overview?.overview?.summary?.cpuUsageDistribution?.medium?.percent}%) </label>
                            <span>{overview?.overview?.summary?.cpuUsageDistribution?.medium?.count}</span>
                            <br />
                            <label>Низкое({overview?.overview?.summary?.cpuUsageDistribution?.low?.percent}%) </label>
                            <span>{overview?.overview?.summary?.cpuUsageDistribution?.low?.count}</span>
                            <br />
                            <label>Распределение использования памяти</label>
                            <br />
                            <label>Высокое({overview?.overview?.summary?.memoryUsageDistribution?.high?.percent}%) </label>
                            <span>{overview?.overview?.summary?.memoryUsageDistribution?.high?.count}</span>
                            <br />
                            <label>Среднее({overview?.overview?.summary?.memoryUsageDistribution?.medium?.percent}%) </label>
                            <span>{overview?.overview?.summary?.memoryUsageDistribution?.medium?.count}</span>
                            <br />
                            <label>Низкое({overview?.overview?.summary?.memoryUsageDistribution?.low?.percent}%) </label>
                            <span>{overview?.overview?.summary?.memoryUsageDistribution?.low?.count}</span>
                            <br />
                            <label>Распределение использования диска</label>
                            <br />
                            <label>Высокое({overview?.overview?.summary?.diskUsageDistribution?.high?.percent}%) </label>
                            <span>{overview?.overview?.summary?.diskUsageDistribution?.high?.count}</span>
                            <br />
                            <label>Среднее({overview?.overview?.summary?.diskUsageDistribution?.medium?.percent}%) </label>
                            <span>{overview?.overview?.summary?.diskUsageDistribution?.medium?.count}</span>
                            <br />
                            <label>Низкое({overview?.overview?.summary?.diskUsageDistribution?.low?.percent}%) </label>
                            <span>{overview?.overview?.summary?.diskUsageDistribution?.low?.count}</span> */}

                            <div className={styles.distPies}>

                                <SimplePieChart
                                    Title="Использование CPU"
                                    PieData={CountDistToPieList(overview.summary.cpuUsageDistribution)}
                                />

                                <SimplePieChart
                                    Title="Использование памяти"
                                    PieData={CountDistToPieList(overview.summary.memoryUsageDistribution)}
                                />

                                <SimplePieChart
                                    Title="Использование диска"
                                    PieData={CountDistToPieList(overview.summary.diskUsageDistribution)}
                                />
                            </div>
                        </div>

                        <div className={styles.topN}>
                            <div className={styles.topNSection}>
                                <h3>Топ по использованию CPU</h3>
                                <table>
                                    {/* <thead>
                                        <tr>
                                            <th>Агент</th>
                                            <th>Использование CPU</th>
                                        </tr>
                                    </thead> */}
                                    <tbody>
                                    {overview.topN?.cpuUsage?.map((agent) => (
                                        agent.name != "" &&
                                        <tr key={agent.id}>
                                            <td><Link to={`/agents/${agent.id}`}>{agent.name}</Link></td>
                                            <td><LoadBar value={agent.cpuUsage} width={70} height={10}/></td>
                                            <td>{agent.cpuUsage.toFixed(2)}%</td>
                                        </tr>
                                    ))}
                                    </tbody>
                                </table>
                            </div>
                            <div className={styles.topNSection}>
                                <h3>Топ по использованию памяти</h3>
                                <table>
                                    {/* <thead>
                                        <tr>
                                            <th>Агент</th>
                                            <th>Использование памяти</th>
                                        </tr>
                                    </thead> */}
                                    <tbody>
                                    {overview?.topN?.memoryUsage?.map((agent) => (
                                        agent.name != "" &&
                                        <tr key={agent.id}>
                                            <td><Link to={`/agents/${agent.id}`}>{agent.name}</Link></td>
                                            <td><LoadBar value={agent.memoryUsage} width={70} height={10}/></td>
                                            <td>{agent.memoryUsage.toFixed(2)}%</td>
                                        </tr>
                                    ))}
                                    </tbody>
                                </table>
                            </div>
                        </div>
                    </div>
                )}
            </div>

            {
                warning != null &&
                <div>
                    <p>{warning}</p>
                </div>
            }
            {
                info != null &&
                <div className={commonStyles.infoMessage}>
                    <p>{info}</p>
                </div>
            }
        </div>
    )
}