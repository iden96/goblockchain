package utils

import "encoding/json"

func JsonStatus(message string) ([]byte, error) {
	return json.Marshal(struct {
		Message string `json:"message"`
	}{
		Message: message,
	})
}
