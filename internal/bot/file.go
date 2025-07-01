package bot

import (
	"bytes"
	"fmt"
	"lockbot/internal/api"
	"strings"

	tele "gopkg.in/telebot.v4"
)

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
	fmt.Println(storage, userID, token, reply)

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

	filename := strings.Join(args, " ")

	if err := api.DeleteFile(token, filename); err != nil {
		return c.Send("File deletion error: " + err.Error())
	}

	return c.Send("The file has been successfully deleted: " + filename)
}

func (b *Bot) downloadHandler(c tele.Context) error {
	args := c.Args()
	if len(args) < 1 {
		return c.Send("Usage: /download <filename>")
	}
	filename := strings.Join(args, " ")

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
