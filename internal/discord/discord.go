package discord

import (
	"log/slog"
)

func SendChannelMessage(channelID string, message string) {
	_, err := Session.ChannelMessageSend(channelID, message)
	if err != nil {
		slog.Warn("Failed to send message to channel", "channelId", channelID, "message", message, "error", err)
	}
}
