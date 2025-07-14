package config

import (
	"bufio"
	"log"
	"os"
	"strings"
	"time"

	tele "gopkg.in/telebot.v4"
)

var (
	Pref     tele.Settings
	Api      string
	LogLevel string
)

func init() {
	loadEnvFile(".env")

	token := os.Getenv("TOKEN")
	if token == "" {
		log.Fatal("TOKEN is not set in environment")
	}

	Api = os.Getenv("API_URL")
	if Api == "" {
		log.Fatal("API_URL is not set in environment")
	}

	LogLevel := os.Getenv("LOG_LEVEL")
	if LogLevel == "" {
		log.Fatal("LOG_LEVEL is not set in env")
	}

	Pref = tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}
}

func loadEnvFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.Trim(strings.TrimSpace(parts[1]), `"'`)

		if _, exists := os.LookupEnv(key); !exists {
			os.Setenv(key, value)
		}
	}
}
