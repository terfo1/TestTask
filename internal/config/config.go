package config

import (
	"TestTask/pkg/logger"
	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		logger.Logger.Fatal("Could not load env file!", err)
	}
}
