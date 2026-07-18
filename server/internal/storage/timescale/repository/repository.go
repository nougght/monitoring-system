package repository

import "github.com/jackc/pgx/v5/pgxpool"

type Repositories struct {
}

func New(pgPool *pgxpool.Pool) *Repositories {
	return &Repositories{}
}
