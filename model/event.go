package model

import "github.com/goccy/go-json"

type ChatEvent struct {
	User    string `binding:"required"`
	Message string `binding:"required"`
}

func (c *ChatEvent) ToJson() (string, error) {
	bytes, err := json.Marshal(c)
	if err != nil {
		return "", nil
	}

	return string(bytes), nil
}

type AlertEvent struct {
	Name string `binding:"required"`
}

func (a *AlertEvent) ToJson() (string, error) {
	bytes, err := json.Marshal(a)
	if err != nil {
		return "", nil
	}

	return string(bytes), nil
}
