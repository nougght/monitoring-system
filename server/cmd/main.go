package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	_ "github.com/nougght/monitoring-system/server/api"
	"github.com/nougght/monitoring-system/server/internal/app"
	"github.com/nougght/monitoring-system/server/internal/config"
)

// @title           Monitoring Server API
// @version         1.0
// @host            localhost:8091
// @BasePath        /api/v1
func main() {
	cfg := config.MustLoadConfig("config.yaml")
	rootCtx, cancel := signal.NotifyContext(
		context.Background(), syscall.SIGINT, syscall.SIGTERM,
	)
	defer cancel()

	app := app.New(rootCtx, cfg)

	if err := app.Run(rootCtx); err != nil {
		log.Fatalf("failed to run app: %v", err)
	}
}
