package discord

import (
	"log/slog"

	"github.com/bwmarrin/discordgo"

	"github.com/nrakhay/ONEsports/internal/config"
)

var Session *discordgo.Session

func InitSession() {
	var err error

	Session, err = discordgo.New("Bot " + config.Token)

	if err != nil {
		slog.Error("Failed to create discord session", "error", err)
		return
	}
}

func InitConnection() {
	if err := Session.Open(); err != nil {
		slog.Error("Failed to create websocket connection to discord", "error", err)
		return
	}
}
