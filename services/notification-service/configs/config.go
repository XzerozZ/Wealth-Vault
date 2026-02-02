package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Configs struct {
	APP        APP
	PostgreSQL PostgreSQL
	NATS       NATS
}

type APP struct {
	Port string
}

type PostgreSQL struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
}

type NATS struct {
	Host string
	Port string
}

func LoadConfigs() *Configs {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, reading from environment variables")
	}

	requireEnv := func(key string) string {
		value := os.Getenv(key)
		if value == "" {
			log.Fatalf("CRITICAL ERROR: Missing required environment variable '%s'", key)
		}
		return value
	}

	return &Configs{
		APP: APP{
			Port: requireEnv("NOTI_PORT"),
		},
		PostgreSQL: PostgreSQL{
			Host:     requireEnv("DB_HOST"),
			Port:     requireEnv("DB_PORT"),
			Username: requireEnv("DB_USER"),
			Password: requireEnv("DB_PASSWORD"),
			Database: requireEnv("DB_NAME"),
		},
		NATS: NATS{
			Host: requireEnv("NATS_HOST"),
			Port: requireEnv("NATS_PORT"),
		},
	}
}
