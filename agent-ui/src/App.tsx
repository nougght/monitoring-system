import {
    useEffect,
    useState
} from 'react'

import TabBar from './components/TabBar';
import Monitoring from './components/Monitoring';
import Specifications from './components/Specifications';
import './App.css';
import type { Metrics } from './domain/metrics';

interface Tab {
    text: string;
    content: React.ReactNode;
}

function App() {
    const [metrics, setMetrics] = useState<Metrics | null>(null);
    const [activeTab, setActiveTab] = useState(0)

    const tabs: Tab[] = [
        {
            text: "Monitoring",
            content: <div><Monitoring
                focusedWindow={metrics?.focusedWindow}
                cpuPercent={metrics?.cpuPercent} /></div>
        },
        {
            text: "Specifications",
            content: <div><Specifications /></div>
        }
    ]
    useEffect(() => {
        const socket = new WebSocket("ws://127.0.0.1:8088/ws");
        console.log("start")
        socket.addEventListener("open", () => {
            socket.send("Hello Server!");
        });

        socket.onclose = () => {
            console.log("connection closed")
        };

        // Listen for messages
        socket.addEventListener("message", (event) => {
            // console.log("Message from server ", event.data);
            const metrics = JSON.parse(event.data) as Metrics;
            // console.log("metrics", metrics);
            setMetrics(metrics);
        });
    }, []);

    return (
        <>
            <header>
                <TabBar tabs={tabs.map((tab) => tab.text)} onSwitch={setActiveTab} activeTab={activeTab} />
            </header>
            <main>
                {tabs[activeTab].content}
            </main>
        </>
    )
}

export default App
