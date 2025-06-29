package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	tele "gopkg.in/telebot.v4"
)

var (
	Pref tele.Settings
	Api  string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	token := os.Getenv("TOKEN")
	if token == "" {
		log.Fatal("TOKEN is not set in environment")
	}

	Api = os.Getenv("API_URL")
	if Api == "" {
		log.Fatal("API_URL is not set in environment")

	}

	Pref = tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}
}
