package rest

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nougght/monitoring-system/server/internal/model"
	"github.com/nougght/monitoring-system/server/internal/service"
)

func handleError(c *gin.Context, err error) {
	responseCode := http.StatusInternalServerError

	switch {
	case errors.Is(err, model.ErrBadRequest):
		responseCode = http.StatusBadRequest
	case errors.Is(err, model.ErrNotFound):
		responseCode = http.StatusNotFound
	case errors.Is(err, model.ErrServiceUnavailable):
		responseCode = http.StatusServiceUnavailable
	}

	log.Println(err)
	c.JSON(responseCode, gin.H{"error": err.Error()})
}

type Handlers struct {
	agentHandler  *AgentHandler
	streamHandler *StreamHandler
}

func newHandlers(services *service.Services) *Handlers {
	agent := newAgentHandler(
		services.AgentRegistry(),
		services.AgentInteractionService(),
	)

	stream := newStreamHandler(
		services.AgentRegistry(),
		services.AgentInteractionService(),
	)

	return &Handlers{
		agentHandler:  agent,
		streamHandler: stream,
	}
}

func (h *Handlers) AgentHandler() *AgentHandler {
	return h.agentHandler
}

func (h *Handlers) StreamHandler() *StreamHandler {
	return h.streamHandler
}
