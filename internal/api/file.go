package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"lockbot/internal/config"
	"lockbot/internal/models"
	"mime/multipart"
	"net/url"
)

func UploadFile(token string, filename string, fileReader io.Reader) ([]byte, error) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, fmt.Errorf("error when creating multipart: %w", err)
	}

	if _, err = io.Copy(part, fileReader); err != nil {
		return nil, fmt.Errorf("error when copying a file: %w", err)
	}

	func() { _ = writer.Close() }()

	headers := models.Headers{
		ContentType:   "application/json",
		Accept:        "application/json",
		Authorization: "Brearer " + token,
	}.ToMap()

	url := fmt.Sprintf("%supload", config.Api)
	body, status, err := doRequest("POST", url, &buf, headers)
	if err != nil {
		return nil, err
	}

	statusCheck(status, body)

	return body, nil
}

func GetStorage(token string) (*models.StorageResponse, error) {
	headers := models.Headers{
		ContentType: "application/json",
		Accept:      "application/json",
	}.ToMap()

	url := fmt.Sprintf("%sstrorage", config.Api)
	body, status, err := doRequest("GET", url, nil, headers)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}

	statusCheck(status, body)

	var resp models.StorageResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("response parsing error: %w", err)
	}

	return &resp, nil
}

func DeleteFile(token, filename string) error {
	escapedName := url.PathEscape(filename)

	headers := models.Headers{
		ContentType: "application/json",
		Accept:      "application/json",
	}.ToMap()

	url := fmt.Sprintf("%sdelete/%s", config.Api, escapedName)
	body, status, err := doRequest("DELETE", url, nil, headers)
	if err != nil {
		return err
	}

	statusCheck(status, body)

	return nil
}

func DownloadFile(token, filename string) ([]byte, string, error) {
	headers := models.Headers{
		Authorization: "Brearer " + token,
		Accept:        "application/octet-stream",
	}.ToMap()

	url := fmt.Sprintf("%sstorage/%s", config.Api, filename)
	body, status, err := doRequest("GET", url, nil, headers)
	if err != nil {
		return nil, "", fmt.Errorf("query error: %w", err)
	}

	statusCheck(status, body)

	return body, filename, nil
}
