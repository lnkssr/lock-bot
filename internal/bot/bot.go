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
	b.api.Handle("/logout", b.logoutHandler)
	b.api.Handle("/register", b.registerHandler)
	b.api.Handle(tele.OnDocument, b.uploadHandler)
	b.api.Handle("/storage", b.storageHandler)
	b.api.Handle("/delete", b.deleteHandler)
	b.api.Handle("/list_user", b.usersHandler)
	b.api.Handle("/download", b.downloadHandler)
}

func (b *Bot) Start() {
	b.api.Start()
}
