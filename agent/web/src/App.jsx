import { useState } from 'react'
import { useEffect } from 'react';


function App() {
    const [activeWindow, setActiveWindow] = useState("");

    useEffect(() => {
        const socket = new WebSocket("ws://localhost:8000/ws");
        console.log("start")
        socket.addEventListener("open", (event) => {
            socket.send("Hello Server!");
        });

        socket.onclose = (event) => {
            console.log("connection closed")
        };

        // Listen for messages
        socket.addEventListener("message", (event) => {
            console.log("Message from server ", event.data);
            setActiveWindow(event.data);
        });
    }, []);

    return (
        <>
            <section>
                <h2>Active window</h2>
                <p id="active-window">{activeWindow}</p>
            </section>
        </>
    )
}

export default App
