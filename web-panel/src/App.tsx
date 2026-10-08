import './App.css'
import {BrowserRouter, Navigate, Outlet, Route, Routes, useLocation} from 'react-router-dom'
import {useMediaQuery} from 'usehooks-ts'
import {AgentsPage} from './pages/AgentsPage'
import {AgentPage} from './pages/AgentPage'
import {NewAgentPage} from './pages/NewAgentPage'
import {SideBar, type SideBarData} from './components/sideBar'
import {useEffect, useRef, useState} from 'react'
import dashIcon from "../public/assets/dashboard.svg"
import agentsIcon from "../public/assets/cpu.svg"
import {
    type AgentDetailedMessage,
    type ClientMessage,
    ClientMessageTypeAgentDetailed,
    FillMetricsFromSeries,
    type Message,
    MessageTypeActivity,
    MessageTypeFleetOverview,
    MessageTypeSeries,
    type Metrics,
    type SeriesDTO
} from './domain/metrics'
import {OverviewPage} from './pages/OverviewPage'
import type {Overview} from './domain/overview'
import type {ActivityUpdate} from './domain/activity'
import {GroupsPage} from "./pages/GroupsPage.tsx";
import {BottomNav, type BottomNavData} from "./components/bottomNav.tsx";
// import type { SeriesDTO } from './domain/metrics'
let sideBarData: SideBarData = {
    iconSrc: "",
    title: "Vigil",
    items: [
        {id: "overview", title: "Обзор", iconSrc: dashIcon, path: "/overview"},
        {id: "agents", title: "Агенты", iconSrc: agentsIcon, path: "/agents", countLabel: 0},
        {id: "groups", title: "Группы", path: "/groups"},
        {id: "reports", title: "Отчеты", path: "/reports"},
        {id: "events", title: "События", path: "/events"},
        {id: "streams", title: "Удаленный просмотр", path: "streams"}

    ]

}
let bottomNavData: BottomNavData = {
    iconSrc: "",
    title: "Vigil",
    items: [
        {id: "overview", title: "Обзор", iconSrc: dashIcon, path: "/overview"},
        {id: "agents", title: "Агенты", iconSrc: agentsIcon, path: "/agents", countLabel: 0},
        {id: "groups", title: "Группы", path: "/groups"},
        {id: "reports", title: "Отчеты", path: "/reports"},

    ]

}

const AppLayout = ({onlineCount}: { onlineCount: number }) => {
    // let location = useLocation()
    const isPhone = useMediaQuery('(max-width:599px)')
    // useEffect(
    //     () => {
    //
    //     },
    //     [location]
    // )
    useEffect(
        () => {
            if (onlineCount < 0) return
            sideBarData.items[1].countLabel = onlineCount
        },
        [onlineCount]
    )
    return (
        <div style={{
            display: `flex`,
            flexDirection: `${isPhone ? 'column' : 'row'}`,
            height: `100%`,
            alignItems: `stretch`
        }}>
            {!isPhone &&
                <SideBar
                    data={sideBarData}
                />
            }
            <div style={{flexGrow: 1, overflow: `auto`}}><Outlet/></div>
            {isPhone &&
                <BottomNav
                    data={bottomNavData}
                />
            }
        </div>
    )
}


function App() {
    const [sendMessage, setSendMessage] = useState<ClientMessage | null>(null)
    const [wsSocket, setWSSocket] = useState<WebSocket | null>(null)
    const [socketConnected, setConnected] = useState<boolean>(false)
    const [metrics, setMetrics] = useState<Metrics | undefined>()
    const [overview, setOverview] = useState<Overview | undefined>()
    const [activity, setActivity] = useState<ActivityUpdate | undefined>()


    const sendMessageRef = useRef(sendMessage);

    const sendMsg = () => {
        const currentMsg = sendMessageRef.current;
        if (currentMsg != null) {
            console.log("send message ", currentMsg)
            try {
                wsSocket?.send(JSON.stringify(currentMsg))
            } catch (ex: any) {
                console.log(`exception: ${ex}`)
            }
        }
    }
    useEffect(() => {
        const wsProto = location.protocol === 'https:' ? 'wss' : 'ws';
        const socket = new WebSocket(`${wsProto}://${location.host}/api/v1/ws`);
        socket.addEventListener("open", () => {
            console.log("start")
            setConnected(true)
        });

        socket.onclose = () => {
            console.log("connection closed")
            setConnected(false)
        };

        socket.addEventListener("message", (event) => {
            // console.log("Message from server ", event.data);
            const msg = JSON.parse(event.data) as Message;
            if (msg.type == MessageTypeSeries) {
                console.log("series message received", msg)
                const series = msg.payload as SeriesDTO
                if (series != undefined) {
                    setMetrics((prev) => FillMetricsFromSeries(series, prev ?? {}, msg.agentID))
                }
            } else if (msg.type == MessageTypeFleetOverview) {
                console.log("fleet overview message received", msg)
                setOverview(msg.payload as Overview)
            } else if (msg.type == MessageTypeActivity) {
                console.log("activity message received")
                setActivity(msg.payload as ActivityUpdate)
            }
        });

        setWSSocket(socket)

        return () => {
            socket.close()
        }
    }, []);

    useEffect(() => {
        if (socketConnected === true) {
            sendMsg()
        }
    }, [socketConnected])

    useEffect(() => {
        if (sendMessage == null)
            return
        sendMessageRef.current = sendMessage
        if (wsSocket == null || wsSocket.readyState != wsSocket.OPEN) {
            console.log("ws socket not ready: ", wsSocket?.readyState)
            return
        }
        console.log("ws message sent ", sendMessage)
        sendMsg()
    }, [sendMessage]);

    const handleLocationChange = (path: string) => {
        const splitted = path.split("/")
        if (splitted[1] == "agents" && splitted.length == 3) {
            console.log("agent id = ", splitted[2])
            const p: AgentDetailedMessage = {
                agents: [splitted[2]]
            }
            const msg: ClientMessage = {
                type: ClientMessageTypeAgentDetailed,
                payload: p
            }
            setSendMessage(msg)
        } else if (splitted[1] == "overview") {
            const msg: ClientMessage = {
                type: MessageTypeFleetOverview,
                payload: undefined
            }
            setSendMessage(msg)
        }
    }

    return (
        <div className='app'>
            <BrowserRouter>
                <RouteChangeTracker
                    handler={handleLocationChange}
                />
                <Routes>
                    <Route element={<AppLayout onlineCount={overview?.summary?.onlineAgents ?? 0}/>}>
                        <Route path="/" element={<Navigate to="/agents" replace/>}/>
                        <Route path="/overview" element={<OverviewPage
                            overviewProp={overview}/>}/>
                        <Route path="/agents" element={<AgentsPage/>}/>
                        <Route path="/agents/:id" element={<AgentPage
                            metricsProp={metrics}
                            activity={activity}/>}/>
                        <Route path="/agents/new" element={<NewAgentPage/>}/>
                        <Route path="/groups" element={<NotImplemented/>}/>
                        <Route path="/reports" element={<NotImplemented/>}/>
                        <Route path="/events" element={<NotImplemented/>}/>
                        <Route path="/streams" element={<NotImplemented/>}/>
                    </Route>
                    {/* <Route path="*" element={<NotFoundPage />} /> */}
                </Routes>
            </BrowserRouter>
        </div>
    )
}


function RouteChangeTracker({handler}: { handler: (path: string) => void }) {
    const location = useLocation();

    useEffect(() => {
        console.log('Route changed to:', location.pathname);
        handler(location.pathname)

    }, [location]);

    return null; // This component doesn't need to render anything visual
}

export default App

const NotImplemented = () => {
    return (
        <div>
            <h3>
                Not implemented
            </h3>
        </div>
    )
}