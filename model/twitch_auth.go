package model

import "time"

type TwitchAuth struct {
	AccessToken           string
	AccessTokenValidTill  time.Time
	RefreshToken          string
	RefreshTokenValidTill time.Time
}
