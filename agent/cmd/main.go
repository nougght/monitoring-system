package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/nougght/monitoring-system/shared/go/util"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"time"

	"agent/internal/config"
	grpc_client "agent/internal/grpc"
	"agent/internal/model"

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
	var configPath string
	var rootCmd = &cobra.Command{
		Use: "agent",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if configPath == "" {
				return fmt.Errorf("config path required")
			}
			// TODO: auto resolve config path
			return nil
		},
	}
	rootCmd.PersistentFlags().StringVar(&configPath, "config", "", "path to start config")

	var runCmd = &cobra.Command{
		Use: "run",
		RunE: func(cmd *cobra.Command, args []string) error {
			var cfg model.StartConfig
			util.MustReadYaml(configPath, cfg)
			runAgent(cfg)
			return nil
		},
	}
	rootCmd.AddCommand(runCmd)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
	// var enrollCmd = &cobra.Command{
	// 	Use: "run",
	// 	RunE: func(cmd *cobra.Command, args []string) error {
	// 		// cfg := loadConfig(configPath)
	// 		// enroll(cfg)
	// 		return nil
	// 	},
	// }

	// var statusCmd = &cobra.Command{
	// 	Use: "status",
	// 	RunE: func(cmd *cobra.Command, args []string) error {
	// 		// return cmdStatus(configPath)
	// 	},
	// }
}

func resolveConfigPaht() {

}
func runAgent(startConfig model.StartConfig) {
	cfg := config.MustLoadConfig("config.yaml")
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

	server := &http.Server{
		Addr:    ":8111",
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

	grpcClient, err := grpc.NewClient("127.0.0.1:8092", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("failed to create grpc client: ", err)
	}
	grpcAgentClient := grpc_client.NewAgentClient(grpcClient, cfg, service.GetMetricsService())
	err = grpcAgentClient.Connect(rootCtx)
	if err != nil {
		log.Println("failed to connect to grpc server: ", err)
	}

	log.Println("http server started on :8111")
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
