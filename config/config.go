package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port       string
	MongoDBURI string
	DBName     string
	JWTSecret  string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	return &Config{
		Port:       os.Getenv("PORT"),
		MongoDBURI: os.Getenv("MONGODB_URI"),
		DBName:     os.Getenv("DB_NAME"),
		JWTSecret:  os.Getenv("JWT_SECRET"),
	}
}
