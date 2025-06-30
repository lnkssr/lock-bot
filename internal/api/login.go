package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"lockbot/internal/config"
	"lockbot/internal/models"
)

func Login(email, password string) (*models.LoginResponse, error) {
	reqBody := models.LoginRequest{Email: email, Password: password}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("JSON encoding error: %w", err)
	}

	headers := models.Headers{
		ContentType: "application/json",
		Accept:      "application/json",
	}.ToMap()

	url := fmt.Sprintf("%slogin", config.Api)
	body, status, err := doRequest("POST", url, bytes.NewBuffer(jsonBody), headers)
	if err != nil {
		return nil, err
	}

	statusCheck(status, body)

	var loginResp models.LoginResponse
	if err := json.Unmarshal(body, &loginResp); err != nil {
		return nil, fmt.Errorf("response parsing error: %w", err)
	}

	return &loginResp, nil
}

func Logout(token string) error {
	headers := models.Headers{
		Authorization: "Brearer " + token,
		Accept:        "application/json",
	}.ToMap()

	url := fmt.Sprintf("%slogout", config.Api)
	body, status, err := doRequest("POST", url, nil, headers)
	if err != nil {
		return err
	}

	statusCheck(status, body)

	return nil
}

func Register(email, name, password string) (*models.RegisterResponse, error) {
	reqBody := models.RegisterRequest{Email: email, Name: name, Password: password}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("JSON encoding error: %w", err)
	}

	headers := models.Headers{
		ContentType: "application/json",
		Accept:      "application/json",
	}.ToMap()

	url := fmt.Sprintf("%sregister", config.Api)
	body, status, err := doRequest("POST", url, bytes.NewBuffer(jsonBody), headers)
	if err != nil {
		return nil, err
	}

	statusCheck(status, body)

	var resp models.RegisterResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("response parsing error: %w", err)
	}

	return &resp, nil
}
