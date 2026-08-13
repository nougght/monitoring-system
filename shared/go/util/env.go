package util

import (
	"log"
	"os"
	"strconv"
)

func MustGetEnvVar(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Panicf(`env variable "%s" not found`, key)
	}
	return value
}

//nolint:unused
func GetOptionalEnvVar(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func MustGetIntEnvVar(key string) int {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Panicf(`env variable "%s" not found`, key)
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		log.Panicf(`env variable "%s" is not a valid integer: %s`, key, err.Error())
	}
	return intValue
}

//nolint:unused
func GetOptionalIntEnvVar(key string, defaultValue int) int {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		log.Panicf(`env variable "%s" is not a valid integer: %s`, key, err.Error())
	}
	return intValue
}

func MustGetPortEnvVar(key string) int {
	value := MustGetIntEnvVar(key)
	if value <= 1024 || value > 65535 {
		log.Panicf("port must be in range 1025-65535")
	}
	return value
}
