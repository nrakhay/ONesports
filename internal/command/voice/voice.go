package voice

import (
	"bytes"
	"fmt"
	"log/slog"
	"time"

	"github.com/nrakhay/ONEsports/internal/command/text"
	"github.com/nrakhay/ONEsports/internal/config"
	"github.com/nrakhay/ONEsports/internal/discord"
	"github.com/nrakhay/ONEsports/internal/repository"
	"github.com/nrakhay/ONEsports/internal/service/s3"

	"github.com/bwmarrin/discordgo"
)

func saveChannelRecording(channelID string, buffer *bytes.Buffer, key string) (s3Url string, err error) {
	channel, err := discord.Session.Channel(channelID)
	if err != nil {
		slog.Error("Failed to retrieve channel: ", err)
		return "", err
	}

	slog.Info("Uploading files to S3")

	s3Url, err = s3.UploadBufferToS3(buffer, key)
	if err != nil {
		slog.Error("Failed to upload to S3", "key", key, "error", err)
		return "", err
	}
	slog.Info("Successfully uploaded to S3.", "Filename", key)

	err = repository.CreateVCRecording(channelID, channel.Name, s3Url)
	if err != nil {
		slog.Error("Failed to create recording in database", "error", err)
		return "", err
	}

	slog.Info("Successfully created recording in database", "Filename", key)

	// this is text channel id for recordings
	text.SendVoiceRecordingToTextChannel(config.RecordingsChannelID, channel.Name, key)

	return s3Url, nil
}

func HandleVoice(c chan *discordgo.Packet, channelID string) {
	buffers := make(map[uint32]*bytes.Buffer)

	defer func() {
		for ssrc, buf := range buffers {
			key := fmt.Sprintf("%d-%s.ogg", ssrc, time.Now().Format("2006-01-02-15-04-05"))

			s3Url, err := saveChannelRecording(channelID, buf, key)
			if err != nil {
				slog.Error("Failed to send SSRC data to S3", "ssrc", ssrc, "error", err)
				continue
			}

			slog.Info("All SSRC data has been saved", "S3 URL", s3Url)
		}
	}()

	for p := range c {
		slog.Info("Bot received packet from", "SSRC", p.SSRC)
		buffer, ok := buffers[p.SSRC]
		if !ok {
			buffer = new(bytes.Buffer)
			buffers[p.SSRC] = buffer
			slog.Info("Created new buffer for", "SSRC", p.SSRC)
		}

		if _, err := buffer.Write(p.Opus); err != nil {
			slog.Error("Failed to write to buffer for", "SSRC", p.SSRC, "error", err)
		}
	}
}
