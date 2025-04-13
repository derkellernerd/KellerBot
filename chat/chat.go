package chat

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/derkellernerd/kellerbot/core"
	"github.com/derkellernerd/kellerbot/handler"
	"github.com/derkellernerd/kellerbot/model"
	"github.com/derkellernerd/kellerbot/repository"
	twitchClient "github.com/derkellernerd/kellerbot/twitch"
	"github.com/joeyak/go-twitch-eventsub/v3"
)

type Chat struct {
	client       *twitch.Client
	env          *core.Environment
	commandRepo  *repository.Command
	alertRepo    *repository.Alert
	eventHandler *handler.Event
	twitchEvent  *repository.TwitchEvent
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
			twitch.SubChannelRaid,
		}

		for _, event := range events {
			fmt.Printf("subscribing to %s as %s\n", event, c.env.Twitch.UserId)
			_, err := twitch.SubscribeEvent(twitch.SubscribeRequest{
				SessionID:   message.Payload.Session.ID,
				ClientID:    c.env.Twitch.ClientId,
				AccessToken: c.env.TwitchSession.AccessToken,
				Event:       event,
				Condition: map[string]string{
					"broadcaster_user_id":    c.env.Twitch.UserId,
					"user_id":                c.env.Twitch.UserId,
					"moderator_user_id":      c.env.Twitch.UserId,
					"to_broadcaster_user_id": c.env.Twitch.UserId,
				},
			})
			if err != nil {
				fmt.Printf("ERROR: %v\n", err)
				return
			}
		}
	})
	c.client.OnEventChannelRaid(func(raid twitch.EventChannelRaid) {
		event, err := c.twitchEvent.TwitchEventFindByTwitchEventSubscripton(string(twitch.SubChannelRaid))
		if err != nil {
			if err == repository.ErrTwitchEventNotFound {
				return
			}
			log.Println(err)
		}

		alert, err := c.alertRepo.AlertFindByName(event.AlertName)
		if err != nil {
			log.Println(err)
		}

		payload := map[string]any{
			"from_broadcaster_user_name": raid.FromBroadcasterUserName,
			"viewers":                    raid.Viewers,
		}

		go func() {
			log.Printf("Sending alert %s", alert.Name)
			alerts, err := handler.GetAlertFinal(c.env, c.alertRepo, &alert, payload, 0)
			if err != nil {
				log.Println(err)
				return
			}
			c.eventHandler.SendAlertEvent(alerts)
		}()
	})
	c.client.OnEventChannelFollow(func(follow twitch.EventChannelFollow) {
		event, err := c.twitchEvent.TwitchEventFindByTwitchEventSubscripton(string(twitch.SubChannelFollow))
		if err != nil {
			if err == repository.ErrTwitchEventNotFound {
				return
			}
			log.Println(err)
		}

		alert, err := c.alertRepo.AlertFindByName(event.AlertName)
		if err != nil {
			log.Println(err)
		}

		payload := map[string]any{
			"user_name": follow.UserName,
		}

		go func() {
			log.Printf("Sending alert %s", alert.Name)
			alerts, err := handler.GetAlertFinal(c.env, c.alertRepo, &alert, payload, 0)
			if err != nil {
				log.Println(err)
				return
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

				twitchClient.SendChatMessage(c.env, fmt.Sprintf("commands available: %s", result))
				return
			}

			cmd, err := c.commandRepo.CommandFindByCommand(command)
			if err != nil {
				log.Println(err)
			}

			durationBetweenExecutions := time.Now().Sub(cmd.LastUsed).Seconds()
			log.Printf("Last used: %f seconds", durationBetweenExecutions)
			if cmd.TimeoutInSeconds > 0 && durationBetweenExecutions < float64(cmd.TimeoutInSeconds) {
				twitchClient.SendChatAnswer(c.env, message.MessageId, "Command noch %d Sekunden nicht verfuegbar", cmd.TimeoutInSeconds-uint64(durationBetweenExecutions))
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
				twitchClient.SendChatMessage(c.env, message.Message)
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
					twitchClient.SendChatAnswer(c.env, message.MessageId, "SirSad: %s", err)
					return
				}
				twitchClient.SendChatAnswer(c.env, message.MessageId, "ICH MACH LICHT!")
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
					alerts, err := handler.GetAlertFinal(c.env, c.alertRepo, &alert, make(map[string]any), 0)
					if err != nil {
						log.Println(err)
						return
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

func NewChat(env *core.Environment, commandRepo *repository.Command, alertRepo *repository.Alert, eventHandler *handler.Event, twitchEventRepo *repository.TwitchEvent) (*Chat, error) {
	client := twitch.NewClient()

	return &Chat{
		client:       client,
		env:          env,
		commandRepo:  commandRepo,
		alertRepo:    alertRepo,
		eventHandler: eventHandler,
		twitchEvent:  twitchEventRepo,
	}, nil
}
