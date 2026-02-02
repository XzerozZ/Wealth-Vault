package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Configs struct {
	GRPC       GRPC
	PostgreSQL PostgreSQL
	SUPA       SUPA
	AssetGRPC  AssetGRPC
	Mail       Mail
	NATS       NATS
}

type GRPC struct {
	Port string
}

type SUPA struct {
	Key    string
	URL    string
	Bucket string
}

type PostgreSQL struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
}

type AssetGRPC struct {
	Host string
	Port string
}

type Mail struct {
	Host   string
	Port   string
	Sender string
	Key    string
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
		GRPC: GRPC{
			Port: requireEnv("GRPC_PORT"),
		},
		PostgreSQL: PostgreSQL{
			Host:     requireEnv("DB_HOST"),
			Port:     requireEnv("DB_PORT"),
			Username: requireEnv("DB_USER"),
			Password: requireEnv("DB_PASSWORD"),
			Database: requireEnv("DB_NAME"),
		},
		SUPA: SUPA{
			Key:    requireEnv("SUPABASE_KEY"),
			URL:    requireEnv("SUPABASE_URL"),
			Bucket: requireEnv("BUCKET_NAME"),
		},
		AssetGRPC: AssetGRPC{
			Host: requireEnv("ASSET_HOST"),
			Port: requireEnv("ASSET_PORT"),
		},
		Mail: Mail{
			Host:   requireEnv("EMAIL_HOST"),
			Port:   requireEnv("EMAIL_PORT"),
			Sender: requireEnv("EMAIL_USER"),
			Key:    requireEnv("EMAIL_PASS"),
		},
		NATS: NATS{
			Host: requireEnv("NATS_HOST"),
			Port: requireEnv("NATS_PORT"),
		},
	}
}
