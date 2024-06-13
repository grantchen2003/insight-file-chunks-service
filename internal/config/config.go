package config

import (
	"errors"

	"github.com/joho/godotenv"
)

func LoadEnvVars(env string) error {
	switch env {

	case "dev", "development":
		return godotenv.Load(".env.dev")

	case "prod", "production":
		return godotenv.Load(".env.prod")

	default:
		return errors.New("no env vars found")
	}
}
