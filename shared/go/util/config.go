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
	defer f.Close()

	defer CloseWithLog(f)
	data, err := io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", path, err)
	}

	err = yaml.Unmarshal(data, out)
	if err != nil {
		return fmt.Errorf("failed to decode yaml: %w", err)
	}
	return nil
}

func SaveYaml(path string, in interface{}) error {
	f, err := os.OpenFile(path, os.O_RDWR, 0644)

	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", path, err)
	}
	defer f.Close()

	defer CloseWithLog(f)
	encoded, err := yaml.Marshal(in)
	if err != nil {
		return fmt.Errorf("failed to encode yaml:%w", err)
	}

	_, err = f.Write(encoded)
	if err != nil {
		return fmt.Errorf("failed to write yaml to file: %w", err)
	}
	return nil
}

func MustReadYaml(path string, out interface{}) {
	err := ReadYaml(path, out)
	if err != nil {
		log.Panicf("failed to read yaml file %s: %s", path, err.Error())
	}
}
