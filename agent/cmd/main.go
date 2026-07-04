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

	_ "agent/api"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Monitoring Agent API
// @version         1.0
// @host            localhost:8088
// @BasePath        /
func main() {
	cfg := config.LoadConfig("config.yaml")

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

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// r.GET("/swagger/*any", func(ctx *gin.Context) {
	// 	ctx.Request.URL.Path = "/doc.json"
	// 	r.HandleContext(ctx)
	// })

	r.GET("/ws", wsHandler.HandleConnection)
	r.GET("/specs", h.GetSpecifications)
	r.GET("/metrics", h.GetMetrics)

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
		err := server.Shutdown(ctx)
		service.StopServices()
		shutdownChan <- err
	}()

	log.Println("http server started on :8088")
	err = server.ListenAndServe()

	if err != nil && err != http.ErrServerClosed {
		log.Fatal("ListenAndServe: ", err)
	}

	err = <-shutdownChan
	if err != nil {
		log.Fatal("shutdown error: ", err)
	}
	log.Println("http server stopped")
}
