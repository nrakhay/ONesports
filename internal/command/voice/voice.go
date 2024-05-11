package voice

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/nrakhay/ONEsports/internal/command/text"
	"github.com/nrakhay/ONEsports/internal/repository"
	"github.com/nrakhay/ONEsports/internal/service/s3"

	"github.com/bwmarrin/discordgo"

	"github.com/pion/rtp"
	"github.com/pion/webrtc/v3/pkg/media"
	"github.com/pion/webrtc/v3/pkg/media/oggwriter"
)

func OnBotLeaveVoiceChannel(channelID string) {
	dir := "recordings"

	files, err := os.ReadDir(dir)
	if err != nil {
		slog.Error("Failed to read directory recordings: %v", err)
		return
	}

	// if no files to upload -> return
	if len(files) == 0 {
		return
	}

	slog.Info("Uploading files to S3")

	for _, f := range files {
		fileName := filepath.Join(dir, f.Name())
		s3Url, err := s3.UploadFileToS3(fileName)
		if err != nil {
			slog.Error("Failed to upload %s to S3: %v", fileName, err)
			continue
		}

		slog.Info("Successfully uploaded to S3.", "Filename", f.Name())

		err = repository.CreateVCRecording(channelID, s3Url)
		if err != nil {
			slog.Error("Failed to create recording in database", "error", err)
			continue
		}

		// this is text channel id for recordings
		text.SendVoiceRecordingToTextChannel("1238876913070637198", f.Name())

		err = os.Remove(fileName)
		if err != nil {
			slog.Info("Failed to delete %s: %v", fileName, err)
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

func HandleVoice(c chan *discordgo.Packet) {
	files := make(map[uint32]media.Writer)
	defer func() {
		for _, f := range files {
			f.Close()
		}
	}()

	err := os.MkdirAll("recordings", os.ModePerm)
	if err != nil {
		slog.Error("Failed to create recordings directory", "error", err)
		return
	}

	var filename string

	for p := range c {
		slog.Info("Bot received packet from SSRC", "SSRC ID", p.SSRC)
		file, ok := files[p.SSRC]
		if !ok {
			var err error
			filename = fmt.Sprintf("recordings/%d-%s.ogg", p.SSRC, time.Now().Format("2006-01-02-15-04-05"))
			file, err = oggwriter.New(filename, 48000, 2)
			if err != nil {
				slog.Error("Failed to create file", "path", filename, "error", err)
				return
			}
			files[p.SSRC] = file
			slog.Info("Created new file for SSRC", "SSRC ID", p.SSRC, "Filename", filename)
		}
		rtp := createPionRTPPacket(p)
		err := file.WriteRTP(rtp)

		if err != nil {
			slog.Error("Failed to write to file", "Filename", filename, "error", err)
		}
	}
}
