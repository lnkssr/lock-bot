package bot

import (
	"fmt"
	"lockbot/internal/api"
	logger "lockbot/internal/log"
	"strconv"
	"strings"

	tele "gopkg.in/telebot.v4"
)

func (b *Bot) makeAdminHandler(c tele.Context) error {
	args := c.Args()
	requesterID := c.Sender().ID

	if len(args) < 1 {
		logger.Warn("makeAdmin: missing arguments", requesterID)
		return c.Send("Usage: /makeadmin <user_id>")
	}

	targetID, err := strconv.Atoi(args[0])
	if err != nil {
		logger.Warn("makeAdmin: invalid user ID", requesterID, args[0])
		return c.Send("Invalid user ID format.")
	}

	token, ok := b.getSession(requesterID)
	if !ok {
		logger.Warn("makeAdmin: unauthorized", requesterID)
		return c.Send("You are not authorized.")
	}

	logger.Debug("Attempting to promote user", requesterID, targetID)

	err = api.MakeAdmin(token, targetID)
	if err != nil {
		logger.Error("makeAdmin failed", requesterID, targetID, err)
		return c.Send("Error promoting user: " + err.Error())
	}

	logger.Info("User promoted to admin", requesterID, targetID)
	return c.Send(fmt.Sprintf("User %d is now an admin", targetID))
}

func (b *Bot) revokeAdminHandler(c tele.Context) error {
	args := c.Args()
	requesterID := c.Sender().ID

	if len(args) < 1 {
		logger.Warn("revokeAdmin: missing arguments", requesterID)
		return c.Send("Usage: /revokeadmin <user_id>")
	}

	targetID, err := strconv.Atoi(args[0])
	if err != nil {
		logger.Warn("revokeAdmin: invalid user ID", requesterID, args[0])
		return c.Send("User ID must be a number.")
	}

	token, ok := b.getSession(requesterID)
	if !ok {
		logger.Warn("revokeAdmin: unauthorized", requesterID)
		return c.Send("You are not authorized.")
	}

	logger.Debug("Attempting to revoke admin rights", requesterID, targetID)

	err = api.RevokeAdmin(token, targetID)
	if err != nil {
		logger.Error("revokeAdmin failed", requesterID, targetID, err)
		return c.Send("Error revoking admin rights: " + err.Error())
	}

	logger.Info("Admin rights revoked", requesterID, targetID)
	return c.Send(fmt.Sprintf("User %d is no longer an admin", targetID))
}

func (b *Bot) updateLimitHandler(c tele.Context) error {
	args := c.Args()
	requesterID := c.Sender().ID

	if len(args) < 2 {
		logger.Warn("updateLimit: missing arguments", requesterID)
		return c.Send("Usage: /limit <user_id> <new_limit>")
	}

	targetID, err := strconv.Atoi(args[0])
	if err != nil {
		logger.Warn("updateLimit: invalid user_id", requesterID, args[0])
		return c.Send("User ID must be a number.")
	}

	newLimit, err := strconv.Atoi(args[1])
	if err != nil {
		logger.Warn("updateLimit: invalid new_limit", requesterID, args[1])
		return c.Send("New limit must be a number.")
	}

	token, ok := b.getSession(requesterID)
	if !ok {
		logger.Warn("updateLimit: unauthorized", requesterID)
		return c.Send("You are not authorized.")
	}

	logger.Debug("Updating user limit", requesterID, targetID, newLimit)

	err = api.UpdateUserLimit(token, targetID, newLimit)
	if err != nil {
		logger.Error("updateLimit failed", requesterID, targetID, newLimit, err)
		return c.Send("Error updating user limit: " + err.Error())
	}

	logger.Info("User limit updated", requesterID, targetID, newLimit)
	return c.Send(fmt.Sprintf("User %d's limit updated to %d MB", targetID, newLimit))
}

func (b *Bot) usersHandler(c tele.Context) error {
	requesterID := c.Sender().ID

	token, ok := b.getSession(requesterID)
	if !ok {
		logger.Warn("usersHandler: unauthorized access", requesterID)
		return c.Send("You are not logged in. Please login with /login")
	}

	logger.Debug("Fetching user list", requesterID)

	users, err := api.GetAllUsers(token)
	if err != nil {
		logger.Error("usersHandler: failed to retrieve users", requesterID, err)
		return c.Send("Error retrieving the list of users: " + err.Error())
	}

	if len(users) == 0 {
		logger.Info("usersHandler: no users found", requesterID)
		return c.Send("No users found.")
	}

	var sb strings.Builder
	sb.WriteString("Users:\n")
	for _, u := range users {
		sb.WriteString(fmt.Sprintf("- %s (ID: %d)\n", u.Email, u.ID))
	}

	logger.Info("usersHandler: user list retrieved", requesterID, len(users))
	return c.Send(sb.String())
}
