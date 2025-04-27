package model

import "time"

type TwitchAuth struct {
	Username              string
	AccessToken           string
	AccessTokenValidTill  time.Time
	RefreshToken          string
	RefreshTokenValidTill time.Time
}
