package rest

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/nougght/monitoring-system/server/internal/config"
	"github.com/nougght/monitoring-system/server/internal/service"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewServer(cfg *config.Config, services service.Services) *http.Server {
	handlers := newHandlers(&services)

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: false,
	}))
	api := r.Group("/api/v1")

	api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	handlers.AgentHandler().RegisterRoutes(api)
	handlers.StreamHandler().RegisterRoutes(api)

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.HTTP.ServerPort),
		Handler: r.Handler(),
	}
	log.Println(server.Addr)
	return &server
}
