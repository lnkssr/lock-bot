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

	body, err := doRequest(
		"GET",
		fmt.Sprintf("%sprofile", config.Api),
		nil,
		headers)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func ProfileV2(token string) ([]byte, error) {
	headers := models.Headers{
		Authorization: "Bearer " + token,
		Accept:        "applications/json",
	}.ToMap()

	body, err := doRequest(
		"GET",
		fmt.Sprint(config.Api+"v2/"+"profile"),
		nil,
		headers)

	if err != nil {
		return nil, err
	}

	return body, nil
}
