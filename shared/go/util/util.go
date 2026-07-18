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
		log.Printf("failed to close file: %s", err.Error())
	}
}
