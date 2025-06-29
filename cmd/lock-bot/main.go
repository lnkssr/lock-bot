package main

import (
	"lockbot/internal/bot"
	"log"
)

func main() {
	bot, err := bot.NewBot()
	if err != nil {
		log.Fatal(err)
		return
	}

	bot.Start()
}
