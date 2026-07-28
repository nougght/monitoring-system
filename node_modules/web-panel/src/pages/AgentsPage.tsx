import { Link } from "react-router-dom"



export const AgentsPage = () => {
    return (
        <div className='agentsPage'>
            <h1>Агенты</h1>
            <main>
                <div></div>
                <div className='bottomArea'>
                    <button className='addAgentButton'>
                        <Link to="./new">Добавить</Link>
                    </button>
                </div>
            </main>
        </div>
    )
}