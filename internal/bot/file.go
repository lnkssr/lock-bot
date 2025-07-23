package bot

import (
	"bytes"
	"fmt"
	"lockbot/internal/api"
	logger "lockbot/internal/log"
	"strings"

	tele "gopkg.in/telebot.v4"
)

func (b *Bot) uploadHandler(c tele.Context) error {
	file := c.Message().Document
	userID := c.Sender().ID

	if file == nil {
		logger.Info("No document attached", userID)
		return c.Send("Please send the file as a document.")
	}

	reader, err := b.api.File(&file.File)
	if err != nil {
		logger.Error("Failed to retrieve file from Telegram", userID, err)
		return c.Send("Error while retrieving the file: " + err.Error())
	}
	defer func() { _ = reader.Close() }()

	token, ok := b.getSession(userID)
	if !ok {
		logger.Warn("Upload attempt without login", userID)
		return c.Send("You are not logged in. Please login with /login")
	}

	logger.Debug("Uploading file", userID, file.FileName)

	go func() {
		_, err = api.UploadFile(token, file.FileName, reader)
		if err != nil {
			logger.Error("File upload failed", userID, file.FileName, err)
			_ = c.Send("File upload error: " + err.Error())
			return
		}

		logger.Info("File uploaded successfully", userID, file.FileName)
	}()
	return c.Send("The file has been successfully uploaded!")
}

func (b *Bot) storageHandler(c tele.Context) error {
	userID := c.Sender().ID

	token, ok := b.getSession(userID)
	if !ok {
		logger.Warn("Storage access attempt without login", userID)
		return c.Send("You are not logged in. Please login with /login")
	}

	storage, err := api.GetStorage(token)
	if err != nil {
		logger.Error("Failed to get storage info", userID, err)
		return c.Send("Error retrieving storage information: " + err.Error())
	}

	reply := fmt.Sprintf(
		"Files in storage (%d):\n%s\n\nUtilized: %d of %d MB",
		len(storage.Files),
		formatFilesList(storage.Files),
		storage.Storage.Used,
		storage.Storage.Limit,
	)

	logger.Info("Storage info retrieved successfully", userID)
	return c.Send(reply)
}

func (b *Bot) deleteHandler(c tele.Context) error {
	userID := c.Sender().ID
	args := c.Args()

	if len(args) < 1 {
		logger.Warn("Filename not provided for deletion", userID)
		return c.Send("Usage: /delete <filename>")
	}

	token, ok := b.getSession(userID)
	if !ok {
		logger.Warn("Delete attempt without login", userID)
		return c.Send("You are not logged in. Please login with /login")
	}

	filename := strings.Join(args, " ")

	go func() {
		if err := api.DeleteFile(token, filename); err != nil {
			logger.Error("Failed to delete file", userID, filename, err)
			_ = c.Send("File deletion error: " + err.Error())
			return
		}

		logger.Info("File deleted", userID, filename)
	}()
	return c.Send("The file has been successfully deleted: " + filename)
}

func (b *Bot) downloadHandler(c tele.Context) error {
	userID := c.Sender().ID
	args := c.Args()

	if len(args) < 1 {
		logger.Warn("Filename not provided for download", userID)
		return c.Send("Usage: /download <filename>")
	}

	filename := strings.Join(args, " ")

	token, ok := b.getSession(userID)
	if !ok {
		logger.Warn("Download attempt without login", userID)
		return c.Send("You are not logged in. Please login with /login")
	}

	data, name, err := api.DownloadFile(token, filename)
	if err != nil {
		logger.Error("File download failed", userID, filename, err)
		return c.Send("File download error: " + err.Error())
	}

	doc := &tele.Document{
		File:     tele.File{FileReader: bytes.NewReader(data)},
		FileName: name,
	}

	logger.Info("File downloaded successfully", userID, filename)
	return c.Send(doc)
}
