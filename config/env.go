package config

import (
	"errors"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func InitENV() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Println("Error loading .env file")
	}
}

func GetEnv(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return "", errors.New("environment variable not found")
	}
	return value, nil
}
