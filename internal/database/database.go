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
		connStr := "host=localhost port=50500 user=postgres password=admin dbname=on_esports_db sslmode=disable"

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

func PopulateDB() {
	tx := DB.MustBegin()

	users := []struct {
		FirstName   string
		LastName    string
		Email       string
		PhoneNumber string
	}{
		{"Daulet", "Kanatuly", "daulet@gmail.com", "+77075005050"},
		{"Nurali", "Rakhay", "nuralirakhay@gmail.com", "+77075336934"},
	}

	// check if already exists
	for _, user := range users {
		var exists bool
		err := tx.Get(&exists, "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 OR phone_number = $2)", user.Email, user.PhoneNumber)
		if err != nil {
			tx.Rollback() // rollback transaction
			return
		}

		// insert new record if doesn't exist
		if !exists {
			tx.MustExec("INSERT INTO users (first_name, last_name, email, phone_number) VALUES ($1, $2, $3, $4)",
				user.FirstName, user.LastName, user.Email, user.PhoneNumber)
		}
	}

	tx.Commit()
}
