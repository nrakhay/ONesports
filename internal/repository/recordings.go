package repository

import (
	"github.com/nrakhay/ONEsports/internal/database"
)

func CreateVCRecording(channelID string, filePath string) error {
	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}

	statement, err := tx.Prepare("INSERT INTO voice_channel_recordings (channel_id, file_path) VALUES ($1, $2)")
	if err != nil {
		tx.Rollback()
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(channelID, filePath)
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
