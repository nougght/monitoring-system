package handler

import (
	"agent/internal/localserver/dto"
	"agent/internal/localserver/mapper"
	"agent/internal/model"
	"agent/internal/service/agentcore"
	"agent/internal/service/metrics"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	core           *agentcore.CoreService
	metricsService *metrics.MetricsService
}

func NewHandler(core *agentcore.CoreService, metricsService *metrics.MetricsService) *Handler {
	return &Handler{
		metricsService: metricsService,
		core:           core,
	}
}

func (h *Handler) LocalAgentStatus(c *gin.Context) {
	connectionStatus := "connected"
	if h.core.State().Connected() {
		connectionStatus = "not connected"
		lastConnectedAt := h.core.State().LastConnectedAt()
		if lastConnectedAt != nil {
			connectionStatus += fmt.Sprintf(", last - %t", lastConnectedAt)
		}
	}
	status := dto.Status{
		AgentID:          h.core.State().AgentID(),
		ConnectionStatus: connectionStatus,
		CertInfo:         "not implemented",
		ServerInfo:       "not implemented",
	}

	c.JSON(http.StatusOK, status)
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
// @Success 200 {object} dto.Metrics
// @Failure 500 {object} map[string]string
// @Router /metrics [get]
func (h *Handler) GetMetrics(c *gin.Context) {
	metrics := mapper.ConvertMetricsToDto(h.metricsService.GetMetrics())
	log.Println("metrics", metrics)
	c.JSON(http.StatusOK, metrics)
}
