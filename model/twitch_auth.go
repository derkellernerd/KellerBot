package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
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
		return fmt.Errorf("failed to unmarshal JSONB value: %v", value)
	}

	var result TwitchAuth
	err := json.Unmarshal(bytes, &result)
	*t = result
	return err
}

func (t TwitchAuth) Value() (driver.Value, error) {
	if len(t.AccessToken) == 0 {
		return nil, nil
	}
	return json.Marshal(&t)
}

func (t TwitchAuth) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case "mysql", "sqlite":
		return "JSON"
	case "postgres":
		return "JSONB"
	}
	return ""
}
