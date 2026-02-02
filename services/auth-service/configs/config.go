package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Configs struct {
	JWT        JWT
	GRPC       GRPC
	UserGRPC   UserGRPC
	PostgreSQL PostgreSQL
	Mail       Mail
}

type GRPC struct {
	Port string
}

type UserGRPC struct {
	Host string
	Port string
}

type JWT struct {
	Secret string
}

type PostgreSQL struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
}

type Mail struct {
	Host   string
	Port   string
	Sender string
	Key    string
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
		JWT: JWT{
			Secret: requireEnv("JWT_SECRET"),
		},
		GRPC: GRPC{
			Port: requireEnv("GRPC_PORT"),
		},
		UserGRPC: UserGRPC{
			Host: requireEnv("User_HOST"),
			Port: requireEnv("User_PORT"),
		},
		PostgreSQL: PostgreSQL{
			Host:     requireEnv("DB_HOST"),
			Port:     requireEnv("DB_PORT"),
			Username: requireEnv("DB_USER"),
			Password: requireEnv("DB_PASSWORD"),
			Database: requireEnv("DB_NAME"),
		},
		Mail: Mail{
			Host:   requireEnv("EMAIL_HOST"),
			Port:   requireEnv("EMAIL_PORT"),
			Sender: requireEnv("EMAIL_USER"),
			Key:    requireEnv("EMAIL_PASS"),
		},
	}
}
