package utils

import (
	"log"
	"os"
)

func RequireEnv(name string) string {
	val := os.Getenv(name)
	if val == "" {
		log.Fatalf("Required environment variable %s not set", name)
	}

	return val
}
