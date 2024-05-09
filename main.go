package main

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/nrakhay/ONEsports/config"
	"github.com/nrakhay/ONEsports/database"
	"github.com/nrakhay/ONEsports/service/bot"
)

var db *sqlx.DB

func main() {
	err := config.ReadConfig()

	if err != nil {
		fmt.Println(err)
		return
	}

	db, err := connectDB()
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	// testing conn
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	} else {
		log.Println("Successfully Connected")
	}

	bot.Start()

	<-make(chan struct{})
	return
}

func connectDB() (*sqlx.DB, error) {
	connStr := "host=localhost port=5432 user=postgres password=admin dbname=on_esports_db sslmode=disable"
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, err
	}

	db.MustExec(database.Schema)

	return db, nil
}

func GetDB() *sqlx.DB {
	return db
}
