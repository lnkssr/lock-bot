package api

import "fmt"

func statusCheck(status int, body []byte) error {
	if status < 200 || status >= 300 {
		return fmt.Errorf("server error: %s", string(body))
	}
	return nil
}
