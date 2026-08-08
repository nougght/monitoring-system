import { useState } from "react"
import { useCreateAgents } from "../hooks/useCreateAgent";
import type { postAgentsResponse, postAgentsResponseSuccess } from "../api/client/monitoringServerAPI";
import type { CreateAgentResponse } from "../api/models";
import { useDownloadConfig } from "../hooks/useDownloadConfig";


interface NewAgentProps {
    createdAt?: string;
    description?: string;
    enrollmentKey?: string;
    id?: string;
    isOnline?: boolean;
    lastSeenAt?: string;
    name?: string;
}

export const NewAgentPage = (_props: NewAgentProps) => {
    const [resp, setResp] = useState<CreateAgentResponse | null>(null)
    const [name, setName] = useState<string>(Math.random().toString(36).substring(2))
    const [description, setDescription] = useState<string | null>(null)
    const [warning, setWarning] = useState<string | null>()
    const [info, setInfo] = useState<string | null>()

    const { mutate, isPending, isError, isSuccess } = useCreateAgents()
    const { mutate: downloadConfig, isPending: isPendingConfig, isError: isErrorConfig, isSuccess: isSuccessError } = useDownloadConfig()


    const handleCreate = () => {
        mutate(
            {
                name: name,
                description: description ?? undefined
            },
            {
                onSuccess: (resp) => {
                    if (resp.status == 200) {
                        setResp(resp.data)
                    } else {
                        setWarning("ошибка")
                    }
                },
                onError: (error) => {
                    setWarning(`ошибка:${error}`)
                }
            }
        )
    }

    const handleDownloadInstaller = () => {
        setInfo("not implemented")
        setTimeout(() => { setInfo(null) }, 1500)
    }
    const handleDownloadConfig = () => {
        downloadConfig(
            {
                agentID: resp?.id ?? "",
                dto: {
                    enrollmentKey: resp?.enrollmentKey ?? undefined
                }
            },
            {
                onError:(error)=> {
                    setWarning(`ошибка:${error}`)
                }
            }
        )
    }

    return (
        <div>
            <div>
                <h1>Создание агента</h1>
                <div>
                    <label htmlFor="name">Название</label>
                    <input type="text" id="name" name="name" value={name}
                        required onChange={e => setName(e.target.value)} />
                    <br />
                    <label htmlFor="description">Описание</label>
                    <input type="text" id="name" name="name" value={description ?? ""} onChange={e => setName(e.target.value)} />
                    <br />
                    <button onClick={handleCreate}>Создать</button>
                </div>
            </div>

            <div>
                {
                    isSuccess &&
                    <div className="keyContainer">
                        <h3>
                            Ключ подключения агента
                        </h3>
                        <p>{isSuccess ? resp?.enrollmentKey : "....."}</p>
                        <div>
                            <button onClick={handleDownloadInstaller}>Скачать установщик</button>
                            <button onClick={handleDownloadConfig}>Скачать конфиг</button>
                        </div>
                    </div>
                }
            </div>
            {
                warning != null &&
                <div>
                <p>{warning}</p>
                </div>
            }
            {
                info != null &&
                <div className="infoMessage">
                    <p>{info}</p>
                </div>
            }
        </div>
    )
}