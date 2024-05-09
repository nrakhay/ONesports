package entity

type VoicChannelRecording struct {
	ID        int    `db:"id"`
	ChannelID string `db:"channel_id"`
	FilePath  string `db:"file_path"`
}
