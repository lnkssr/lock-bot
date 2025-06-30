package bot

import (
	"lockbot/internal/models"
	"strings"
	"time"
)

func (b *Bot) saveSession(userID int64, token string, ttl time.Duration) {
	b.sessions[userID] = models.SessionData{
		Token:     token,
		ExpiresAt: time.Now().Add(ttl),
	}
}

func (b *Bot) getSession(userID int64) (string, bool) {
	session, ok := b.sessions[userID]
	if !ok {
		return "", false
	}
	if time.Now().After(session.ExpiresAt) {
		delete(b.sessions, userID)
		return "", false
	}
	return session.Token, true
}

func formatFilesList(files []string) string {
	if len(files) == 0 {
		return "(empty)"
	}
	return "- " + strings.Join(files, "\n- ")
}
