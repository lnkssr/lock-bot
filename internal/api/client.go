package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"lockbot/internal/config"
	"lockbot/internal/models"
	"mime/multipart"
	"net/http"
)

func Login(email, password string) (*models.LoginResponse, error) {
	reqBody := models.LoginRequest{Email: email, Password: password}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("ошибка кодирования JSON: %w", err)
	}

	headers := map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}

	body, status, err := doRequest("POST", config.Api+"login", bytes.NewBuffer(jsonBody), headers)
	if err != nil {
		return nil, err
	}

	if status != http.StatusOK {
		return nil, fmt.Errorf("ошибка сервера: %s", string(body))
	}

	var loginResp models.LoginResponse
	if err := json.Unmarshal(body, &loginResp); err != nil {
		return nil, fmt.Errorf("ошибка парсинга ответа: %w", err)
	}

	return &loginResp, nil
}

func Profile(token string) ([]byte, error) {
	headers := map[string]string{
		"Authorization": "Bearer " + token,
		"Accept":        "application/json",
	}

	body, status, err := doRequest("GET", config.Api+"profile", nil, headers)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("ошибка сервера: %d %s", status, string(body))
	}

	return body, nil
}

func Logout(token string) error {
	headers := map[string]string{
		"Authorization": "Bearer " + token,
		"Accept":        "application/json",
	}

	body, status, err := doRequest("POST", config.Api+"logout", nil, headers)
	if err != nil {
		return err
	}
	if status != http.StatusOK {
		return fmt.Errorf("ошибка сервера: %d %s", status, string(body))
	}

	return nil
}

func Register(email, name, password string) (*models.RegisterResponse, error) {
	reqBody := models.RegisterRequest{Email: email, Name: name, Password: password}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("ошибка кодирования JSON: %w", err)
	}

	headers := map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}

	body, status, err := doRequest("POST", config.Api+"register", bytes.NewBuffer(jsonBody), headers)
	if err != nil {
		return nil, err
	}

	if status < 200 || status >= 300 {
		return nil, fmt.Errorf("ошибка сервера: %s", string(body))
	}

	var resp models.RegisterResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("ошибка парсинга ответа: %w", err)
	}

	return &resp, nil
}

func UploadFile(token string, filename string, fileReader io.Reader) ([]byte, error) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании multipart: %w", err)
	}

	if _, err = io.Copy(part, fileReader); err != nil {
		return nil, fmt.Errorf("ошибка при копировании файла: %w", err)
	}

	writer.Close()

	headers := map[string]string{
		"Content-Type":  writer.FormDataContentType(),
		"Accept":        "application/json",
		"Authorization": "Bearer " + token,
	}

	body, status, err := doRequest("POST", config.Api+"upload", &buf, headers)
	if err != nil {
		return nil, err
	}

	if status < 200 || status >= 300 {
		return nil, fmt.Errorf("ошибка сервера: %s", string(body))
	}

	return body, nil
}

func GetStorage(token string) (*models.StorageResponse, error) {
	headers := map[string]string{
		"Authorization": "Bearer " + token,
		"Accept":        "application/json",
	}

	body, status, err := doRequest("GET", config.Api+"storage", nil, headers)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса: %w", err)
	}
	if status < 200 || status >= 300 {
		return nil, fmt.Errorf("ошибка сервера: %d %s", status, string(body))
	}

	var resp models.StorageResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("ошибка парсинга ответа: %w", err)
	}

	return &resp, nil
}
