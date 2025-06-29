package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"lockbot/internal/config"
	"lockbot/internal/models"
	"net/http"
)

func Login(email, password string) (*models.LoginResponse, error) {
	reqBody := models.LoginRequest{
		Email:    email,
		Password: password,
	}

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
	reqBody := models.RegisterRequest{
		Email:    email,
		Name:     name,
		Password: password,
	}

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
