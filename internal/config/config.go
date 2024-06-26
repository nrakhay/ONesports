package config

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

var (
	Token               string
	BotPrefix           string
	RecordingsChannelID string

	BucketName      string
	AccessKeyID     string
	SecretAccessKey string
	Region          string
)

type Config struct {
	Token               string `json:"token"`
	BotPrefix           string `json:"botPrefix"`
	RecordingsChannelID string `json: recordingsChannelID`
	BucketName          string `json:"bucketName"`
	AccessKeyID         string `json:"accessKeyID"`
	SecretAccessKey     string `json:"secret`
	Region              string `json:"region"`
}

func goDotEnvVariable(key string) string {
	err := godotenv.Load(".env")

	if err != nil {
		slog.Error("Error loading .env file")
	}

	return os.Getenv(key)
}

func ReadConfig() error {
	Token = goDotEnvVariable("BOT_TOKEN")
	BotPrefix = "!"

	// CHANGE THIS TO YOUR OWN CHANNEL ID
	RecordingsChannelID = goDotEnvVariable("RECORDINGS_CHANNEL_ID")

	BucketName = goDotEnvVariable("AWS_BUCKET_NAME")
	AccessKeyID = goDotEnvVariable("AWS_ACCESS_KEY")
	SecretAccessKey = goDotEnvVariable("AWS_SECRET_KEY")
	Region = goDotEnvVariable("AWS_REGION")

	return nil
}
