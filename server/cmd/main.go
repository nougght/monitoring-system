package main

import (
	"fmt"
	"log"
	"net"

	grpc_handler "github.com/nougght/monitoring-system/server/internal/handler/grpc"
	"google.golang.org/grpc"
)

func main() {

	s := grpc.NewServer()
	agentService := grpc_handler.NewAgentService()
	agentService.Register(s)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", 8090))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Println("server started")
	if err := s.Serve(l); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
