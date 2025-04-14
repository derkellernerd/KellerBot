package model

import (
	"github.com/goccy/go-json"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

const (
	ACTION_TYPE_CHAT_MESSAGE ActionType = "CHAT_MESSAGE"
	ACTION_TYPE_CHAT_ANSWER  ActionType = "CHAT_ANSWER"
	ACTION_TYPE_GIF          ActionType = "GIF"
	ACTION_TYPE_SOUND        ActionType = "SOUND"
	ACTION_TYPE_HTTP         ActionType = "HTTP"
	ACTION_TYPE_TEXT         ActionType = "TEXT"
	ACTION_TYPE_COMPOSITION  ActionType = "COMPOSITION"
)

type ActionType string

type ActionTypeHttp struct {
	Uri        string
	HttpMethod string
	Payload    any
}

type ActionTypeChatMessage struct {
	ChatMessage string
}

type ActionTypeChatAnswer struct {
	ChatMessage string
}

type ActionTypeComposition struct {
	Actions []string
}

type ActionTypeSound struct {
	FileName string
	Gain     float64
}

type ActionTypeGif struct {
	FileName string
}

type ActionTypeText struct {
	Text string
}

type ActionTypes interface {
	ActionTypeChatMessage | ActionTypeChatAnswer | ActionTypeGif | ActionTypeSound | ActionTypeHttp
}

type Action struct {
	gorm.Model
	ActionName string `gorm:"uniqueIndex"`
	ActionType ActionType
	Data       datatypes.JSON
}

func (a *Action) SetData(data any) error {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	a.Data = datatypes.JSON(jsonBytes)
	return nil
}

func ActionGetData[V ActionTypes](action *Action) (V, error) {
	var actionData V
	jsonBytes, err := action.Data.MarshalJSON()
	if err != nil {
		return actionData, err
	}
	err = json.Unmarshal(jsonBytes, &actionData)
	return actionData, err
}

func NewAction(actionName string, actionType ActionType, data any) (Action, error) {
	action := Action{
		ActionName: actionName,
		ActionType: actionType,
	}
	err := action.SetData(data)
	return action, err
}

type ActionCreateRequest struct {
	ActionName string     `binding:"required"`
	ActionType ActionType `binding:"required"`
	Data       any        `binding:"required"`
}

type ActionUpdateRequest struct {
	Data any `binding:"required"`
}
