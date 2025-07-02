package bot

import (
	"fmt"
	"lockbot/internal/api"
	logger "lockbot/internal/log"
	"strings"
	"time"

	tele "gopkg.in/telebot.v4"
)

func (b *Bot) help(c tele.Context) error {
	return c.Send(`Help page:
/login <email password> - login to your account
/register <email name password> - register a new account
/logout - logout
/profile - view your profile
/storage - view your storage
/delete <filename> - delete a file
/download <filename> - download a file`)
}

func (b *Bot) loginHandler(c tele.Context) error {
	args := c.Args()
	userID := c.Sender().ID

	if len(args) < 2 {
		logger.Warn("Login attempt with insufficient arguments", userID)
		return c.Send("Usage: /login <email> <password>")
	}

	email := args[0]
	password := strings.Join(args[1:], " ")

	logger.Debug("Attempting login", userID, email)

	resp, err := api.Login(email, password)
	if err != nil {
		logger.Error("Login failed", userID, email, err)
		return c.Send("Login error: " + err.Error())
	}

	b.saveSession(userID, resp.Token, 24*time.Hour)
	logger.Info("Login successful", userID, email)

	reply := fmt.Sprintf(
		"%s\nWelcome, %s!",
		resp.Message,
		resp.User.Email,
	)

	return c.Send(reply)
}

func (b *Bot) logoutHandler(c tele.Context) error {
	userID := c.Sender().ID

	token, ok := b.getSession(userID)
	if !ok {
		logger.Warn("Logout attempt without active session", userID)
		return c.Send("You are not authorized.")
	}

	err := api.Logout(token)
	if err != nil {
		logger.Error("Logout failed", userID, err)
		return c.Send("Logout error: " + err.Error())
	}

	delete(b.sessions, userID)
	logger.Info("Logout successful", userID)

	return c.Send("You have successfully logged out of your account.")
}

func (b *Bot) registerHandler(c tele.Context) error {
	args := c.Args()
	userID := c.Sender().ID

	if len(args) < 3 {
		logger.Warn("Register attempt with insufficient arguments", userID)
		return c.Send("Usage: /register <email> <name> <password>")
	}

	email := args[0]
	name := args[1]
	password := strings.Join(args[2:], " ")

	logger.Debug("Attempting registration", userID, email, name)

	resp, err := api.Register(email, name, password)
	if err != nil {
		logger.Error("Registration failed", userID, email, err)
		return c.Send("Registration error: " + err.Error())
	}

	logger.Info("Registration successful", userID, email)

	reply := fmt.Sprintf(
		"%s\nYou have successfully registered. Now you can log in.",
		resp.Message,
	)

	return c.Send(reply)
}
