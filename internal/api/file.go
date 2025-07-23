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
		ContentType:   writer.FormDataContentType(),
		Accept:        "application/json",
		Authorization: "Bearer " + token,
	}.ToMap()

	body, err := doRequest(
		"POST",
		fmt.Sprintf("%supload", config.Api),
		&buf,
		headers)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func GetStorage(token string) (*models.StorageResponse, error) {
	headers := models.Headers{
		ContentType:   "application/json",
		Accept:        "application/json",
		Authorization: "Bearer " + token,
	}.ToMap()

	body, err := doRequest(
		"GET",
		fmt.Sprintf("%sstorage", config.Api),
		nil,
		headers)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}

	var resp models.StorageResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("response parsing error: %w", err)
	}

	return &resp, nil
}

func DeleteFile(token, filename string) error {
	escapedName := url.PathEscape(filename)

	headers := models.Headers{
		ContentType:   "application/json",
		Accept:        "application/json",
		Authorization: "Bearer " + token,
	}.ToMap()

	_, err := doRequest(
		"DELETE",
		fmt.Sprintf("%sdelete/%s", config.Api, escapedName),
		nil,
		headers)
	if err != nil {
		return err
	}

	return nil
}

func DownloadFile(token, filename string) ([]byte, string, error) {
	headers := models.Headers{
		Authorization: "Bearer " + token,
		Accept:        "application/octet-stream",
	}.ToMap()

	body, err := doRequest(
		"GET",
		fmt.Sprintf("%sstorage/%s", config.Api, filename),
		nil,
		headers)
	if err != nil {
		return nil, "", fmt.Errorf("query error: %w", err)
	}

	return body, filename, nil
}
