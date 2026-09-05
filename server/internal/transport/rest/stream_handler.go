package rest

import (
	"bytes"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nougght/monitoring-system/server/internal/model"
	agent "github.com/nougght/monitoring-system/server/internal/service/agent_interaction"
	agentregistry "github.com/nougght/monitoring-system/server/internal/service/agent_registry"
)

type StreamHandler struct {
	agentRegistryService *agentregistry.AgentRegistryService
	agentService         *agent.AgentInteractionService
}

func newStreamHandler(agentRegistryServcie *agentregistry.AgentRegistryService,
	agentService *agent.AgentInteractionService) *StreamHandler {
	if agentRegistryServcie == nil || agentService == nil {
		log.Panicf("agent handler params required")
	}
	return &StreamHandler{
		agentRegistryService: agentRegistryServcie,
		agentService:         agentService,
	}
}

func (s *StreamHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/agents/:agentID/frames", s.StreamFrames)
}

// StreamFrames godoc
// @Id getStreamFrames
// @Summary Get stream frames
// @Produce multipart/x-mixed-replace
// @Param agentID path string true "Agent ID"
// @Success      200  {file}    binary  "MJPEG stream"
// @Failure      400  {object}  map[string]any
// @Failure      404  {object}  map[string]any
// @Failure      500  {object}  map[string]any
// @Router /agents/{agentID}/frames [get]
func (s *StreamHandler) StreamFrames(c *gin.Context) {
	agentID, err := uuid.Parse(c.Param("agentID"))
	if err != nil {
		handleError(c, fmt.Errorf("invalid 'agentID' path parameter: %w", model.ErrBadRequest))
		return
	}

	viewerID := uuid.New()
	framesChan, err := s.agentService.SubStreaming(agentID, viewerID)
	if err != nil {
		handleError(c, err)
		return
	}
	boundary := uuid.NewString()
	c.Header("Content-Type", fmt.Sprintf("multipart/x-mixed-replace; boundary=%s", boundary))
	c.Header("Cache-Control", "no-store")
	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Flush()

	var buff bytes.Buffer
	for {
		select {
		case <-c.Request.Context().Done():
			s.agentService.UnsubAllStreaming(viewerID)
		case frame := <-framesChan:
			log.Println("frame")

			_, err := fmt.Fprintf(c.Writer,
				"--%s\r\nContent-Type: image/jpeg\r\nContent-Length: %d\r\n\r\n",
				boundary, len(frame))

			_, err = c.Writer.Write(frame)
			if err != nil {
				log.Println("failed to send frame by http: %w", err)
			}
			buff.Reset()
		}
	}
}
