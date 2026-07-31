import { useState } from "react"
import { useCreateAgents } from "../hooks/useCreateAgent";
import type { postAgentsResponse, postAgentsResponseSuccess } from "../api/client/monitoringServerAPI";


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
    const [key, setKey] = useState<string | null>(null)
    const [name, setName] = useState<string>(Math.random().toString(36).substring(2))
    const [description, setDescription] = useState<string | null>(null)
    const [warning, setWarning] = useState<string | null>()

    const { mutate, isPending, isError, isSuccess } = useCreateAgents()


    const handleCreate = () => {
        mutate(
            {
                name: name,
                description: description ?? undefined
            },
            {
                // onSuccess/onError можно передать прямо в mutate
                // они сработают после глобальных из useMutation
                onSuccess: (resp) => {
                    if (resp.status == 200) {
                        setKey(resp.data.enrollmentKey ?? null)
                    } else {
                        setWarning("ошибка")
                    }
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
                    <br/>
                    <label htmlFor="description">Описание</label>
                    <input type="text" id="name" name="name" value={description ?? ""} onChange={e => setName(e.target.value)} />
                    <br/>
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
                        <p>{isSuccess ? key : "....."}</p>
                    </div>
                }
            </div>
            <p>{warning}</p>
        </div>
    )
}