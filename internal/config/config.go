package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	DbUrl  string
	ApiKey string
}

func LoadEnv() Env {
	log.Println("Loading environment configuration...")

	err := godotenv.Load()
	if err != nil {
		log.Printf("WARNING: .env file not found or could not be loaded: %v", err)
		log.Println("Continuing with system environment variables...")
	} else {
		log.Println("SUCCESS: .env file loaded successfully")
	}

	DbUrl := os.Getenv("DB_URL")
	ApiKey := os.Getenv("API_KEY")

	if DbUrl == "" {
		log.Println("ERROR: DB_URL environment variable is not set")
	} else {
		log.Printf("SUCCESS: DB_URL found (length: %d characters)", len(DbUrl))
	}

	if ApiKey == "" {
		log.Println("ERROR: API_KEY environment variable is not set")
	} else {
		log.Printf("SUCCESS: API_KEY found (length: %d characters)", len(ApiKey))
	}

	env := Env{
		DbUrl:  DbUrl,
		ApiKey: ApiKey,
	}

	log.Println("Environment configuration loaded successfully")
	return env
}
