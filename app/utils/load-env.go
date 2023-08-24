package utils

import (
	"github.com/joho/godotenv"
)

// LoadEnvVariables loads the environment variables from the .env file
func LoadEnvVariables() error {
	err := godotenv.Load(".env")
	if err != nil {
		return err
	}

	return nil
}
