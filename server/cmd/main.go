package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/nougght/monitoring-system/server/internal/app"
	"github.com/nougght/monitoring-system/server/internal/config"
)

func main() {
	cfg := config.MustLoadConfig()
	rootCtx, cancel := signal.NotifyContext(
		context.Background(), syscall.SIGINT, syscall.SIGTERM,
	)
	defer cancel()

	app := app.New(rootCtx, cfg)

	if err := app.Run(rootCtx); err != nil {
		log.Fatalf("failed to run app: %v", err)
	}
}
