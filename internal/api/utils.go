package api

import (
	"fmt"
	"io"
	"net/http"
)

func statusCheck(status int, respBody []byte) error {
	if status < 200 || status >= 300 {
		return fmt.Errorf("server error: %s", respBody)
	}
	return nil
}

func doRequest(method, url string, body io.Reader, headers map[string]string) ([]byte, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("query creation error: %w", err)
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request sending error: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("response reading error: %w", err)
	}
	err = statusCheck(resp.StatusCode, respBody)
	if err != nil {
		return nil, fmt.Errorf("status code is not avalible: %w", err)
	}
	return respBody, nil
}
