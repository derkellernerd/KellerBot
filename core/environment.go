package core

import "os"

type Environment struct {
	Twitch struct {
		UserId       string
		AccessToken  string
		ClientId     string
		ClientSecret string
	}
}

func NewEnvironment() *Environment {
	return &Environment{
		Twitch: struct {
			UserId       string
			AccessToken  string
			ClientId     string
			ClientSecret string
		}{
			UserId:       os.Getenv("TWITCH_USER_ID"),
			ClientId:     os.Getenv("TWITCH_CLIENT_ID"),
			ClientSecret: os.Getenv("TWITCH_CLIENT_SECRET"),
		},
	}
}
