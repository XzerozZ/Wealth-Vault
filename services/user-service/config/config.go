package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Configs struct {
	GRPC       GRPC
	PostgreSQL PostgreSQL
}

type GRPC struct {
	Port string
}

type PostgreSQL struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
	SSLMode  string
}

func LoadConfigs() *Configs {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, reading from environment variables")
	}

	return &Configs{
		GRPC: GRPC{
			Port: os.Getenv("GRPC_PORT"),
		},
		PostgreSQL: PostgreSQL{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			Username: os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Database: os.Getenv("DB_NAME"),
			SSLMode:  os.Getenv("SSL_Mode"),
		},
	}
}
