package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/nougght/monitoring-system/server/internal/model"
	agent_model "github.com/nougght/monitoring-system/server/internal/model/agent"
)

type SpecsRepository struct {
	pool DB
}

func NewSpecsRepository(db DB) *SpecsRepository {
	return &SpecsRepository{
		pool: db,
	}

}

func (r *SpecsRepository) db(ctx context.Context) DB {
	res := r.pool
	if tx := ctx.Value(model.ContextKeyTx); tx != nil {
		res = tx.(DB)
	}
	return res
}

func (r *SpecsRepository) CreateOrUpdateSpecs(ctx context.Context, specs *agent_model.Specs) (*agent_model.Specs, error) {
	query := `
	INSERT INTO agent_specs (agent_id, hostname, os_type, os, os_arch, cpu_cores_count, memory_total, full_specs, updated_at) 
	VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)
	RETURNING updated_at
	`
	err := r.db(ctx).QueryRow(ctx, query,
		specs.AgentID,
		specs.HostSpecs.Hostname,
		specs.HostSpecs.OSType,
		specs.HostSpecs.OS,
		specs.HostSpecs.OSArch,
		specs.CpuSpecs.NumberOfCores,
		specs.MemorySpecs.Total,
		specs,
		specs.UpdatedAt).Scan(&specs.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("insert failed: %w", err)
	}

	return specs, nil
}

func (r *SpecsRepository) GetCurrentSpecs(ctx context.Context, agentID uuid.UUID) (specs *agent_model.Specs, err error) {
	query := `
		SELECT agent_id, full_specs, updated_at FROM agent_specs WHERE agent_id = $1;
		`
	rows, err := r.db(ctx).Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("select failed: %w", err)
	}
	defer rows.Close()

	specs = &agent_model.Specs{}
	err = rows.Scan(&specs.AgentID, &specs, &specs.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("scan failed: %w", err)
	}

	return specs, nil
}
