package model

import (
	"github.com/goccy/go-json"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

const (
	COMMAND_TYPE_MESSAGE     CommandType = "MESSAGE_ACTION"
	COMMAND_TYPE_HTTP_ACTION CommandType = "HTTP_ACTION"
)

type CommandType string

type Command struct {
	gorm.Model
	Command string `gorm:"unique"`
	Type    CommandType
	Data    datatypes.JSON
	Used    uint64
}

func NewCommand(command string, commandType CommandType, data any) Command {
	cmd := Command{
		Command: command,
		Type:    commandType,
	}

	cmd.SetData(data)
	return cmd
}

func NewCommandMessage(command string, message string) Command {
	data := CommandActionMessage{
		Message: message,
	}

	return NewCommand(command, COMMAND_TYPE_MESSAGE, data)
}

func (c *Command) Increment() {
	c.Used += 1
}

func (c *Command) SetData(data any) error {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	c.Data = datatypes.JSON(jsonBytes)
	return nil
}

func (c *Command) GetDataActionMessage() (*CommandActionMessage, error) {
	var actionMessage CommandActionMessage
	jsonBytes, err := c.Data.MarshalJSON()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(jsonBytes, &actionMessage)
	return &actionMessage, err
}

type CommandCreateRequest struct {
	Command string      `binding:"required"`
	Type    CommandType `binding:"required"`
	Data    any         `binding:"required"`
}

type CommandUpdateRequest struct {
	Data any `binding:"required"`
}

type CommandActionMessage struct {
	Message string
}

type CommandActionHttp struct {
	Url     string
	Method  string
	Payload any
}
