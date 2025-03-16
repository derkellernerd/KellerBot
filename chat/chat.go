package chat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/derkellernerd/dori/core"
	"github.com/joeyak/go-twitch-eventsub/v3"
)

type Chat struct {
	client *twitch.Client
	env    *core.Environment
}

type TwitchChatMessage struct {
	BroadcasterId string `json:"broadcaster_id"`
	SenderId      string `json:"sender_id"`
	Message       string `json:"message"`
}

func (c *Chat) httpApiRequest(endpoint string, payload any) error {
	uri := fmt.Sprintf("https://api.twitch.tv/helix/%s", endpoint)

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, uri, bytes.NewReader(payloadBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Client-Id", c.env.Twitch.ClientId)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.env.Twitch.AccessToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := io.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))

	return nil
}

func (c *Chat) SendChatMessage(message string) error {
	payload := TwitchChatMessage{
		BroadcasterId: c.env.Twitch.UserId,
		SenderId:      c.env.Twitch.UserId,
		Message:       message,
	}
	return c.httpApiRequest("chat/messages", &payload)
}

func (c *Chat) Start() error {
	c.client.OnError(func(err error) {
		fmt.Printf("ERROR: %v\n", err)
	})
	c.client.OnWelcome(func(message twitch.WelcomeMessage) {
		fmt.Printf("WELCOME: %v\n", message)

		events := []twitch.EventSubscription{
			twitch.SubChannelChatMessage,
		}

		for _, event := range events {
			fmt.Printf("subscribing to %s as %s\n", event, c.env.Twitch.UserId)
			_, err := twitch.SubscribeEvent(twitch.SubscribeRequest{
				SessionID:   message.Payload.Session.ID,
				ClientID:    c.env.Twitch.ClientId,
				AccessToken: c.env.Twitch.AccessToken,
				Event:       event,
				Condition: map[string]string{
					"broadcaster_user_id": c.env.Twitch.UserId,
					"user_id":             c.env.Twitch.UserId,
				},
			})
			if err != nil {
				fmt.Printf("ERROR: %v\n", err)
				return
			}
		}
	})
	c.client.OnEventChannelChatMessage(func(message twitch.EventChannelChatMessage) {
		fmt.Printf("Hello chat incoming: %s\n", message.Message.Text)
		if message.Message.Text == "!ping" {
			c.SendChatMessage("pong")
		}
	})

	c.client.OnKeepAlive(func(message twitch.KeepAliveMessage) {
		fmt.Printf("KEEPALIVE: %v\n", message)
	})
	c.client.OnRevoke(func(message twitch.RevokeMessage) {
		fmt.Printf("REVOKE: %v\n", message)
	})
	c.client.OnRawEvent(func(event string, metadata twitch.MessageMetadata, subscription twitch.PayloadSubscription) {
		fmt.Printf("EVENT[%s]: %s: %s\n", subscription.Type, metadata, event)
	})

	err := c.client.Connect()
	return err
}

func NewChat(env *core.Environment) (*Chat, error) {
	client := twitch.NewClient()

	return &Chat{
		client: client,
		env:    env,
	}, nil
}
