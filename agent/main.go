package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"

	"sort"

	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/process"

	// "agent/ui"
	"agent/utils"
	// "syscall"
	// "os/exec"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	log.Println("new connection")
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	go wsHandler(ws)
}

func wsHandler(conn *websocket.Conn) {
	defer conn.Close()
	for {
		ticker := time.NewTicker(time.Second)
		<-ticker.C
		activeWindow := utils.GetActiveWindowTitle()
		if err := conn.WriteMessage(websocket.TextMessage, []byte(activeWindow)); err != nil {
			log.Println(err)
			break
		}
	}
}

func main() {
	v, _ := mem.VirtualMemory()

	// almost every return value is a struct
	fmt.Printf("Total: %v, Free:%v, UsedPercent:%f%%\n", v.Total, v.Free, v.UsedPercent)

	c, _ := cpu.Percent(time.Second*2, true)
	fmt.Println(c)

	fmt.Println(v)

	p, _ := process.Processes()
	sort.Slice(p, func(i, j int) bool {
		perc1, _ := p[i].MemoryPercent()
		perc2, _ := p[j].MemoryPercent()
		return perc1 > perc2
	})
	for _, proc := range p {
		name, _ := proc.Name()
		perc, _ := proc.MemoryPercent()
		fmt.Println(name, perc)
	}

	fmt.Println("window")
	fmt.Println(utils.GetActiveWindowTitle())

	// convert to JSON. String() is also implemented
	// ui.Start()

	http.HandleFunc("/ws", handleConnections)
	log.Println("http server started on :8000")
	err := http.ListenAndServe(":8000", nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
