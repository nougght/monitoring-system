package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"time"

	"agent/internal/config"
	"agent/internal/localserver/handler"
	ws "agent/internal/localserver/websocket"
	"agent/internal/service"
)

func main() {
	cfg := config.LoadConfig("config.yaml")

	// fmt.Println(cfg)
	// v, _ := mem.VirtualMemory()

	// // almost every return value is a struct
	// fmt.Printf("Total: %v, Free:%v, UsedPercent:%f%%\n", v.Total, v.Free, v.UsedPercent)

	// c, _ := cpu.Percent(time.Second*2, true)
	// fmt.Println(c)

	// fmt.Println(v)

	// p, _ := process.Processes()
	// sort.Slice(p, func(i, j int) bool {
	// 	perc1, _ := p[i].MemoryPercent()
	// 	perc2, _ := p[j].MemoryPercent()
	// 	return perc1 > perc2
	// })
	// // for _, proc := range p {
	// // 	name, _ := proc.Name()
	// // 	perc, _ := proc.MemoryPercent()
	// // 	fmt.Println(name, perc)
	// // }
	// getSpecifications()

	// fmt.Println("window")
	// fmt.Println(utils.GetActiveWindowTitle())

	// convert to JSON. String() is also implemented
	// ui.Start()

	rootCtx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	service, err := service.GetServices(cfg)
	if err != nil {
		log.Fatal("failed to get services: ", err)
	}
	service.StartServices(rootCtx)

	h := handler.NewHandler(service.GetMetricsService())
	wsHandler := ws.NewWsHandler(service.GetMetricsService(), rootCtx)

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // или "*" для всех
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: false,
	}))
	r.GET("/ws", wsHandler.HandleConnection)
	r.GET("/specs", h.GetSpecifications)

	server := &http.Server{
		Addr:    ":8088",
		Handler: r.Handler(),
	}

	shutdownChan := make(chan error)
	go func() {
		defer close(shutdownChan)
		<-rootCtx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			shutdownChan <- err
		}
	}()

	log.Println("http server started on :8088")
	err = server.ListenAndServe()

	if err != nil && err != http.ErrServerClosed {
		log.Fatal("ListenAndServe: ", err)
	}

	err = <-shutdownChan
	if err != nil {
		log.Fatal("shutdown: ", err)
	}
	log.Println("http server stopped")
}
