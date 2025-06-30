package api

import (
	"fmt"
	"lockbot/internal/config"
	"lockbot/internal/models"
)

func Profile(token string) ([]byte, error) {
	headers := models.Headers{
		Authorization: "Brearer " + token,
		Accept:        "application/json",
	}.ToMap()

	url := fmt.Sprintf("%sprofile", config.Api)
	body, status, err := doRequest("GET", url, nil, headers)
	if err != nil {
		return nil, err
	}

	statusCheck(status, body)

	return body, nil
}
