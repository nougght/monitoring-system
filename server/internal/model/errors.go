package model

import "errors"

var (
	ErrNotFound   = errors.New("error not found")
	ErrBadRequest = errors.New("error bad request")
)
