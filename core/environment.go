package core

import (
	"os"

	"github.com/derkellernerd/kellerbot/database"
)

type Environment struct {
	DatabaseManager *database.DatabaseManager
	Twitch          struct {
		UserId       string
		ClientId     string
		ClientSecret string
	}
	TwitchSession *TwitchSession
}

func NewEnvironment(twitchSession *TwitchSession) *Environment {
	return &Environment{
		Twitch: struct {
			UserId       string
			ClientId     string
			ClientSecret string
		}{
			UserId:       os.Getenv("TWITCH_USER_ID"),
			ClientId:     os.Getenv("TWITCH_CLIENT_ID"),
			ClientSecret: os.Getenv("TWITCH_CLIENT_SECRET"),
		},
		TwitchSession: twitchSession,
	}
}
