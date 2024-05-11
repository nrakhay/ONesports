package bot

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/nrakhay/ONEsports/internal/discord"
	"github.com/nrakhay/ONEsports/internal/handlers"
)

func Start() {
	discord.InitSession()
	addHandlers()
	discord.InitConnection()

	defer discord.Session.Close()

	slog.Info("Bot is running!")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sc

	slog.Info("Shutting down.", "signal", sig)
}

func addHandlers() {
	discord.Session.AddHandler(handlers.MessageHandler)
	discord.Session.AddHandler(handlers.ChannelCreateHandler)
	discord.Session.AddHandler(handlers.VoiceStateUpdateHandler)
}
