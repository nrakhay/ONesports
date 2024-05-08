package config

import (
	"log"
	"os"
	"github.com/joho/godotenv"
)

var (
	Token     string
	BotPrefix string

	config *Config
)

type Config struct {
	Token     string `json:"token"`
	BotPrefix string `json:"botPrefix"`
}

func goDotEnvVariable(key string) string {
	err := godotenv.Load(".env")
  
	if err != nil {
	  log.Fatalf("Error loading .env file")
	}
  
	return os.Getenv(key)
  }

func ReadConfig() error {
	Token = goDotEnvVariable("BOT_TOKEN")
	BotPrefix = "!"

	return nil
}