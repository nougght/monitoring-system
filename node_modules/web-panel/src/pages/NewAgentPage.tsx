import { useState } from "react"


interface NewAgentProps{

}

export const NewAgentPage = (_props: NewAgentProps) => {
    const [key, _setKey] = useState<string | null>(null)

    return (
        <div>
            <div className="keyContainer">
                <h3>
                    Ключ подключения агента
                </h3>
                <p>{key ?? "....."}</p>
            </div>
        </div>
    )
}