package model

import (
	"github.com/goccy/go-json"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type TwitchEventLog struct {
	gorm.Model
	TwitchEventSubscription string
	Data                    datatypes.JSON
}

func (t *TwitchEventLog) SetData(data any) error {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	t.Data = datatypes.JSON(jsonBytes)
	return nil
}

type TwitchEvent struct {
	gorm.Model
	TwitchEventSubscription string `gorm:"unique"`
	ActionName              string
}

type TwitchEventCreateRequest struct {
	TwitchEventSubscription string `binding:"required"`
	ActionName              string `binding:"required"`
}

type TwitchEventUpdateRequest struct {
	ActionName string `binding:"required"`
}
