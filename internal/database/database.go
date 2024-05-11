package database

import (
	"log/slog"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

func ConnectDB() {
	var err error
	for i := 0; i < 4; i++ {
		connStr := "host=localhost port=5432 user=postgres password=admin dbname=on_esports_db sslmode=disable"

		DB, err = sqlx.Connect("postgres", connStr)
		if err != nil {
			slog.Warn("Failed to connect to database, retrying in 5 seconds", "error", err)
			time.Sleep(time.Second * 5)
		}
	}

	if err != nil {
		slog.Error("Failed to establish connection to database", "error", err)
	}

	DB.MustExec(schema)

	slog.Info("Successfully connected to database")
}
