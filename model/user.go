package model

import (
	"encoding/json"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Identifier     string
	Username       string
	TwitchUserAuth TwitchAuth
	TwitchBotAuth  TwitchAuth
}

func (u *User) SetTwitchUserAuth(data any) error {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	u.Payload = datatypes.JSON(jsonBytes)
	return nil
}

func (u *User) GetTwitchUserAuth() (map[string]any, error) {
	var payload map[string]any
	jsonBytes, err := e.Payload.MarshalJSON()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(jsonBytes, &payload)
	return payload, err
}
