package bot

import (
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/nrakhay/ONEsports/config"
	"github.com/pion/rtp"
	"github.com/pion/webrtc/v3/pkg/media"
	"github.com/pion/webrtc/v3/pkg/media/oggwriter"
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

func channelCreateHandler(s *discordgo.Session, c *discordgo.ChannelCreate) {
	if c.Type == discordgo.ChannelTypeGuildVoice {
		log.Printf("A new voice channel was created: %s", c.Name)
		vc, err := s.ChannelVoiceJoin(c.GuildID, c.ID, false, false)
		if err != nil {
			log.Printf("Failed to join the voice channel: %s", err)
			return
		}

		log.Println("Joined the new voice channel:", c.Name)
		go handleVoice(vc.OpusRecv)

		go func() {
			select {
			case <-time.After(10 * time.Second):
				vc.Disconnect()
			}
		}()
	}
}

func createPionRTPPacket(p *discordgo.Packet) *rtp.Packet {
	return &rtp.Packet{
		Header: rtp.Header{
			Version:        2,
			PayloadType:    0x78,
			SequenceNumber: p.Sequence,
			Timestamp:      p.Timestamp,
			SSRC:           p.SSRC,
		},
		Payload: p.Opus,
	}
}

func handleVoice(c chan *discordgo.Packet) {
	files := make(map[uint32]media.Writer)
	defer func() {
		for _, f := range files {
			f.Close()
		}
	}()

	for p := range c {
		log.Printf("Received packet from SSRC: %d", p.SSRC)
		file, ok := files[p.SSRC]
		if !ok {
			var err error
			file, err = oggwriter.New(fmt.Sprintf("%d.ogg", p.SSRC), 48000, 2)
			if err != nil {
				log.Printf("Failed to create file %d.ogg, giving up on recording: %v", p.SSRC, err)
				return
			}
			files[p.SSRC] = file
			log.Printf("Created new file for SSRC: %d", p.SSRC)
		}

		rtp := createPionRTPPacket(p)
		err := file.WriteRTP(rtp)
		if err != nil {
			log.Printf("Failed to write to file %d.ogg, giving up on recording: %v", p.SSRC, err)
		}
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
