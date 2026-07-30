package rest

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	agent "github.com/nougght/monitoring-system/server/internal/service/agent_interaction"
	agentregistry "github.com/nougght/monitoring-system/server/internal/service/agent_registry"
	"github.com/nougght/monitoring-system/server/internal/transport/dto/mapper"
	dto "github.com/nougght/monitoring-system/server/internal/transport/dto/types"
	"github.com/nougght/monitoring-system/shared/go/util"
)

type AgentHandler struct {
	agentRegistryService *agentregistry.AgentRegistryService
	agentService         *agent.AgentInteractionService
}

func newAgentHandler(agentRegistryServcie *agentregistry.AgentRegistryService,
	agentService *agent.AgentInteractionService) *AgentHandler {
	if agentRegistryServcie == nil || agentService == nil {
		log.Panicf("agent handler params required")
	}
	return &AgentHandler{
		agentRegistryService: agentRegistryServcie,
		agentService:         agentService,
	}
}

func (h *AgentHandler) RegisterRoutes(r *gin.RouterGroup) {
	group := r.Group("/agents")

	group.POST("", h.CreateAgent)
	group.GET("", h.GetAllAgents)
}

// CreateAgent godoc
// @Summary Create new agent
// @Accept json
// @Produce json
// @Param request body dto.CreateAgentBody true "Create agent body"
// @Success 200 {object} dto.CreateAgentResponse
// @Failure      400  {object}  gin.H
// @Failure      404  {object}  gin.H
// @Failure      500  {object}  gin.H
// @Router /agents [post]
func (h *AgentHandler) CreateAgent(c *gin.Context) {
	var body *dto.CreateAgentBody
	err := c.ShouldBindJSON(&body)
	if err != nil {
		handleError(c, err)
		return
	}

	res, err := h.agentRegistryService.CreateAgent(c.Request.Context(), body.Name, body.Description)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, mapper.CreateAgentResultToDTO(res))
}

// GetAllAgents godoc
// @Summary Get all agents
// @Produce json
// @Success 200 {array} dto.AgentDTO
// @Failure      400  {object}  gin.H
// @Failure      404  {object}  gin.H
// @Failure      500  {object}  gin.H
// @Router /agents [get]
func (h *AgentHandler) GetAllAgents(c *gin.Context) {
	res, err := h.agentRegistryService.GetAllAgents(c.Request.Context())
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.Map(res, mapper.AgentToDTO))
}
