package bot

import (
	"fmt"
	"lockbot/internal/api"
	"strconv"
	"strings"

	tele "gopkg.in/telebot.v4"
)

func (b *Bot) makeAdminHandler(c tele.Context) error {
	args := c.Args()
	if len(args) < 1 {
		return c.Send("Usage: /makeadmin <user_id>")
	}

	userIDStr := args[0]
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return c.Send("Invalid user ID")
	}

	token, ok := b.getSession(c.Sender().ID)
	if !ok {
		return c.Send("You are not authorized.")
	}

	err = api.MakeAdmin(token, userID)
	if err != nil {
		return c.Send("Error when raising rights: " + err.Error())
	}

	return c.Send(fmt.Sprintf("User %d is now an admin", userID))
}

func (b *Bot) revokeAdminHandler(c tele.Context) error {
	args := c.Args()
	if len(args) < 1 {
		return c.Send("Usage: /revokeadmin <user_id>")
	}

	userID, err := strconv.Atoi(args[0])
	if err != nil {
		return c.Send("Incorrect format user_id")
	}

	token, ok := b.getSession(c.Sender().ID)
	if !ok {
		return c.Send("You are not authorized.")
	}

	err = api.RevokeAdmin(token, userID)
	if err != nil {
		return c.Send("License revocation error: " + err.Error())
	}

	return c.Send(fmt.Sprintf("User %d is no longer an admin", userID))
}

func (b *Bot) updateLimitHandler(c tele.Context) error {
	args := c.Args()
	if len(args) < 2 {
		return c.Send("Usage: /limit <user_id> <new_limit>")
	}

	userID, err := strconv.Atoi(args[0])
	if err != nil {
		return c.Send("user_id should be a number")
	}
	limit, err := strconv.Atoi(args[1])
	if err != nil {
		return c.Send("new_limit should be a number")
	}

	token, ok := b.getSession(c.Sender().ID)
	if !ok {
		return c.Send("You are not authorized.")
	}

	err = api.UpdateUserLimit(token, userID, limit)
	if err != nil {
		return c.Send("Limit update error: " + err.Error())
	}

	return c.Send(fmt.Sprintf("User limit %d updated to %d", userID, limit))
}

func (b *Bot) usersHandler(c tele.Context) error {
	userID := c.Sender().ID
	token, ok := b.getSession(userID)
	if !ok {
		return c.Send("You are not logged in. Please login with /login")
	}

	users, err := api.GetAllUsers(token)
	if err != nil {
		return c.Send("Error getting the list of users:" + err.Error())
	}

	if len(users) == 0 {
		return c.Send("No users found.")
	}

	var sb strings.Builder
	sb.WriteString("Users:\n")
	for _, u := range users {
		sb.WriteString(fmt.Sprintf("- %s (ID: %d)\n", u.Email, u.ID))
	}

	return c.Send(sb.String())
}
