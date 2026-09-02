import './App.css'
import { BrowserRouter, Navigate, Route, Routes, Outlet, useNavigate, useLocation } from 'react-router-dom'
import { AgentsPage } from './pages/AgentsPage'
import { AgentPage } from './pages/AgentPage'
import { NewAgentPage } from './pages/NewAgentPage'
import { SideBar, type SideBarData, type NavItem } from './components/sideBar'
import { useEffect, useState } from 'react'
import dashIcon from "./assets/dashboard.svg"
import agentsIcon from "./assets/cpu.svg"

const sideBarData: SideBarData = {
    iconSrc: "",
    title: "Vigil",
    items: [
        { id: "dashboard", title: "Обзор", iconSrc: dashIcon, path: "/dashboard", countLabel: 0 },
        { id: "agents", title: "Агенты", iconSrc: agentsIcon, path: "/agents", countLabel: 0 },
    ]

}
const AppLayout = () => {
    const navigate = useNavigate()
    let location = useLocation()

    useEffect(
        () => {

        },
        [location]
    )
    return (
        <div style={{ display: `flex`, flexDirection: `row`, height: `100%`, alignItems: `stretch` }}>
            <SideBar
                data={sideBarData}
            />
            <div style={{ flexGrow: 1, overflow: `auto` }}><Outlet /></div>

        </div>
    )
}

function App() {

    return (
        <div className='app'>
            <BrowserRouter>
                <Routes>
                    <Route element={<AppLayout />}>
                        <Route path="/" element={<Navigate to="/agents" replace />} />
                        <Route path="/dashboard" element={<p>not implemented</p>} />
                        <Route path="/agents" element={<AgentsPage />} />
                        <Route path="/agents/:id" element={<AgentPage />} />
                        <Route path="/agents/new" element={<NewAgentPage />} />
                    </Route>
                    {/* <Route path="*" element={<NotFoundPage />} /> */}
                </Routes>
            </BrowserRouter>
        </div>
    )
}

export default App

