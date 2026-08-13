import './App.css'
import { BrowserRouter, Navigate, Route, Routes } from 'react-router-dom'
import { AgentsPage } from './pages/AgentsPage'
import { NewAgentPage } from './pages/NewAgentPage'

function App() {

    return (
        <div className='app'>
            <BrowserRouter>
                <Routes>
                    <Route path="/" element={<Navigate to="/agents" replace />} />
                    <Route path="/agents" element={<AgentsPage />} />
                    <Route path="/agents/new" element={<NewAgentPage />} />
                    {/* <Route path="*" element={<NotFoundPage />} /> */}
                </Routes>
            </BrowserRouter>
        </div>
    )
}

export default App
