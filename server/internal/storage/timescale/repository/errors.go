package repository

import "errors"

var (
	ErrNotFound       = errors.New("repository error not found")
	ErrNoAffectedRows = errors.New("repository error no affected rows")
)
