package app

import (
	"agent/internal/config"
	grpc_client "agent/internal/grpc/agent_client"
	"agent/internal/localserver/handler"
	ws "agent/internal/localserver/websocket"
	"agent/internal/service"
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func RunAgent(setupConfig *config.SetupConfig) error {
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		return err
	}
	rootCtx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	service, err := service.GetServices(setupConfig, cfg)
	if err != nil {
		return fmt.Errorf("failed to setup services: %w", err)
	}
	service.StartServices(rootCtx)

	h := handler.NewHandler(service.GetCoreService(), service.GetMetricsService())
	wsHandler := ws.NewWsHandler(service.GetMetricsService(), rootCtx)

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
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

	r.GET("/status", h.LocalAgentStatus)

	server := &http.Server{
		Addr:    ":8111",
		Handler: r.Handler(),
	}

	grpcTLS := service.GetCoreService().TLSConfigForGRPC()
	grpcClient, err := grpc.NewClient(setupConfig.ServerAddress,
		grpc.WithTransportCredentials(credentials.NewTLS(grpcTLS)),
	)
	if err != nil {
		return fmt.Errorf("failed to create grpc client: %w", err)
	}

	grpcAgentClient := grpc_client.NewAgentClient(grpcClient,
		cfg,
		service.GetCoreService().State(),
		service.GetMetricsService(),
	)

	grpcCtx, cancel := context.WithCancel(rootCtx)
	defer cancel()

	go func() {
		err = grpcAgentClient.Connect(grpcCtx)
		// TODO: add retrying
		if err != nil {
			log.Println("failed to connect to grpc server: ", err)
		}
		cancel()
	}()
	go func() {
		err = grpcAgentClient.StartStreamMJPEG(grpcCtx)
		if err != nil {
			log.Println("failed to connect to grpc streaming: ", err)
		}
	}()

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

	log.Println("http server started on :8111")
	err = server.ListenAndServe()

	if err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("http server error: %w", err)
	}

	err = <-shutdownChan
	if err != nil {
		return fmt.Errorf("http server shutdown error: %w", err)
	}
	log.Println("http server stopped")
	return nil
}
