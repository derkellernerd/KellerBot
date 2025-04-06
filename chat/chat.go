package chat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/derkellernerd/dori/core"
	"github.com/derkellernerd/dori/handler"
	"github.com/derkellernerd/dori/model"
	"github.com/derkellernerd/dori/repository"
	"github.com/joeyak/go-twitch-eventsub/v3"
)

type Chat struct {
	client       *twitch.Client
	env          *core.Environment
	commandRepo  *repository.Command
	alertRepo    *repository.Alert
	eventHandler *handler.Event
}

type TwitchChatMessage struct {
	BroadcasterId   string `json:"broadcaster_id"`
	SenderId        string `json:"sender_id"`
	Message         string `json:"message"`
	ParentMessageId string `json:"reply_parent_message_id"`
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
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.env.TwitchSession.AccessToken))

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

func (c *Chat) SendChatAnswer(messageId string, message string, a ...any) error {
	payload := TwitchChatMessage{
		BroadcasterId: c.env.Twitch.UserId,
		SenderId:      c.env.Twitch.UserId,
		Message:       fmt.Sprintf(message, a...),
	}

	if messageId != "" {
		payload.ParentMessageId = messageId
	}

	return c.httpApiRequest("chat/messages", &payload)
}

func (c *Chat) SendChatMessage(message string, a ...any) error {
	return c.SendChatAnswer("", message, a...)
}

func (c *Chat) Start() error {
	c.client.OnError(func(err error) {
		fmt.Printf("ERROR: %v\n", err)
	})
	c.client.OnWelcome(func(message twitch.WelcomeMessage) {
		fmt.Printf("WELCOME: %v\n", message)

		events := []twitch.EventSubscription{
			twitch.SubChannelChatMessage,
			twitch.SubChannelFollow,
		}

		for _, event := range events {
			fmt.Printf("subscribing to %s as %s\n", event, c.env.Twitch.UserId)
			_, err := twitch.SubscribeEvent(twitch.SubscribeRequest{
				SessionID:   message.Payload.Session.ID,
				ClientID:    c.env.Twitch.ClientId,
				AccessToken: c.env.TwitchSession.AccessToken,
				Event:       event,
				Condition: map[string]string{
					"broadcaster_user_id": c.env.Twitch.UserId,
					"user_id":             c.env.Twitch.UserId,
					"moderator_user_id":   c.env.Twitch.UserId,
				},
			})
			if err != nil {
				fmt.Printf("ERROR: %v\n", err)
				return
			}
		}
	})
	c.client.OnEventChannelFollow(func(follow twitch.EventChannelFollow) {
		alert, err := c.alertRepo.AlertFindByName("minion_horn")
		if err != nil {
			log.Println(err)
		}

		go func() {
			log.Printf("Sending alert %s", alert.Name)
			alerts := []model.Alert{}

			if alert.Type == model.ALERT_TYPE_COMPOSITION {
				alertComposition, err := alert.GetDataComposition()
				if err != nil {
					log.Println(err)
				}

				for _, alertName := range alertComposition.AlertNames {
					childAlert, err := c.alertRepo.AlertFindByName(alertName)
					if err != nil {
						log.Println(err)
					}

					alerts = append(alerts, childAlert)
				}
			} else {
				alerts = append(alerts, alert)
			}
			c.eventHandler.SendAlertEvent(alerts)
		}()
	})
	c.client.OnEventChannelChatMessage(func(message twitch.EventChannelChatMessage) {
		chatEvent := model.ChatEvent{
			Message: message.Message.Text,
			User:    message.ChatterUserName,
		}
		go func() {
			c.eventHandler.SendChatEvent(&chatEvent)
		}()

		if strings.HasPrefix(message.Message.Text, "!") {
			log.Println("Incoming command")
			commandParts := strings.Split(message.Message.Text, " ")
			command := strings.ToLower(commandParts[0])
			args := commandParts[1:]
			command = command[1:]

			if command == "commands" {
				commands, err := c.commandRepo.CommandFindAll()
				if err != nil {
					log.Println(err)
				}

				result := ""
				for _, cmd := range commands {
					result = fmt.Sprintf("%s !%s", result, cmd.Command)
				}

				c.SendChatMessage(fmt.Sprintf("commands available: %s", result))
				return
			}

			cmd, err := c.commandRepo.CommandFindByCommand(command)
			if err != nil {
				log.Println(err)
			}

			durationBetweenExecutions := time.Now().Sub(cmd.LastUsed).Seconds()
			log.Printf("Last used: %f seconds", durationBetweenExecutions)
			if cmd.TimeoutInSeconds > 0 && durationBetweenExecutions < float64(cmd.TimeoutInSeconds) {
				c.SendChatAnswer(message.MessageId, "Command noch %d Sekunden nicht verfuegbar", cmd.TimeoutInSeconds-uint64(durationBetweenExecutions))
				return
			}

			cmd.Increment()
			cmd.LastUsed = time.Now()

			switch cmd.Type {
			case model.COMMAND_TYPE_MESSAGE:
				message, err := cmd.GetDataActionMessage()
				if err != nil {
					log.Println(err)
				}
				c.SendChatMessage(message.Message)
				break
			case model.COMMAND_TYPE_HTTP_ACTION:
				httpAction, err := cmd.GetDataActionHttp()
				if err != nil {
					log.Println(err)
				}

				numberArgs := []any{}
				for _, arg := range args {
					val, err := strconv.ParseUint(arg, 10, 64)
					if err != nil {
						log.Println(err)
					}

					numberArgs = append(numberArgs, val)
				}

				err = httpAction.Do(numberArgs...)
				if err != nil {
					c.SendChatAnswer(message.MessageId, "SirSad: %s", err)
					return
				}
				c.SendChatAnswer(message.MessageId, "ICH MACH LICHT!")
				break
			case model.COMMAND_TYPE_ALERT_ACTION:
				alertAction, err := cmd.GetDataActionAlert()
				if err != nil {
					log.Println(err)
				}

				alert, err := c.alertRepo.AlertFindByName(alertAction.Alert)
				if err != nil {
					log.Println(err)
				}

				go func() {
					log.Printf("Sending alert %s", alert.Name)
					alerts := []model.Alert{}

					if alert.Type == model.ALERT_TYPE_COMPOSITION {
						alertComposition, err := alert.GetDataComposition()
						if err != nil {
							log.Println(err)
						}

						for _, alertName := range alertComposition.AlertNames {
							childAlert, err := c.alertRepo.AlertFindByName(alertName)
							if err != nil {
								log.Println(err)
							}

							alerts = append(alerts, childAlert)
						}
					} else {
						alerts = append(alerts, alert)
					}
					c.eventHandler.SendAlertEvent(alerts)
				}()
			}

			c.commandRepo.CommandUpdate(&cmd)
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

func NewChat(env *core.Environment, commandRepo *repository.Command, alertRepo *repository.Alert, eventHandler *handler.Event) (*Chat, error) {
	client := twitch.NewClient()

	return &Chat{
		client:       client,
		env:          env,
		commandRepo:  commandRepo,
		alertRepo:    alertRepo,
		eventHandler: eventHandler,
	}, nil
}
