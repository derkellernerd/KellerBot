package core

import "time"

type TwitchSession struct {
	AccessToken          string
	AccessTokenValidTill time.Time
	RefreshToken         string
}

func (ts *TwitchSession) IsAuthenticated() bool {
	if ts.AccessToken == "" {
		return false
	}

	return true
}

func (ts *TwitchSession) NeedsRefresh() bool {
	if ts.AccessTokenValidTill.Before(time.Now()) {
		return true
	}

	return false
}
