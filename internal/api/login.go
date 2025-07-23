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

	body, err := doRequest(
		"POST",
		fmt.Sprintf("%slogin", config.Api),
		bytes.NewBuffer(jsonBody),
		headers)
	if err != nil {
		return nil, err
	}

	var loginResp models.LoginResponse
	if err := json.Unmarshal(body, &loginResp); err != nil {
		return nil, fmt.Errorf("response parsing error: %w", err)
	}

	return &loginResp, nil
}

func Logout(token string) error {
	headers := models.Headers{
		Authorization: "Bearer " + token,
		Accept:        "application/json",
	}.ToMap()

	_, err := doRequest("POST", fmt.Sprintf("%slogout", config.Api), nil, headers)
	if err != nil {
		return err
	}

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

	body, err := doRequest(
		"POST",
		fmt.Sprintf("%sregister", config.Api),
		bytes.NewBuffer(jsonBody),
		headers)
	if err != nil {
		return nil, err
	}

	var resp models.RegisterResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("response parsing error: %w", err)
	}

	return &resp, nil
}
