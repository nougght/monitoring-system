package rest

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nougght/monitoring-system/server/internal/model"
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
	group.POST("/:agentID/setupconfig", h.DownloadAgentConfig)
	group.GET("/:agentID/specifications", h.GetAgentSpecs)
}

// CreateAgent godoc
// @Id createAgent
// @Summary Create new agent
// @Accept json
// @Produce json
// @Param request body dto.CreateAgentBody true "Create agent body"
// @Success 200 {object} dto.CreateAgentResponse
// @Failure      400  {object}  map[string]any
// @Failure      404  {object}  map[string]any
// @Failure      500  {object}  map[string]any
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

// DownloadAgentConfig godoc
// @Id downloadAgentConfig
// @Summary Download agent setup config
// @Accept json
// @Produce application/x-yaml
// @Param agentID path string true "Agent ID"
// @Param request body dto.AgentConfigBody true "Download agent config body"
// @Success 200 {object} dto.AgentConfigResponse
// @Failure      400  {object}  map[string]any
// @Failure      404  {object}  map[string]any
// @Failure      500  {object}  map[string]any
// @Router /agents/{agentID}/setupconfig [post]
func (h *AgentHandler) DownloadAgentConfig(c *gin.Context) {
	agentID, err := uuid.Parse(c.Param("agentID"))
	if err != nil {
		handleError(c, fmt.Errorf("invalid 'agentID' path parameter: %w", model.ErrBadRequest))
		return
	}

	var body dto.AgentConfigBody
	err = c.ShouldBindJSON(&body)
	if err != nil {
		handleError(c, fmt.Errorf("invalid request body: %w", model.ErrBadRequest))
		return
	}

	config := h.agentRegistryService.GenerateAgentSetupConfig(c.Request.Context(), agentID, body.EnrollmentKey)
	c.Writer.Header().Set("Content-Disposition", `attachment; filename="config.yaml"`)
	c.YAML(http.StatusOK, config)
}

// GetAllAgents godoc
// @Id getAllAgents
// @Summary Get all agents
// @Produce json
// @Success 200 {array} dto.AgentDTO
// @Failure      400  {object}  map[string]any
// @Failure      404  {object}  map[string]any
// @Failure      500  {object}  map[string]any
// @Router /agents [get]
func (h *AgentHandler) GetAllAgents(c *gin.Context) {
	res, err := h.agentRegistryService.GetAllAgents(c.Request.Context())
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.Map(res, mapper.AgentToDTO))
}

// GetAllAgents godoc
// @Id getAgentSpecs
// @Summary Get agent specifications
// @Produce json
// @Param agentID path string true "Agent ID"
// @Success 200 {object} dto.SpecsDTO
// @Failure      400  {object}  map[string]any
// @Failure      404  {object}  map[string]any
// @Failure      500  {object}  map[string]any
// @Router /agents/{agentID}/specifications [get]
func (h *AgentHandler) GetAgentSpecs(c *gin.Context) {
	agentID, err := uuid.Parse(c.Param("agentID"))
	if err != nil {
		handleError(c, fmt.Errorf("invalid 'agentID' path parameter: %w", model.ErrBadRequest))
		return
	}
	res, err := h.agentRegistryService.GetSpecifications(c.Request.Context(), agentID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, mapper.SpecsToDTO(res))
}
