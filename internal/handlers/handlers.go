package handlers

import (
	"log/slog"
	"strings"
	"time"

	"github.com/nrakhay/ONEsports/internal/command/voice"
	"github.com/nrakhay/ONEsports/internal/discord"

	"github.com/bwmarrin/discordgo"
)

func MessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// ignore messages from the bot itself to avoid loops
	if m.Author.ID == s.State.User.ID {
		return
	}

	slog.Info("Received message:", "content: ", m.Content)

	//  command patterns
	directCommand := "!"
	idCommandWithBang := "<@" + s.State.User.ID + "> !"
	idCommandWithoutBang := "<@" + s.State.User.ID + ">" + " "

	// check for command patterns and extract the command
	var cmd string
	if strings.HasPrefix(m.Content, directCommand) {
		cmd = strings.TrimPrefix(m.Content, directCommand)
	} else if strings.HasPrefix(m.Content, idCommandWithBang) {
		cmd = strings.TrimPrefix(m.Content, idCommandWithBang)
	} else if strings.HasPrefix(m.Content, idCommandWithoutBang) {
		cmd = strings.TrimPrefix(m.Content, idCommandWithoutBang)
	} else {
		return
	}

	// split into base command and arguments
	cmdArgs := strings.Fields(cmd)
	if len(cmdArgs) == 0 {
		return // if no command after prefix
	}

	slog.Info("Processing command", "command: ", cmdArgs[0])

	// Handle commands
	switch cmdArgs[0] {
	case "ping":
		discord.SendChannelMessage(m.ChannelID, "pong!")
	}
}

// Listen to channel create events
func ChannelCreateHandler(s *discordgo.Session, c *discordgo.ChannelCreate) {
	if c.Type == discordgo.ChannelTypeGuildVoice {
		slog.Info("A new voice channel was created.", "Name of channel", c.Name)
		vc, err := s.ChannelVoiceJoin(c.GuildID, c.ID, false, false)
		if err != nil {
			slog.Error("Failed to join the voice channel", "VC ID", err)
			return
		}

		go voice.HandleVoice(vc.OpusRecv)

		go func() {
			select {
			case <-time.After(5 * time.Second):
				vc.Disconnect()
			}
		}()
	}
}

// Listen to voice state update events
func VoiceStateUpdateHandler(s *discordgo.Session, vs *discordgo.VoiceStateUpdate) {
	if vs.UserID == s.State.User.ID {
		if vs.ChannelID == "" {
			slog.Info("Bot has left the voice channel")
			voice.OnBotLeaveVoiceChannel(vs.ChannelID)
		} else {
			slog.Info("Bot is now in voice channel.", "VC ID", vs.ChannelID)
		}
	}
}
