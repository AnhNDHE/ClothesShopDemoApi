package config

import (
	"os"
)

type Config struct {
	DBUrl     string
	JWTSecret string
}

func LoadConfig() Config {
	dbUrl := "postgres://" +
		os.Getenv("DB_USER") + ":" +
		os.Getenv("DB_PASSWORD") + "@" +
		os.Getenv("DB_HOST") + ":" +
		os.Getenv("DB_PORT") + "/" +
		os.Getenv("DB_NAME") +
		"?sslmode=disable"

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default-secret-key" // For development, change in production
	}

	return Config{
		DBUrl:     dbUrl,
		JWTSecret: jwtSecret,
	}
}
