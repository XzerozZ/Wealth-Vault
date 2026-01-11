package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Configs struct {
	APP         APP
	JWT         JWT
	UserService User
	AuthService Auth
}

type APP struct {
	Port string
}

type JWT struct {
	Secret string
}

type User struct {
	Host string
	Port string
}

type Auth struct {
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
		JWT: JWT{
			Secret: os.Getenv("JWT_SECRET"),
		},
		UserService: User{
			Host: requireEnv("User_HOST"),
			Port: requireEnv("User_PORT"),
		},
		AuthService: Auth{
			Host: requireEnv("Auth_HOST"),
			Port: requireEnv("Auth_PORT"),
		},
	}
}
