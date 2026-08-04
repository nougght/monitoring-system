package service

import (
	"agent/internal/config"
	"agent/internal/service/agentcore"
	"agent/internal/service/collector"
	"agent/internal/service/metrics"
	"context"
	"fmt"
)

type Service struct {
	coreService      *agentcore.CoreService
	collectorService *collector.CollectorService
	metricsService   *metrics.MetricsService
}

func GetServices(setupCfg *config.SetupConfig, cfg *config.Config) (*Service, error) {
	coreService, err := agentcore.NewCore(setupCfg, agentcore.NewCertStore(setupCfg.CaPath))
	if err != nil {
		return nil, fmt.Errorf("failed to init agent core service: %w", err)
	}
	collectorService := collector.NewCollectorService(cfg)
	metricsService := metrics.NewMetricsService(cfg, collectorService.GetSpecifications)
	collectorService.SetMetricsConsumer(metricsService)
	return &Service{
		coreService:      coreService,
		collectorService: collectorService,
		metricsService:   metricsService,
	}, nil
}

func (s *Service) StartServices(ctx context.Context) {
	s.collectorService.StartCollectors(ctx)
}

func (s *Service) StopServices() {
	s.collectorService.StopCollectors()
}

func (s *Service) GetCoreService() *agentcore.CoreService {
	return s.coreService
}

func (s *Service) GetMetricsService() *metrics.MetricsService {
	return s.metricsService
}

func (s *Service) GetCollectorService() *collector.CollectorService {
	return s.collectorService
}
