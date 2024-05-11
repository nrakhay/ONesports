package main

import (
	"log/slog"

	"github.com/nrakhay/ONEsports/internal/bot"
	"github.com/nrakhay/ONEsports/internal/config"
	"github.com/nrakhay/ONEsports/internal/database"
	"github.com/nrakhay/ONEsports/internal/service/s3"
)

func main() {
	err := config.ReadConfig()

	if err != nil {
		slog.Error("Error occured reading config:", err)
		return
	}

	database.ConnectDB()
	database.PopulateDB()

	s3.StartS3Session()

	bot.Start()
}
