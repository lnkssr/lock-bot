package bot

import (
	"fmt"
	"lockbot/internal/api"
	"strings"
	"time"

	tele "gopkg.in/telebot.v4"
)

func (b *Bot) help(c tele.Context) error {
	// TODO: help list command
	return c.Send("Help page")
}

func (b *Bot) loginHandler(c tele.Context) error {
	args := c.Args()
	if len(args) < 2 {
		return c.Send("Использование: /logn email password")
	}

	email := args[0]
	password := strings.Join(args[1:], " ")

	resp, err := api.Login(email, password)
	if err != nil {
		return c.Send("Ошибка: " + err.Error())
	}

	b.saveSession(c.Sender().ID, resp.Token, 24*time.Hour)

	reply := fmt.Sprintf(
		"%s\nПользователь: %s (ID: %d)\nТокен: %s",
		resp.Message,
		resp.User.Email,
		resp.User.ID,
		resp.Token,
	)

	return c.Send(reply)
}

func (b *Bot) profileHandler(c tele.Context) error {
	userID := c.Sender().ID

	token, ok := b.getSession(userID)
	if !ok {
		return c.Send("Вы не авторизованы. Введите /logn email password")
	}

	profileData, err := api.GetProfile(token)
	if err != nil {
		return c.Send("Ошибка получения профиля: " + err.Error())
	}

	return c.Send("Профиль:\n" + string(profileData))
}
