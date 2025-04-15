package model

import (
	"encoding/json"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

const (
	EVENT_SOURCE_TWITCH EventSource = "TWITCH"
)

type EventSource string

type Event struct {
	gorm.Model
	EventName           string
	Source              EventSource
	Payload             datatypes.JSON
	ExecutingActionName string
}

func (e *Event) SetPayload(data any) error {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	e.Payload = datatypes.JSON(jsonBytes)
	return nil
}

func (e *Event) GetPayload() (map[string]any, error) {
	var payload map[string]any
	jsonBytes, err := e.Payload.MarshalJSON()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(jsonBytes, &payload)
	return payload, err
}
