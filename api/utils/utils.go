package utils

import (
	"fmt"
	"os"
)

func GetEnv(name string) string {
	return GetEnvWithFallback(name, "")
}

func GetEnvWithFallback(name string, fallback string) string {
	val := os.Getenv(name)

	if val == "" {
		if fallback == "" {
			panic(fmt.Sprintf("Missing required environment variable %s", name))
		}

		return fallback
	}

	return val
}
