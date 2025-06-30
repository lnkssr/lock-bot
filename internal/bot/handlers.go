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
	// TODO: help list command
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
		return c.Send("Использование: /logn email password")
	}

	email := args[0]
	password := strings.Join(args[1:], " ")

	resp, err := api.Login(email, password)
	if err != nil {
		return c.Send("Ошибка: " + err.Error())
	}

	b.saveSession(c.Sender().ID, resp.Token, 24*time.Hour)

	reply := fmt.Sprintf(
		"%s\nПользователь: %s (ID: %d)\nТокен: %s",
		resp.Message,
		resp.User.Email,
		resp.User.ID,
		resp.Token,
	)

	return c.Send(reply)
}

func (b *Bot) profileHandler(c tele.Context) error {
	userID := c.Sender().ID

	token, ok := b.getSession(userID)
	if !ok {
		return c.Send("Вы не авторизованы. Введите /logn email password")
	}

	profileData, err := api.Profile(token)
	if err != nil {
		return c.Send("Ошибка получения профиля: " + err.Error())
	}

	return c.Send("Профиль:\n" + string(profileData))
}

func (b *Bot) logoutHandler(c tele.Context) error {
	userID := c.Sender().ID

	token, ok := b.getSession(userID)
	if !ok {
		return c.Send("Вы не авторизованы.")
	}

	err := api.Logout(token)
	if err != nil {
		return c.Send("Ошибка выхода: " + err.Error())
	}

	delete(b.sessions, userID)

	return c.Send("Вы успешно вышли из аккаунта.")
}

func (b *Bot) registerHandler(c tele.Context) error {
	args := c.Args()
	if len(args) < 3 {
		return c.Send("Использование: /register email name password")
	}

	email := args[0]
	name := args[1]
	password := strings.Join(args[2:], " ")

	resp, err := api.Register(email, name, password)
	if err != nil {
		return c.Send("Ошибка регистрации: " + err.Error())
	}

	reply := fmt.Sprintf(
		"%s\nПользователь: %s (ID: %d)",
		resp.Message,
		resp.User.Email,
		resp.User.ID,
	)

	return c.Send(reply)
}

func (b *Bot) uploadHandler(c tele.Context) error {
	file := c.Message().Document
	if file == nil {
		return c.Send("Пожалуйста, отправьте файл в виде документа.")
	}

	reader, err := b.api.File(&file.File)
	if err != nil {
		return c.Send("Ошибка при получении файла: " + err.Error())
	}
	defer reader.Close()

	userID := c.Sender().ID
	token, ok := b.getSession(userID)
	if !ok {
		return c.Send("Вы не авторизованы. Пожалуйста, войдите с помощью /login")
	}

	respBody, err := api.UploadFile(token, file.FileName, reader)
	if err != nil {
		return c.Send("Ошибка загрузки файла: " + err.Error())
	}

	return c.Send("Файл успешно загружен!\nОтвет: " + string(respBody))
}

func (b *Bot) storageHandler(c tele.Context) error {
	userID := c.Sender().ID
	token, ok := b.getSession(userID)
	if !ok {
		return c.Send("Вы не авторизованы. Пожалуйста, войдите с помощью /login")
	}

	storage, err := api.GetStorage(token)
	if err != nil {
		return c.Send("Ошибка получения информации о хранилище: " + err.Error())
	}

	reply := fmt.Sprintf(
		"Файлы в хранилище (%d):\n%s\n\nИспользовано: %d из %d байт",
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
		return c.Send("Использование: /delete filename.txt")
	}

	userID := c.Sender().ID
	token, ok := b.getSession(userID)
	if !ok {
		return c.Send("Вы не авторизованы. Пожалуйста, выполните вход с помощью /login")
	}

	filename := args[0]

	if err := api.DeleteFile(token, filename); err != nil {
		return c.Send("Ошибка удаления файла: " + err.Error())
	}

	return c.Send("Файл успешно удалён: " + filename)
}

func (b *Bot) usersHandler(c tele.Context) error {
	userID := c.Sender().ID
	token, ok := b.getSession(userID)
	if !ok {
		return c.Send("Вы не авторизованы. Используйте /login")
	}

	users, err := api.GetAllUsers(token)
	if err != nil {
		return c.Send("Ошибка получения списка пользователей: " + err.Error())
	}

	if len(users) == 0 {
		return c.Send("Пользователи не найдены.")
	}

	var sb strings.Builder
	sb.WriteString("Пользователи:\n")
	for _, u := range users {
		sb.WriteString(fmt.Sprintf("• %s (ID: %d)\n", u.Email, u.ID))
	}

	return c.Send(sb.String())
}

func (b *Bot) downloadHandler(c tele.Context) error {
	args := c.Args()
	if len(args) < 1 {
		return c.Send("Использование: /download <имя_файла>")
	}
	filename := args[0]

	userID := c.Sender().ID
	token, ok := b.getSession(userID)
	if !ok {
		return c.Send("Вы не авторизованы. Используйте /login")
	}

	data, name, err := api.DownloadFile(token, filename)
	if err != nil {
		return c.Send("Ошибка загрузки файла: " + err.Error())
	}

	doc := &tele.Document{
		File:     tele.File{FileReader: bytes.NewReader(data)},
		FileName: name,
	}

	return c.Send(doc)
}
