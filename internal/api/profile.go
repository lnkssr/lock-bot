package api

import (
	"fmt"
	"lockbot/internal/config"
	"lockbot/internal/models"
)

func Profile(token string) ([]byte, error) {
	headers := models.Headers{
		Authorization: "Bearer " + token,
		Accept:        "application/json",
	}.ToMap()

	body, status, err := doRequest(
		"GET",
		fmt.Sprintf("%sprofile", config.Api),
		nil,
		headers)
	if err != nil {
		return nil, err
	}

	statusCheck(status, body)

	return body, nil
}
