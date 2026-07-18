package util

import (
	"fmt"
	"io"
	"log"
	"os"

	"go.yaml.in/yaml/v3"
)

func ReadYaml(path string, out interface{}) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", path, err)
	}
	defer CloseWithLog(f)
	data, err := io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", path, err)
	}

	err = yaml.Unmarshal(data, out)
	if err != nil {
		return fmt.Errorf("failed to decode config: %w", err)
	}
	return nil
}

func MustReadYaml(path string, out interface{}) {
	err := ReadYaml(path, out)
	if err != nil {
		log.Panicf("failed to read yaml file %s: %s", path, err.Error())
	}
}
