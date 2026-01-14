package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Configs struct {
	APP         APP
	JWT         JWT
	SUPA        SUPA
	UserService User
	AuthService Auth
}

type APP struct {
	Port string
}

type JWT struct {
	Secret string
}

type SUPA struct {
	Key    string
	URL    string
	Bucket string
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
			Secret: requireEnv("JWT_SECRET"),
		},
		SUPA: SUPA{
			Key:    requireEnv("SUPABASE_KEY"),
			URL:    requireEnv("SUPABASE_URL"),
			Bucket: requireEnv("BUCKET_NAME"),
		},
		UserService: User{
			Host: requireEnv("USER_HOST"),
			Port: requireEnv("USER_PORT"),
		},
		AuthService: Auth{
			Host: requireEnv("AUTH_HOST"),
			Port: requireEnv("AUTH_PORT"),
		},
	}
}
