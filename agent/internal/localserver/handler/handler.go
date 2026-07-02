package handler

import (
	"agent/internal/service/metrics"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	metricsService *metrics.MetricsService
}

func NewHandler(metricsService *metrics.MetricsService) *Handler {
	return &Handler{
		metricsService: metricsService,
	}
}

func (h *Handler) GetSpecifications(c *gin.Context) {
	specs, err := h.metricsService.GetSpecs(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, specs)
}
