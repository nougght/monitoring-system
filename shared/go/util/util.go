package util

import (
	"io"
	"log"
)

func CloseWithLog(closer io.Closer) {
	if closer == nil {
		return
	}
	err := closer.Close()
	if err != nil {
		log.Printf("failed to close: %s", err.Error())
	}
}

func Map[T any, R any](slice []T, f func(T) R, args ...any) []R {
	result := make([]R, len(slice))
	for i, v := range slice {
		result[i] = f(v)
	}
	return result
}

func Ptr[T any](value T) *T {
	return &value
}
