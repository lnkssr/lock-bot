package api

import (
	"fmt"
	"io"
	"net/http"
)

func statusCheck(status int, body []byte) error {
	if status < 200 || status >= 300 {
		return fmt.Errorf("server error: %s", string(body))
	}
	return nil
}

func doRequest(method, url string, body io.Reader, headers map[string]string) ([]byte, int, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, 0, fmt.Errorf("query creation error: %w", err)
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("request sending error: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("response reading error: %w", err)
	}

	return respBody, resp.StatusCode, nil
}
