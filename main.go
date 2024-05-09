package main

import (
	"fmt"

	"github.com/nrakhay/ONEsports/config"
	"github.com/nrakhay/ONEsports/service/bot"
)

func main() {
	err := config.ReadConfig()

	if err != nil {
		fmt.Println(err)
		return
	}

	bot.Start()

	<-make(chan struct{})
	return
}
