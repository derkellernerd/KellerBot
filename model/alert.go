package model

import (
	"github.com/goccy/go-json"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

const (
	ALERT_TYPE_SOUND       AlertType = "SOUND"
	ALERT_TYPE_VIDEO       AlertType = "VIDEO"
	ALERT_TYPE_GIF_SOUND   AlertType = "GIF_SOUND"
	ALERT_TYPE_GIF         AlertType = "GIF"
	ALERT_TYPE_TEXT        AlertType = "TEXT"
	ALERT_TYPE_COMPOSITION AlertType = "COMPOSITION"
	ALERT_TYPE_CHAT        AlertType = "CHAT"
)

type AlertType string

type Alert struct {
	gorm.Model
	Name              string `gorm:"unique"`
	Type              AlertType
	Data              datatypes.JSON
	Used              uint64
	DurationInSeconds float64
}

func (a *Alert) Increment() {
	a.Used += 1
}

func (a *Alert) SetData(data any) error {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	a.Data = datatypes.JSON(jsonBytes)
	return nil
}

func (c *Alert) GetDataComposition() (*AlertTypeComposition, error) {
	var sound AlertTypeComposition
	jsonBytes, err := c.Data.MarshalJSON()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(jsonBytes, &sound)
	return &sound, err
}

func (c *Alert) GetDataText() (*AlertTypeText, error) {
	var sound AlertTypeText
	jsonBytes, err := c.Data.MarshalJSON()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(jsonBytes, &sound)
	return &sound, err
}

func (c *Alert) GetDataSound() (*AlertTypeSound, error) {
	var sound AlertTypeSound
	jsonBytes, err := c.Data.MarshalJSON()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(jsonBytes, &sound)
	return &sound, err
}

func (c *Alert) GetDataVideo() (*AlertTypeVideo, error) {
	var video AlertTypeVideo
	jsonBytes, err := c.Data.MarshalJSON()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(jsonBytes, &video)
	return &video, err
}

func (c *Alert) GetDataGifSound() (*AlertTypeGifSound, error) {
	var video AlertTypeGifSound
	jsonBytes, err := c.Data.MarshalJSON()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(jsonBytes, &video)
	return &video, err
}

func (c *Alert) GetDataGif() (*AlertTypeGif, error) {
	var gif AlertTypeGif
	jsonBytes, err := c.Data.MarshalJSON()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(jsonBytes, &gif)
	return &gif, err
}

func (c *Alert) GetDataChatText() (*AlertTypeChat, error) {
	var gif AlertTypeChat
	jsonBytes, err := c.Data.MarshalJSON()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(jsonBytes, &gif)
	return &gif, err
}

type AlertTypeSound struct {
	SoundPath string
}

type AlertTypeVideo struct {
	VideoPath string
}

type AlertTypeGifSound struct {
	GifPath   string
	SoundPath string
}

type AlertTypeGif struct {
	GifPath string
}

type AlertTypeComposition struct {
	AlertNames []string
}

type AlertTypeText struct {
	Text string
}

type AlertTypeChat struct {
	Chat string
}

type AlertCreateRequest struct {
	Name string    `binding:"required"`
	Type AlertType `binding:"required"`
	Data any
}
