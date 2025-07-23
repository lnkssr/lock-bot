package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"lockbot/internal/config"
	logger "lockbot/internal/log"
	"lockbot/internal/models"
)

func GetAllUsers(token string) ([]models.User, error) {
	headers := models.Headers{
		Authorization: "Bearer " + token,
		Accept:        "application/json",
	}.ToMap()

	logger.Debug(headers)

	body, err := doRequest(
		"GET",
		fmt.Sprintf("%sadmin/users", config.Api),
		nil,
		headers)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}

	var users []models.User
	if err := json.Unmarshal(body, &users); err != nil {
		return nil, fmt.Errorf("response parsing error: %w", err)
	}

	logger.Debug(users, body)

	return users, nil
}

func MakeAdmin(token string, userID int) error {
	headers := models.Headers{
		Authorization: "Bearer " + token,
		Accept:        "application/json",
	}.ToMap()

	_, err := doRequest(
		"PUT",
		fmt.Sprintf("%sadmin/make_admin/%d", config.Api, userID),
		nil,
		headers)
	if err != nil {
		return fmt.Errorf("query error: %w", err)
	}

	return nil
}

func RevokeAdmin(token string, userID int) error {
	headers := models.Headers{
		Authorization: "Bearer " + token,
		Accept:        "application/json",
	}.ToMap()

	_, err := doRequest(
		"PUT",
		fmt.Sprintf("%sadmin/revoke_admin/%d", config.Api, userID),
		nil,
		headers)
	if err != nil {
		return fmt.Errorf("query error: %w", err)
	}

	return nil
}

func UpdateUserLimit(token string, userID, newLimit int) error {
	reqBody := models.LimitRequest{User: userID, Limit: newLimit}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("JSON encoding error: %w", err)
	}

	headers := models.Headers{
		ContentType:   "application/json",
		Accept:        "application/json",
		Authorization: "Bearer " + token,
	}.ToMap()

	_, err = doRequest(
		"PUT",
		fmt.Sprintf("%sadmin/update_limit", config.Api),
		bytes.NewBuffer(jsonBody),
		headers)
	if err != nil {
		return fmt.Errorf("query error: %w", err)
	}

	return nil
}
