package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"lockbot/internal/config"
	"lockbot/internal/models"
)

func GetAllUsers(token string) ([]models.User, error) {
	headers := models.Headers{
		Authorization: "Brearer " + token,
		Accept:        "application/json",
	}.ToMap()

	url := fmt.Sprintf("%sadmin/users", config.Api)
	body, status, err := doRequest("GET", url, nil, headers)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}

	statusCheck(status, body)

	var users []models.User
	if err := json.Unmarshal(body, &users); err != nil {
		return nil, fmt.Errorf("response parsing error: %w", err)
	}

	return users, nil
}

func MakeAdmin(token string, userID int) error {
	headers := models.Headers{
		Authorization: "Brearer " + token,
		Accept:        "application/json",
	}.ToMap()

	url := fmt.Sprintf("%sadmin/make_admin/%d", config.Api, userID)
	body, status, err := doRequest("PUT", url, nil, headers)
	if err != nil {
		return fmt.Errorf("query error: %w", err)
	}

	statusCheck(status, body)

	return nil
}

func RevokeAdmin(token string, userID int) error {
	headers := models.Headers{
		Authorization: "Brearer " + token,
		Accept:        "application/json",
	}.ToMap()

	url := fmt.Sprintf("%sadmin/revoke_admin/%d", config.Api, userID)
	body, status, err := doRequest("PUT", url, nil, headers)
	if err != nil {
		return fmt.Errorf("query error: %w", err)
	}

	statusCheck(status, body)

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
		Authorization: "Brearer " + token,
	}.ToMap()

	url := fmt.Sprintf("%sadmin/update_limit", config.Api)
	body, status, err := doRequest("PUT", url, bytes.NewBuffer(jsonBody), headers)
	if err != nil {
		return fmt.Errorf("query error: %w", err)
	}
	statusCheck(status, body)

	return nil
}
