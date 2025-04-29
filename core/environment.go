package core

import (
	"os"

	"github.com/derkellernerd/kellerbot/database"
)

type EnvironmentTwitchSettings struct {
	ClientId     string
	ClientSecret string
	CookieSecret string
	RedirectUrl  string
}

type Environment struct {
	DatabaseManager *database.DatabaseManager
	Twitch          EnvironmentTwitchSettings
}

func NewEnvironment() *Environment {
	return &Environment{
		Twitch: EnvironmentTwitchSettings{
			ClientId:     os.Getenv("TWITCH_CLIENT_ID"),
			ClientSecret: os.Getenv("TWITCH_CLIENT_SECRET"),
			CookieSecret: os.Getenv("TWITCH_COOKIE_SECRET"),
			RedirectUrl:  os.Getenv("TWITCH_REDIRECT_URL"),
		},
	}
}
