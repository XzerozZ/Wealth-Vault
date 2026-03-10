package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Configs struct {
	APP         APP
	PostgreSQL  PostgreSQL
	AuthService Auth
	NATS        NATS
	Line        Line
	FCM         FCM
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

type Auth struct {
	Host string
	Port string
}

type NATS struct {
	Host string
	Port string
}

type FCM struct {
	CredentialsFile string
}

type Line struct {
	Token string
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
		AuthService: Auth{
			Host: requireEnv("AUTH_HOST"),
			Port: requireEnv("AUTH_PORT"),
		},
		NATS: NATS{
			Host: requireEnv("NATS_HOST"),
			Port: requireEnv("NATS_PORT"),
		},
		Line: Line{
			Token: requireEnv("CHANNEL_ACCESS_TOKEN"),
		},
		FCM: FCM{
			CredentialsFile: requireEnv("FCM_CLIENT_FILE"),
		},
	}
}
