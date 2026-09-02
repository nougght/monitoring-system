package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/nougght/monitoring-system/server/internal/model"
	metrics_model "github.com/nougght/monitoring-system/server/internal/model/metrics"
)

type SeriesRepository struct {
	database DB
}

func NewSeriesRepository(db DB) *SeriesRepository {
	return &SeriesRepository{
		database: db,
	}

}

func (r *SeriesRepository) db(ctx context.Context) DB {
	res := r.database
	if tx := ctx.Value(model.ContextKeyTx); tx != nil {
		res = tx.(DB)
	}
	return res
}

func (r *SeriesRepository) BatchCreateOrLoadSeries(ctx context.Context, seriesList []metrics_model.MetricSeriesKey) (resultSeries []metrics_model.MetricSeries, err error) {
	if len(seriesList) == 0 {
		return nil, nil
	}
	query := `
            INSERT INTO metric_series (agent_id, kind, label) VALUES ($1,$2,$3)
            ON CONFLICT (agent_id, kind, label) DO NOTHING
            RETURNING id
        `

	batch := &pgx.Batch{}
	for _, s := range seriesList {
		batch.Queue(query, s.AgentID, s.Kind, s.Label)
	}

	tx, err := r.db(ctx).Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed begin transaction: %w", err)
	}
	defer func() {
		err := tx.Rollback(ctx)
		if err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			log.Printf("rollback failed: %s", err.Error())
		}
	}()

	res := r.db(ctx).SendBatch(ctx, batch)
	defer func() {
		_ = res.Close()
	}()

	missingKeysWithIndexes := make(map[metrics_model.MetricSeriesKey]int, 0)
	for i := range seriesList {
		createdSeries := metrics_model.MetricSeries{
			MetricSeriesKey: seriesList[i],
		}
		err := res.QueryRow().Scan(&createdSeries.ID)
		if errors.Is(err, sql.ErrNoRows) {
			missingKeysWithIndexes[seriesList[i]] = i
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("batch QueryRow scan error: %w", err)
		}

		resultSeries = append(resultSeries, createdSeries)
	}

	loadedIDs, err := r.GetSeriesIDsByKeys(ctx, seriesList)
	if err != nil {
		return nil, fmt.Errorf("failed to load series ids: %w", err)
	}

	// fill ids for loaded series
	for key, id := range loadedIDs {
		ind, ok := missingKeysWithIndexes[key]
		if !ok {
			return nil, fmt.Errorf("loaded unknown metric key")
		}
		resultSeries = append(resultSeries, metrics_model.MetricSeries{
			ID:              metrics_model.MetricSeriesID(id),
			MetricSeriesKey: seriesList[ind],
		})
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("commit transaction error: %w", err)
	}
	return resultSeries, nil
}

func (r *SeriesRepository) GetSeriesIDsByKeys(ctx context.Context, seriesKeys []metrics_model.MetricSeriesKey) (map[metrics_model.MetricSeriesKey]int64, error) {
	query := `
		SELECT s.id, s.agent_id, s.kind, s.label
		FROM metric_series s
		JOIN unnest($1::uuid[], $2::smallint[], $3::text[]) AS u(agent_id, kind, label)
  		ON s.agent_id = u.agent_id AND s.kind = u.kind AND s.label = u.label;`
	agentIDList, kindList, labelList := make([]uuid.UUID, len(seriesKeys)), make([]int32, len(seriesKeys)), make([]string, len(seriesKeys))
	for i, key := range seriesKeys {
		agentIDList[i] = key.AgentID
		labelList[i] = key.Label
		kindList[i] = key.Kind
	}
	rows, err := r.db(ctx).Query(ctx, query, agentIDList, kindList, labelList)
	if err != nil {
		return nil, fmt.Errorf("select failed: %w", err)
	}
	defer func() {
		rows.Close()
	}()

	var id int64
	var key metrics_model.MetricSeriesKey
	res := make(map[metrics_model.MetricSeriesKey]int64, len(seriesKeys))
	for rows.Next() {
		err = rows.Scan(&id, &key.AgentID, &key.Kind, &key.Label)
		if err != nil {
			return nil, fmt.Errorf("rows scan failed: %w", err)
		}
		res[key] = id
	}
	return res, nil
}

func (r *SeriesRepository) GetAllSeries(ctx context.Context) ([]*metrics_model.MetricSeries, error) {
	query := `SELECT * FROM metric_series`
	rows, err := r.db(ctx).Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("select failed: %w", err)
	}
	defer rows.Close()
	series, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[metrics_model.MetricSeries])
	if err != nil {
		return nil, fmt.Errorf("collect rows failed: %w", err)
	}
	return series, nil
}
