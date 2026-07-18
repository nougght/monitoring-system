package util

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

func MustGetEnvVar(key string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	panic(fmt.Sprintf(`env variable "%s" not found`, key))
}

//nolint:unused
func GetOptionalEnvVar(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

//nolint:unused
func GetOptionalInt(key string, defaultValue int) int {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		log.Panicf(`env variable "%s" is not a valid integer: %s`, key, err)
	}
	return intValue
}
