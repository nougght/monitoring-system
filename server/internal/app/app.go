package app

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nougght/monitoring-system/server/internal/config"
	grpc_handler "github.com/nougght/monitoring-system/server/internal/transport/grpc"
	"github.com/nougght/monitoring-system/server/internal/service"
	"github.com/nougght/monitoring-system/server/internal/storage/timescale"
	"github.com/nougght/monitoring-system/server/internal/storage/timescale/repository"
	"google.golang.org/grpc"
)

type App struct {
	DB *pgxpool.Pool

	Repositories *repository.Repositories

	Services *service.Services
}

func New(ctx context.Context, cfg *config.Config) *App {
	db, err := timescale.ConnectToDB(ctx, cfg.Postgres)
	if err != nil {
		log.Panicf("failed to connect to database: %v", err)
	}

	return &App{
		DB:           db,
		Repositories: repository.New(db),
		Services: service.New(service.ServicesOptions{
			Repositories: repository.New(db),
		}),
	}
}

func (a *App) Run(ctx context.Context) error {

	grpcServer := grpc.NewServer()
	agentService := grpc_handler.NewAgentService()
	agentService.Register(grpcServer)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", 8092))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcErrChan := make(chan error)
	go func() {
		log.Println("starting gRPC server")
		if err := grpcServer.Serve(l); err != nil {
			grpcErrChan <- err
		}
	}()

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: false,
	}))

	httpServer := http.Server{
		Addr:    ":8091",
		Handler: r.Handler(),
	}
	httpErrChan := make(chan error)
	go func() {
		log.Println("starting HTTP server")
		if err := httpServer.ListenAndServe(); err != nil {
			httpErrChan <- err
		}
	}()

	select {
	case err := <-grpcErrChan:
		log.Fatalf("failed to serve gRPC: %v", err)
	case err := <-httpErrChan:
		log.Fatalf("failed to serve HTTP: %v", err)
	}

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	grpcServer.GracefulStop()
	err = httpServer.Shutdown(shutdownCtx)
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("failed to shutdown HTTP server: %v", err)
	}

	return nil
}
