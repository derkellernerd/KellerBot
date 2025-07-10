package chat

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/derkellernerd/kellerbot/core"
	"github.com/derkellernerd/kellerbot/model"
	"github.com/derkellernerd/kellerbot/repository"
	twitchClient "github.com/derkellernerd/kellerbot/twitch"
	"github.com/derkellernerd/kellerbot/worker"
	"github.com/joeyak/go-twitch-eventsub/v3"
)

type Chat struct {
	client          *twitch.Client
	env             *core.Environment
	twitchEvent     *repository.TwitchEvent
	actionWorker    *worker.Action
	chatCommand     *repository.ChatCommand
	eventRepository *repository.Event
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
			twitch.SubChannelSubscribe,
			twitch.SubChannelSubscriptionGift,
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
	c.client.OnEventChannelSubscribe(func(sub twitch.EventChannelSubscribe) {
		event, err := c.twitchEvent.TwitchEventFindByTwitchEventSubscripton(string(twitch.SubChannelSubscribe))
		if err != nil {
			if err == repository.ErrTwitchEventNotFound {
				return
			}
			log.Println(err)
		}

		payload := map[string]any{
			"user_name": sub.UserName,
			"tier":      sub.Tier,
		}

		err = c.actionWorker.HandleActionByName(event.ActionName, payload)
		if err != nil {
			log.Println(err)
		}

		eventLog := model.Event{
			EventName:           string(twitch.SubChannelSubscribe),
			Source:              model.EVENT_SOURCE_TWITCH,
			ExecutingActionName: event.ActionName,
		}

		eventLog.SetPayload(payload)
		err = c.eventRepository.EventInsert(&eventLog)
		if err != nil {
			log.Println(err)
		}
	})

	c.client.OnEventChannelSubscriptionGift(func(sub twitch.EventChannelSubscriptionGift) {
		event, err := c.twitchEvent.TwitchEventFindByTwitchEventSubscripton(string(twitch.SubChannelSubscriptionGift))
		if err != nil {
			if err == repository.ErrTwitchEventNotFound {
				return
			}
			log.Println(err)
		}

		log.Printf("%#v", sub)

		payload := map[string]any{
			"user_name": sub.UserName,
			"tier":      sub.Tier,
		}

		err = c.actionWorker.HandleActionByName(event.ActionName, payload)
		if err != nil {
			log.Println(err)
		}

		eventLog := model.Event{
			EventName:           string(twitch.SubChannelSubscriptionGift),
			Source:              model.EVENT_SOURCE_TWITCH,
			ExecutingActionName: event.ActionName,
		}

		eventLog.SetPayload(payload)
		err = c.eventRepository.EventInsert(&eventLog)
		if err != nil {
			log.Println(err)
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

		payload := map[string]any{
			"from_broadcaster_user_name": raid.FromBroadcasterUserName,
			"viewers":                    raid.Viewers,
		}

		err = c.actionWorker.HandleActionByName(event.ActionName, payload)
		if err != nil {
			log.Println(err)
		}

		eventLog := model.Event{
			EventName:           string(twitch.SubChannelRaid),
			Source:              model.EVENT_SOURCE_TWITCH,
			ExecutingActionName: event.ActionName,
		}

		eventLog.SetPayload(payload)
		err = c.eventRepository.EventInsert(&eventLog)
		if err != nil {
			log.Println(err)
		}
	})
	c.client.OnEventChannelFollow(func(follow twitch.EventChannelFollow) {
		event, err := c.twitchEvent.TwitchEventFindByTwitchEventSubscripton(string(twitch.SubChannelFollow))
		if err != nil {
			if err == repository.ErrTwitchEventNotFound {
				return
			}
			log.Println(err)
		}

		payload := map[string]any{
			"user_name": follow.UserName,
		}
		err = c.actionWorker.HandleActionByName(event.ActionName, payload)
		if err != nil {
			log.Println(err)
		}

		eventLog := model.Event{
			EventName:           string(twitch.SubChannelFollow),
			Source:              model.EVENT_SOURCE_TWITCH,
			ExecutingActionName: event.ActionName,
		}

		eventLog.SetPayload(payload)
		err = c.eventRepository.EventInsert(&eventLog)
		if err != nil {
			log.Println(err)
		}
	})
	c.client.OnEventChannelChatMessage(func(message twitch.EventChannelChatMessage) {
		if strings.HasPrefix(message.Message.Text, "!") {
			commandParts := strings.Split(message.Message.Text, " ")
			command := strings.ToLower(commandParts[0])
			args := commandParts[1:]
			command = command[1:]

			if command == "commands" {
				commands, err := c.chatCommand.ChatCommandFindAll()
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
			cmd, err := c.chatCommand.ChatCommandFindByChatCommand(command)
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

			payload := map[string]any{
				"args":       args,
				"message_id": message.MessageId,
			}

			err = c.actionWorker.HandleActionByName(cmd.Action, payload)
			if err != nil {
				log.Println(err)
			}

			c.chatCommand.ChatCommandUpdate(&cmd)
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

func NewChat(env *core.Environment, twitchEventRepo *repository.TwitchEvent, chatCommandRepo *repository.ChatCommand, actionWorker *worker.Action, eventRepo *repository.Event) (*Chat, error) {
	client := twitch.NewClient()

	return &Chat{
		client:          client,
		env:             env,
		twitchEvent:     twitchEventRepo,
		actionWorker:    actionWorker,
		chatCommand:     chatCommandRepo,
		eventRepository: eventRepo,
	}, nil
}
