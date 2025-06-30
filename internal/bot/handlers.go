package bot

import (
	"bytes"
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

func (b *Bot) uploadHandler(c tele.Context) error {
	file := c.Message().Document
	if file == nil {
		return c.Send("Please send the file as a document.")
	}

	reader, err := b.api.File(&file.File)
	if err != nil {
		return c.Send("Error while retrieving a file: " + err.Error())
	}
	defer reader.Close()

	userID := c.Sender().ID
	token, ok := b.getSession(userID)
	if !ok {
		return c.Send("You are not logged in. Please login with /login")
	}

	_, err = api.UploadFile(token, file.FileName, reader)
	if err != nil {
		return c.Send("File upload error:" + err.Error())
	}

	return c.Send("The file has been successfully uploaded!")
}

func (b *Bot) storageHandler(c tele.Context) error {
	userID := c.Sender().ID
	token, ok := b.getSession(userID)
	if !ok {
		return c.Send("You are not logged in. Please login with /login")
	}

	storage, err := api.GetStorage(token)
	if err != nil {
		return c.Send("Error retrieving storage information: " + err.Error())
	}

	reply := fmt.Sprintf(
		"Files in storage (%d):\n%s\n\n Utilized: %d from %d MB",
		len(storage.Files),
		formatFilesList(storage.Files),
		storage.Storage.Used,
		storage.Storage.Limit,
	)

	return c.Send(reply)
}

func (b *Bot) deleteHandler(c tele.Context) error {
	args := c.Args()
	if len(args) < 1 {
		return c.Send("Useage: /delete filename.txt")
	}

	userID := c.Sender().ID
	token, ok := b.getSession(userID)
	if !ok {
		return c.Send("You are not logged in. Please login with /login")
	}

	filename := args[0]

	if err := api.DeleteFile(token, filename); err != nil {
		return c.Send("File deletion error: " + err.Error())
	}

	return c.Send("The file has been successfully deleted: " + filename)
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

func (b *Bot) downloadHandler(c tele.Context) error {
	args := c.Args()
	if len(args) < 1 {
		return c.Send("Usage: /download <filename>")
	}
	filename := args[0]

	userID := c.Sender().ID
	token, ok := b.getSession(userID)
	if !ok {
		return c.Send("You are not logged in. Please login with /login")
	}

	data, name, err := api.DownloadFile(token, filename)
	if err != nil {
		return c.Send("File upload error: " + err.Error())
	}

	doc := &tele.Document{
		File:     tele.File{FileReader: bytes.NewReader(data)},
		FileName: name,
	}

	return c.Send(doc)
}
