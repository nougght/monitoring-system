package app

import (
	"context"
	"crypto"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nougght/monitoring-system/server/internal/config"
	"github.com/nougght/monitoring-system/server/internal/model"
	"github.com/nougght/monitoring-system/server/internal/service"
	"github.com/nougght/monitoring-system/server/internal/storage/timescale"
	"github.com/nougght/monitoring-system/server/internal/storage/timescale/repository"
	agent_grpc "github.com/nougght/monitoring-system/server/internal/transport/grpc/agent_server"
	enrollment_grpc "github.com/nougght/monitoring-system/server/internal/transport/grpc/enrollment_server"
	"github.com/nougght/monitoring-system/server/internal/transport/rest"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/nougght/monitoring-system/shared/go/cert_store"
)

type App struct {
	Config *config.Config
	DB     *pgxpool.Pool

	Repositories *repository.Repositories

	Services *service.Services

	HTTPServer           *http.Server
	AgentGRPCServer      *grpc.Server
	EnrollmentGRPCServer *grpc.Server
	CertStore            *model.CertStore

	ca     *x509.Certificate
	key    crypto.Signer
	rootCA *x509.CertPool
}

func New(ctx context.Context, cfg *config.Config) *App {
	certStore := cert_store.NewCertStore(cfg.Cert.IntCAPath, cfg.Cert.IntKeyPath, cfg.Cert.CAPath)

	cert, err := certStore.LoadCertificate()
	if err != nil {
		log.Panicf("failed load certs: %s", err.Error())
	}
	intCA := cert.Leaf
	intKey, err := certStore.LoadKey()
	if err != nil {
		log.Panicf("failed load certs: %s", err.Error())
	}
	rootCA, err := certStore.LoadCA()
	if err != nil {
		log.Panicf("failed load certs: %s", err.Error())
	}

	err = verifyIntermediateCA(intCA, rootCA)
	if err != nil {
		log.Panic("failed to validate intermediate CA")
	}

	db, err := timescale.ConnectToDB(ctx, cfg.Postgres)
	if err != nil {
		log.Println("failed to connect to database, retry after 500ms")
		select {
		case <-time.After(time.Microsecond * 500):
			db, err = timescale.ConnectToDB(ctx, cfg.Postgres)
			if err != nil {
				log.Panicf("failed to connect to database: %s", err.Error())
			}
		case <-ctx.Done():
			log.Panicf("%s", ctx.Err().Error())
		}
	}

	services := service.New(service.ServicesOptions{
		Config:       cfg,
		Repositories: repository.New(db),
		Transactor:   db,
		Cert: &model.Certs{
			CA:     intCA,
			Key:    intKey,
			RootCA: rootCA,
		},
	})

	httpServer := rest.NewServer(cfg, *services)

	agentServer := grpc.NewServer(
		grpc.StreamInterceptor(agent_grpc.BidirectionalServerInterceptor),
		grpc.Creds(
			credentials.NewTLS(&tls.Config{Certificates: []tls.Certificate{*cert},
				ClientAuth: tls.RequireAndVerifyClientCert, //  mTLS
				ClientCAs:  rootCA,
				MinVersion: tls.VersionTLS12,
			}),
		))

	agentService := agent_grpc.NewAgentService(services.AgentInteractionService())
	agentService.Register(agentServer)

	services.AgentInteractionService().SetRequester(agentService)

	enrollmentServer := grpc.NewServer(grpc.Creds(
		credentials.NewTLS(&tls.Config{
			Certificates: []tls.Certificate{*cert},
			ClientAuth:   tls.NoClientCert, // явно НЕ требуем клиентский сертификат
			MinVersion:   tls.VersionTLS12,
		})),
	)

	enrollmentService := enrollment_grpc.NewEnrollmentService(services.AgentInteractionService())
	enrollmentService.Register(enrollmentServer)

	return &App{
		Config:               cfg,
		DB:                   db,
		Repositories:         repository.New(db),
		Services:             services,
		HTTPServer:           httpServer,
		AgentGRPCServer:      agentServer,
		EnrollmentGRPCServer: enrollmentServer,
		ca:                   intCA,
		key:                  intKey,
		rootCA:               rootCA,
	}
}

func (a *App) Run(ctx context.Context) error {
	httpErrChan := make(chan error)
	go func() {
		log.Println("starting HTTP server")
		if err := a.HTTPServer.ListenAndServe(); err != nil {
			httpErrChan <- err
		}
	}()

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.Config.GRPC.MainPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcAgentErrChan := make(chan error)
	go func() {
		log.Printf("starting gRPC server :%d", a.Config.GRPC.MainPort)
		if err := a.AgentGRPCServer.Serve(l); err != nil {
			grpcAgentErrChan <- err
		}
	}()

	l2, err := net.Listen("tcp", fmt.Sprintf(":%d", a.Config.GRPC.EnrollmentPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpEnrollmentErrChan := make(chan error)
	go func() {
		log.Printf("starting gRPC enrollment server :%d", a.Config.GRPC.EnrollmentPort)
		if err := a.EnrollmentGRPCServer.Serve(l2); err != nil {
			grpEnrollmentErrChan <- err
		}
	}()

	select {
	case err := <-grpcAgentErrChan:
		log.Fatalf("failed to serve gRPC agent server: %v", err)
	case err := <-grpEnrollmentErrChan:
		log.Fatalf("failed to serve gRPC enrollment server: %v", err)
	case err := <-httpErrChan:
		log.Fatalf("failed to serve HTTP: %v", err)
	}

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	a.AgentGRPCServer.GracefulStop()
	a.EnrollmentGRPCServer.GracefulStop()
	err = a.HTTPServer.Shutdown(shutdownCtx)
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("failed to shutdown HTTP server: %v", err)
	}

	return nil
}

func verifyIntermediateCA(intermediateCA *x509.Certificate, rootPool *x509.CertPool) error {
	_, err := intermediateCA.Verify(x509.VerifyOptions{Roots: rootPool})
	return err
}
