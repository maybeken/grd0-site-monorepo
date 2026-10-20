package utils

import (
	"fmt"
	"os"
	"strings"
)

func GetEnv(name string) string {
	return GetEnvWithFallback(name, "")
}

func GetEnvOrEmpty(name string) string {
	return os.Getenv(name)
}

func IsLocalEnvironment() bool {
	return strings.EqualFold(strings.TrimSpace(os.Getenv("APP_ENV")), "LOCAL")
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
