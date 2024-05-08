package bot

import (
	"fmt"

	"log"

	"time"

	"github.com/nrakhay/ONEsports/config"

	"github.com/bwmarrin/discordgo"
)

var BotId string
var goBot *discordgo.Session

func Start() {
	goBot, err := discordgo.New("Bot " + config.Token)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	user, err := goBot.User("@me")

	if err != nil {
		fmt.Println(err.Error())
	}

	BotId = user.ID

	goBot.AddHandler(messageHandler)

	goBot.AddHandler(channelCreateHandler)

	err = goBot.Open()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Bot is running!")
}

func channelCreateHandler(session *discordgo.Session, channel *discordgo.ChannelCreate) {
	if channel.Type == discordgo.ChannelTypeGuildVoice {
		log.Printf("A new voice channel was created: %s", channel.Name)

		// Join the voice channel
		vc, err := session.ChannelVoiceJoin(channel.GuildID, channel.ID, false, true)
		if err != nil {
			log.Printf("Failed to join the voice channel: %s", err)
			return
		}

		
		log.Println("Joined the new voice channel:", channel.Name)

		// TODO: Start capturing audio from the voice channel

		// Disconnect after 10 seconds
		go func() {
			select {
			case <-time.After(10 * time.Second):
				vc.Disconnect()
			}
		}()
	}
}

func messageHandler(session *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.ID == BotId {
		return
	}

	if message.Content == "<@"+BotId+"> !ping" || message.Content == "<@"+BotId+"> ping" {
		_, _ = session.ChannelMessageSend(message.ChannelID, "pong!")
	}

	if message.Content == config.BotPrefix+"ping" {
		_, _ = session.ChannelMessageSend(message.ChannelID, "pong!")
	}
}