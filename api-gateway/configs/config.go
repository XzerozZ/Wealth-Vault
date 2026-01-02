package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Configs struct {
	APP         APP
	UserService User
}

type APP struct {
	Port string
}

type User struct {
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
			Port: requireEnv("APP_PORT"),
		},
		UserService: User{
			Host: requireEnv("User_HOST"),
			Port: requireEnv("User_PORT"),
		},
	}
}
