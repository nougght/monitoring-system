package service

import (
	"agent/internal/config"
	"agent/internal/service/collector"
	"agent/internal/service/metrics"
	"context"
)

type Service struct {
	collectorService *collector.CollectorService
	metricsService   *metrics.MetricsService
}

func GetServices(cfg *config.Config) (*Service, error) {
	collectorService := collector.NewCollectorService(cfg)
	metricsService := metrics.NewMetricsService(cfg, collectorService.GetSpecifications)
	collectorService.SetMetricsConsumer(metricsService)
	return &Service{
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

func (s *Service) GetMetricsService() *metrics.MetricsService {
	return s.metricsService
}

func (s *Service) GetCollectorService() *collector.CollectorService {
	return s.collectorService
}
