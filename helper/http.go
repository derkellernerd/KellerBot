package helper

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func BasicHttpRequest(method string, url string, payload any) error {
	var reader io.Reader
	if payload != nil {
		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			return err
		}

		reader = bytes.NewReader(payloadBytes)
	}

	log.Printf("HTTP: Call: %s %s", method, url)
	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
