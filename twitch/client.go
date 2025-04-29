package twitch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/derkellernerd/kellerbot/core"
)

type TwitchChatMessage struct {
	BroadcasterId   string `json:"broadcaster_id"`
	SenderId        string `json:"sender_id"`
	Message         string `json:"message"`
	ParentMessageId string `json:"reply_parent_message_id"`
}

func httpApiRequest(env *core.Environment, endpoint string, payload any) error {
	uri := fmt.Sprintf("https://api.twitch.tv/helix/%s", endpoint)

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, uri, bytes.NewReader(payloadBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Client-Id", env.Twitch.ClientId)
	//	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", env.TwitchSession.AccessToken))

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

func SendChatAnswer(env *core.Environment, messageId string, message string, a ...any) error {
	payload := TwitchChatMessage{
		//	BroadcasterId: env.Twitch.UserId,
		//	SenderId:      env.Twitch.UserId,
		Message: fmt.Sprintf(message, a...),
	}

	if messageId != "" {
		payload.ParentMessageId = messageId
	}

	return httpApiRequest(env, "chat/messages", &payload)
}

func SendChatMessage(env *core.Environment, message string, a ...any) error {
	return SendChatAnswer(env, "", message, a...)
}
