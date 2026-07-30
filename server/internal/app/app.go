package app

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nougght/monitoring-system/server/internal/config"
	"github.com/nougght/monitoring-system/server/internal/service"
	"github.com/nougght/monitoring-system/server/internal/storage/timescale"
	"github.com/nougght/monitoring-system/server/internal/storage/timescale/repository"
	grpc_handler "github.com/nougght/monitoring-system/server/internal/transport/grpc"
	"github.com/nougght/monitoring-system/server/internal/transport/rest"
	"google.golang.org/grpc"
)

type App struct {
	DB *pgxpool.Pool

	Repositories *repository.Repositories

	Services *service.Services

	HTTPServer *http.Server
	GRPCServer *grpc.Server
}

func New(ctx context.Context, cfg *config.Config) *App {
	db, err := timescale.ConnectToDB(ctx, cfg.Postgres)
	if err != nil {
		log.Panicf("failed to connect to database: %v", err)
	}

	services := service.New(service.ServicesOptions{
		Config:       cfg,
		Repositories: repository.New(db),
		Transactor:   db,
	})

	httpServer := rest.NewServer(cfg, *services)

	grpcServer := grpc.NewServer()
	agentService := grpc_handler.NewAgentService()
	agentService.Register(grpcServer)

	return &App{
		DB:           db,
		Repositories: repository.New(db),
		Services:     services,
		HTTPServer:   httpServer,
		GRPCServer:   grpcServer,
	}
}

func (a *App) Run(ctx context.Context) error {
	httpErrChan := make(chan error)
	go func() {
		log.Println("starting HTTP server")
		if err := a.HTTPServer.ListenAndServe(); err != nil {
			httpErrChan <- err
		}
	}()

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", 8092))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcErrChan := make(chan error)
	go func() {
		log.Println("starting gRPC server")
		if err := a.GRPCServer.Serve(l); err != nil {
			grpcErrChan <- err
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

	a.GRPCServer.GracefulStop()
	err = a.HTTPServer.Shutdown(shutdownCtx)
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("failed to shutdown HTTP server: %v", err)
	}

	return nil
}
