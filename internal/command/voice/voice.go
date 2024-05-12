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

	"github.com/pion/rtp"
	"github.com/pion/webrtc/v3/pkg/media/oggwriter"

	"github.com/bwmarrin/discordgo"
)

func saveChannelRecording(channelID string, buf *bytes.Buffer, key string) (s3Url string, err error) {
	channel, err := discord.Session.Channel(channelID)
	if err != nil {
		slog.Error("Failed to retrieve channel: ", err)
		return "", err
	}

	slog.Info("Uploading files to S3")

	s3Url, err = s3.UploadBufferToS3(buf, key)
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
	files := make(map[uint32]*oggwriter.OggWriter)
	buffers := make(map[uint32]*bytes.Buffer)
	defer func() {
		for ssrc, writer := range files {
			if writer != nil {
				err := writer.Close() // close the writer
				if err != nil {
					slog.Error("Failed to close OggWriter", "SSRC ID", ssrc, "error", err)
				}
			}
			buffer, ok := buffers[ssrc]
			if ok && buffer.Len() > 0 {
				key := fmt.Sprintf("%d-%s.ogg", ssrc, time.Now().Format("2006-01-02-15-04-05"))
				fileURL, err := saveChannelRecording(channelID, buffer, key)
				if err != nil {
					slog.Error("Failed to upload to S3", "error", err)
				} else {
					slog.Info("Successfully uploaded file to S3", "URL", fileURL)
				}
			}
		}
	}()

	for p := range c {
		slog.Info("Bot received packet from SSRC", "SSRC ID", p.SSRC)
		writer, ok := files[p.SSRC]
		if !ok {
			var err error
			buffer := new(bytes.Buffer)
			writer, err = oggwriter.NewWith(buffer, 48000, 2)
			if err != nil {
				slog.Error("Failed to initialize OggWriter", "error", err)
				return
			}
			files[p.SSRC] = writer
			buffers[p.SSRC] = buffer
			slog.Info("Initialized OggWriter for SSRC", "SSRC ID", p.SSRC)
		}
		rtp := createPionRTPPacket(p)
		err := writer.WriteRTP(rtp)
		if err != nil {
			slog.Error("Failed to write to OggWriter", "SSRC ID", p.SSRC, "error", err)
		}
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
