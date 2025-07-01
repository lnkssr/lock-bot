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
	b.registerAdminHandlers()
	return b, nil
}

func (b *Bot) registerAdminHandlers() {
	b.api.Handle("/users", b.usersHandler)
	b.api.Handle("/limit", b.updateLimitHandler)
	b.api.Handle("/makeadmin", b.makeAdminHandler)
	b.api.Handle("/revokeadmin", b.revokeAdminHandler)
}

func (b *Bot) registerHandlers() {
	b.api.Handle("/help", b.help)
	b.api.Handle("/login", b.loginHandler)
	b.api.Handle("/logout", b.logoutHandler)
	b.api.Handle("/delete", b.deleteHandler)
	b.api.Handle("/storage", b.storageHandler)
	b.api.Handle("/profile", b.profileHandler)
	b.api.Handle("/register", b.registerHandler)
	b.api.Handle("/download", b.downloadHandler)
	b.api.Handle(tele.OnDocument, b.uploadHandler)
}

func (b *Bot) Start() {
	b.api.Start()
}
