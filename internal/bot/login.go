package bot

import (
	"fmt"
	"lockbot/internal/api"
	"strings"
	"time"

	tele "gopkg.in/telebot.v4"
)

func (b *Bot) help(c tele.Context) error {
	return c.Send(`Help page: 
	/login <email password> - login in account
	/register <email name password> - reginster in account
	/logout - logout
	/profile - weiw your profile
	/storage - weiw your storage
	/delete <filename> - delete your file
	/download <flename> - download file`)
}

func (b *Bot) loginHandler(c tele.Context) error {
	args := c.Args()
	if len(args) < 2 {
		return c.Send("Usage: /logn email password")
	}

	email := args[0]
	password := strings.Join(args[1:], " ")

	resp, err := api.Login(email, password)
	if err != nil {
		return c.Send("Error: " + err.Error())
	}

	b.saveSession(c.Sender().ID, resp.Token, 24*time.Hour)

	reply := fmt.Sprintf(
		"%s\n Welcome: %s",
		resp.Message,
		resp.User.Email,
	)

	return c.Send(reply)
}

func (b *Bot) logoutHandler(c tele.Context) error {
	userID := c.Sender().ID

	token, ok := b.getSession(userID)
	if !ok {
		return c.Send("You are not authorized.")
	}

	err := api.Logout(token)
	if err != nil {
		return c.Send("Login error: " + err.Error())
	}

	delete(b.sessions, userID)

	return c.Send("You have successfully logged out of your account.")
}

func (b *Bot) registerHandler(c tele.Context) error {
	args := c.Args()
	if len(args) < 3 {
		return c.Send("Usage: /register email name password")
	}

	email := args[0]
	name := args[1]
	password := strings.Join(args[2:], " ")

	resp, err := api.Register(email, name, password)
	if err != nil {
		return c.Send("Registration error: " + err.Error())
	}

	reply := fmt.Sprintf(
		"%s\n You have successfully registered, now log in.",
		resp.Message,
	)

	return c.Send(reply)
}
