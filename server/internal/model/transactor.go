package model

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type Transactor interface {
	Begin(ctx context.Context) (pgx.Tx, error)
}
