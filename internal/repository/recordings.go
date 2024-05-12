package repository

import (
	"github.com/nrakhay/ONEsports/internal/database"
)

func CreateVCRecording(channelID string, channelName string, filePath string) error {
	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}

	statement, err := tx.Prepare("INSERT INTO voice_channel_recordings (channel_id, channel_name, file_path) VALUES ($1, $2, $3)")
	if err != nil {
		tx.Rollback()
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(channelID, channelName, filePath)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
