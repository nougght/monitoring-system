package agent_server

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/nougght/monitoring-system/server/internal/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

// wrappedServerStream оборачивает стандартный поток для перехвата сообщений
type wrappedServerStream struct {
	grpc.ServerStream
}

func (w *wrappedServerStream) RecvMsg(m any) error {
	err := w.ServerStream.RecvMsg(m)
	if err == nil {
		log.Printf("[Server Stream] Получено сообщение типа: %T", m)
		_, err := authInterceptor(w.ServerStream.Context())
		if err != nil {
			log.Printf("[Server Stream] Ошибка аутентификации: %v", err)
			return status.Error(codes.Unauthenticated, "authentication failed")
		}
	}
	return err
}

func (w *wrappedServerStream) SendMsg(m any) error {
	log.Printf("[Server Stream] Отправка сообщения типа: %T", m)
	return w.ServerStream.SendMsg(m)
}

func BidirectionalServerInterceptor(
	srv any,
	ss grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	start := time.Now()
	log.Printf("[Server] Открыт bidirectional поток для метода: %s", info.FullMethod)

	// Оборачиваем оригинальный стрим в нашу структуру
	wrapper := &wrappedServerStream{ServerStream: ss}

	// Передаем обертку в обработчик бизнес-логики
	err := handler(srv, wrapper)

	log.Printf("[Server] Поток %s закрыт. Длительность: %v. Ошибка: %v", info.FullMethod, time.Since(start), err)
	return err
}

func authInterceptor(ctx context.Context) (context.Context, error) {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "no peer info")
	}
	tlsInfo, ok := p.AuthInfo.(credentials.TLSInfo)
	if !ok || len(tlsInfo.State.PeerCertificates) == 0 {
		return nil, status.Error(codes.Unauthenticated, "no client certificate")
	}
	agentID, err := uuid.Parse(tlsInfo.State.PeerCertificates[0].Subject.CommonName)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid client CN")
	}
	return context.WithValue(ctx, model.ContextKeyAgentID, agentID), nil
}
