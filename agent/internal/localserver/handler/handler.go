package handler

import (
	"agent/internal/model"
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

// GetSpecifications godoc
// @Summary System specifications
// @Produce json
// @Success 200 {object} model.Specs
// @Failure 500 {object} map[string]string
// @Router /specs [get]
func (h *Handler) GetSpecifications(c *gin.Context) {
	var specs *model.Specs
	specs, err := h.metricsService.GetSpecs(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, specs)
}

// GetMetrics godoc
// @Summary Current system metrics
// @Produce json
// @Success 200 {object} model.Metrics
// @Failure 500 {object} map[string]string
// @Router /metrics [get]
func (h *Handler) GetMetrics(c *gin.Context) {
	var metrics *model.Metrics
	metrics = h.metricsService.GetMetrics()
	c.JSON(http.StatusOK, metrics)
}
