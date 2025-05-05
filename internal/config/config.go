package config

import (
	"TestTask/pkg/logger"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	URL struct {
		Age         string `yaml:"age"`
		Gender      string `yaml:"gender"`
		Nationality string `yaml:"nationality"`
	} `yaml:"url"`
}

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		logger.Logger.Fatal("Could not load env file!", err)
	}
}

func LoadYaml(path string) *Config {
	data, err := os.ReadFile(path)
	if err != nil {
		logger.Logger.Fatal("Could not load config file!", err)
		return nil
	}
	var cfg Config
	if err = yaml.Unmarshal(data, &cfg); err != nil {
		logger.Logger.Fatal("Could not load config file!", err)
		return nil
	}
	return &cfg
}
