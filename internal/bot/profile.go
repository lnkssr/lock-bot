package bot

import (
	"lockbot/internal/api"

	tele "gopkg.in/telebot.v4"
)

func (b *Bot) profileHandler(c tele.Context) error {
	userID := c.Sender().ID

	token, ok := b.getSession(userID)
	if !ok {
		return c.Send("You are not logged in, please try: /logn email password")
	}

	profileData, err := api.Profile(token)
	if err != nil {
		return c.Send("Profile retrieval error: " + err.Error())
	}

	return c.Send("Profile:\n" + string(profileData))
}
