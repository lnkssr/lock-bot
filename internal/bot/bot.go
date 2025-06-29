package bot

import (
	"lockbot/internal/config"
	"lockbot/internal/models"

	tele "gopkg.in/telebot.v4"
)

type Bot struct {
	api      *tele.Bot
	sessions map[int64]models.SessionData
}

func NewBot() (*Bot, error) {
	api, err := tele.NewBot(config.Pref)
	if err != nil {
		return nil, err
	}

	b := &Bot{
		api:      api,
		sessions: make(map[int64]models.SessionData),
	}
	b.registerHandlers()
	return b, nil
}

func (b *Bot) registerHandlers() {
	b.api.Handle("/help", b.help)
	b.api.Handle("/login", b.loginHandler)
	b.api.Handle("/profile", b.profileHandler)
}

func (b *Bot) Start() {
	b.api.Start()
}
