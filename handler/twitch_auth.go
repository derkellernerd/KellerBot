package handler

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/derkellernerd/kellerbot/core"
	"github.com/derkellernerd/kellerbot/repository"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/twitch"
)

const (
	stateCallbackKey = "oauth-state-callback"
	oauthSessionName = "oauth-session"
	oauthTokenKey    = "oauth-token"
)

var (
	twitchScopes = []string{
		"user:read:chat",
		"user:bot",
		"user:write:chat",
		"moderator:read:followers",
	}
)

type TwitchAuth struct {
	env  *core.Environment
	user *repository.User
}

func NewTwitchAuth(env *core.Environment, userRepo *repository.User) *TwitchAuth {
	return &TwitchAuth{
		env: env,
	}
}

func (h *TwitchAuth) getTwitchOAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     h.env.Twitch.ClientId,
		ClientSecret: h.env.Twitch.ClientSecret,
		Scopes:       twitchScopes,
		Endpoint:     twitch.Endpoint,
		RedirectURL:  h.env.Twitch.RedirectUrl,
	}
}

func (h *TwitchAuth) Login(c *gin.Context) {
	_, err := c.Cookie(oauthSessionName)
	if err != nil {
		log.Printf("corrupted session %s -- generated new", err)
		err = nil
	}

	var tokenBytes [255]byte
	if _, err := rand.Read(tokenBytes[:]); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	state := hex.EncodeToString(tokenBytes[:])

	fmt.Printf("STATE: %s\n", state)
	c.SetCookie(stateCallbackKey, state, 3600, "/", "localhost", false, true)
	c.Redirect(http.StatusTemporaryRedirect, h.getTwitchOAuthConfig().AuthCodeURL(state))
}

func (h *TwitchAuth) Callback(c *gin.Context) {
	_, err := c.Cookie(oauthSessionName)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, NewErrorResponse(err))
		return
	}

	stateChallenge, _ := c.Cookie(stateCallbackKey)
	state := c.Query("state")
	switch stateChallenge, state := stateChallenge, state; {
	case state == "", len(stateChallenge) < 1:
		err = errors.New("missing state challenge")
	case state != stateChallenge:
		err = fmt.Errorf("invalid oauth state, expected '%s', got '%s'\n", state, stateChallenge)
	}

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, NewErrorResponse(err))
		return
	}

	token, err := h.getTwitchOAuthConfig().Exchange(context.Background(), c.Query("code"))
	if err != nil {
		return
	}

	fmt.Printf("%#v\n", token)
	c.Redirect(http.StatusTemporaryRedirect, "/ui")

	return
}
