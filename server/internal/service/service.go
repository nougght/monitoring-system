package service

import (
	"github.com/nougght/monitoring-system/server/internal/config"
	"github.com/nougght/monitoring-system/server/internal/storage/timescale/repository"
)

type Services struct {
}

//nolint:unused
type ServicesOptions struct {
	cfg          *config.Config
	Repositories *repository.Repositories
}

func New(opts ServicesOptions) *Services {
	return &Services{}
}
