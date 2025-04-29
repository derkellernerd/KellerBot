package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type twitchAuth struct {
	Username              string
	AccessToken           string
	AccessTokenValidTill  time.Time
	RefreshToken          string
	RefreshTokenValidTill time.Time
}

type TwitchAuth twitchAuth

func (t *TwitchAuth) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	result := &twitchAuth{}
	err := json.Unmarshal(bytes, &result)
	*t = TwitchAuth(*result)
	return err
}

func (t TwitchAuth) Value() (driver.Value, error) {
	if len(t.Username) == 0 {
		return nil, nil
	}

	bytes, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}
