package auth

import (
	"context"
	"crypto/rand"
	"encoding/gob"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/derkellernerd/kellerbot/core"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/twitch"
)

type TwitchAuth struct {
	env      *core.Environment
	callback func()
}

var (
	scopes       = []string{"user:read:chat", "user:bot", "user:write:chat", "moderator:read:followers"}
	redirectURL  = "http://localhost:8080/redirect"
	oauth2Config *oauth2.Config
	cookieSecret = []byte("verysensible")
	cookieStore  = sessions.NewCookieStore(cookieSecret)
)

// HandleRoot is a Handler that shows a login button. In production, if the frontend is served / generated
// by Go, it should use html/template to prevent XSS attacks.
func (a *TwitchAuth) HandleRoot(c *gin.Context) {

	c.Header("Content-Type", "text/html; charset=utf-8")

	c.String(http.StatusOK, `<html><body><a href="/login">Login using Twitch</a></body></html>`)

	return
}

// HandleLogin is a Handler that redirects the user to Twitch for login, and provides the 'state'
// parameter which protects against login CSRF.
func (a *TwitchAuth) HandleLogin(c *gin.Context) {
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
	c.Redirect(http.StatusTemporaryRedirect, oauth2Config.AuthCodeURL(state))

	return
}

// HandleOauth2Callback is a Handler for oauth's 'redirect_uri' endpoint;
// it validates the state token and retrieves an OAuth token from the request parameters.
func (a *TwitchAuth) HandleOAuth2Callback(c *gin.Context) {
	_, err := c.Cookie(oauthSessionName)
	if err != nil {
		log.Printf("corrupted session %s -- generated new", err)
		err = nil
	}

	stateChallenge, _ := c.Cookie(stateCallbackKey)
	state := c.Query("state")
	fmt.Printf("Post State: %s\n", state)
	switch stateChallenge, state := stateChallenge, state; {
	case state == "", len(stateChallenge) < 1:
		err = errors.New("missing state challenge")
	case state != stateChallenge:
		err = fmt.Errorf("invalid oauth state, expected '%s', got '%s'\n", state, stateChallenge)
	}

	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	token, err := oauth2Config.Exchange(context.Background(), c.Query("code"))
	if err != nil {
		return
	}

	// add the oauth token to session
	//session.Values[oauthTokenKey] = token

	a.env.TwitchSession.AccessToken = token.AccessToken
	a.env.TwitchSession.AccessTokenValidTill = token.Expiry
	a.env.TwitchSession.RefreshToken = token.RefreshToken

	fmt.Println("Access token loaded")

	c.Redirect(http.StatusTemporaryRedirect, "/")
	a.callback()

	return
}

// HumanReadableError represents error information
// that can be fed back to a human user.
//
// This prevents internal state that might be sensitive
// being leaked to the outside world.
type HumanReadableError interface {
	HumanError() string
	HTTPCode() int
}

// HumanReadableWrapper implements HumanReadableError
type HumanReadableWrapper struct {
	ToHuman string
	Code    int
	error
}

func (h HumanReadableWrapper) HumanError() string { return h.ToHuman }
func (h HumanReadableWrapper) HTTPCode() int      { return h.Code }

// AnnotateError wraps an error with a message that is intended for a human end-user to read,
// plus an associated HTTP error code.
func AnnotateError(err error, annotation string, code int) error {
	if err == nil {
		return nil
	}
	return HumanReadableWrapper{ToHuman: annotation, error: err}
}

type Handler func(http.ResponseWriter, *http.Request) error

func NewTwitchAuth(env *core.Environment, router *gin.Engine, callback func()) *TwitchAuth {
	gob.Register(&oauth2.Token{})

	oauth2Config = &oauth2.Config{
		ClientID:     env.Twitch.ClientId,
		ClientSecret: env.Twitch.ClientSecret,
		Scopes:       scopes,
		Endpoint:     twitch.Endpoint,
		RedirectURL:  redirectURL,
	}

	ta := &TwitchAuth{
		env:      env,
		callback: callback,
	}

	router.GET("/", ta.HandleRoot)
	router.GET("/login", ta.HandleLogin)
	router.GET("/redirect", ta.HandleOAuth2Callback)

	return ta
}
