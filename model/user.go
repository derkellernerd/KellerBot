package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Identifier     string
	Username       string
	TwitchUserAuth TwitchAuth
	TwitchBotAuth  TwitchAuth
}
