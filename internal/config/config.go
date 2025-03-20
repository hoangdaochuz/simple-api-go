package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// config load env

type Config struct {
	DBHost		 string
	DBPort		 string
	DBUser		 string
	DBPassword	 string
	DBName		 string
	ServerPort	 string
}


func LoadConfig() *Config {
	// Load env
	err := godotenv.Load();
	if err != nil {
		log.Println("Error loading .env file")
	}

	return &Config{
		DBHost: os.Getenv("DB_HOST"),
		DBPort: os.Getenv("DB_PORT"),
		DBUser: os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASS"),
		DBName: os.Getenv("DB_NAME"),
		ServerPort: os.Getenv("SERVER_PORT"),
	}
}