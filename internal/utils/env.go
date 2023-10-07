package utils

import (
	"os"

	"github.com/joho/godotenv"
)

var defaultEnvVars = map[string]string{
	"PG_USER":     "postgres",
	"PG_PASSWORD": "postgres",
	"PG_DATABASE": "postgres",
	"PG_HOST":     "localhost",
	"PG_PORT":     "5432",
	"PG_SSLMODE":  "disable",
}

func setEnvVars(envVars map[string]string) {
	for key, value := range envVars {
		if os.Getenv(key) == "" {
			os.Setenv(key, value)
		}
	}
}

func PrepareEnv() error {
	err := godotenv.Load()

	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	}

	setEnvVars(defaultEnvVars)
	return nil
}
