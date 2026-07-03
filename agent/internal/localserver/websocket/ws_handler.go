package ws

import (
	"agent/internal/service/metrics"
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type WsHandler struct {
	upgrader       *websocket.Upgrader
	metricsService *metrics.MetricsService
	rootCtx        context.Context
}

func NewWsHandler(metricsService *metrics.MetricsService, rootCtx context.Context) *WsHandler {
	return &WsHandler{
		upgrader: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		},
		metricsService: metricsService,
		rootCtx:        rootCtx,
	}
}
func (h *WsHandler) HandleConnection(c *gin.Context) {
	log.Printf("new connection: %s", c.Request.Host)
	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("failed to upgrade to websocket: ", err)
		return
	}
	go h.HandleMessages(conn)
}

func (h *WsHandler) HandleMessages(conn *websocket.Conn) {
	defer conn.Close()

	ticker := time.NewTicker(time.Millisecond * 100)
	defer ticker.Stop()
	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("failed to read message: ", err)
				return
			}
			log.Printf("received message: %s", message)
		}
	}()

	for {
		select {
		case <-ticker.C:
			metrics := h.metricsService.GetMetrics()
			err := conn.WriteJSON(metrics)
			if err != nil {
				log.Println("failed to write message: ", err)
				return
			}
		case <-h.rootCtx.Done():
			log.Println("websocket ctx done")
			return
		}
	}
}
