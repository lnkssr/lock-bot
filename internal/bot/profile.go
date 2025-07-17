package bot

import (
	"lockbot/internal/api"
	logger "lockbot/internal/log"

	tele "gopkg.in/telebot.v4"
)

func (b *Bot) profileHandler(c tele.Context) error {
	userID := c.Sender().ID

	token, ok := b.getSession(userID)
	if !ok {
		logger.Warn("Profile access attempt without login", userID)
		return c.Send("You are not logged in. Please use: /login <email> <password>")
	}

	logger.Debug("Fetching profile", userID)

	profileData, err := api.Profile(token)
	if err != nil {
		logger.Error("Profile retrieval failed", userID, err)
		return c.Send("Profile retrieval error: " + err.Error())
	}

	logger.Info("Profile retrieved successfully", userID)
	return c.Send("Profile:\n" + string(profileData))
}

func (b *Bot) profileHandlerV2(c tele.Context) error {
	userID := c.Sender().ID

	token, ok := b.getSession(userID)
	if !ok {
		logger.Warn("Profile access attempt without login", userID)
		return c.Send("You are not logged in. Please use: /login <email> <password>")
	}

	logger.Debug("Fetching profile", userID)

	profileData, err := api.ProfileV2(token)
	if err != nil {
		logger.Error("Profile retrieval failed", userID, err)
		return c.Send("Profile retrieval error: " + err.Error())
	}

	logger.Info("Profile retrieved successfully", userID)
	return c.Send("Profile:\n" + string(profileData))
}
