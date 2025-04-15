package worker

import (
	"bytes"
	"log"
	"text/template"

	"github.com/derkellernerd/kellerbot/core"
	"github.com/derkellernerd/kellerbot/model"
	"github.com/derkellernerd/kellerbot/repository"
	twitchClient "github.com/derkellernerd/kellerbot/twitch"
	"github.com/google/uuid"
)

type Action struct {
	env           *core.Environment
	actionRepo    *repository.Action
	actionChannel map[string]chan model.Action
}

func NewAction(env *core.Environment, actionRepo *repository.Action) *Action {
	return &Action{
		env:           env,
		actionRepo:    actionRepo,
		actionChannel: make(map[string]chan model.Action),
	}
}

func (a *Action) RegisterListener() (string, chan model.Action) {
	id := uuid.NewString()
	log.Printf("Action Worker>RegisterListener: %s", id)
	a.actionChannel[id] = make(chan model.Action)

	return id, a.actionChannel[id]
}

func (a *Action) UnregisterListener(id string) {
	log.Printf("Action Worker>UnregisterListener: %s", id)
	delete(a.actionChannel, id)
}

func (a *Action) renderText(templateText string, payload map[string]any) (string, error) {
	tmpl, err := template.New("test").Parse(templateText)
	if err != nil {
		return "", err
	}
	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, payload)
	return tpl.String(), nil
}

func (a *Action) HandleActionByName(actionName string, payload map[string]any) error {
	log.Printf("Action Worker>HandleActionByName: %s", actionName)
	action, err := a.actionRepo.ActionFindByActionName(actionName)
	if err != nil {
		return err
	}
	return a.HandleAction(action, payload, 0)
}

func (a *Action) HandleAction(action model.Action, payload map[string]any, forceDurationMs int64) error {
	log.Printf("Action Worker>HandleAction: %s (%d)", action.ActionName, action.ID)
	publishAction := true
	switch action.ActionType {
	case model.ACTION_TYPE_COMPOSITION:
		actionComposition, err := model.ActionGetData[model.ActionTypeComposition](&action)
		if err != nil {
			return err
		}

		if forceDurationMs > 0 {
			actionComposition.DurationMs = forceDurationMs
		}

		for _, actionName := range actionComposition.Actions {
			childAction, err := a.actionRepo.ActionFindByActionName(actionName)
			if err != nil {
				return err
			}

			err = a.HandleAction(childAction, payload, actionComposition.DurationMs)
			if err != nil {
				return err
			}
		}
		publishAction = false
		break
	case model.ACTION_TYPE_CHAT_ANSWER:
		actionChatAnswer, err := model.ActionGetData[model.ActionTypeChatAnswer](&action)
		if err != nil {
			return err
		}

		rendered, err := a.renderText(actionChatAnswer.ChatMessage, payload)
		if err != nil {
			return err
		}
		err = twitchClient.SendChatAnswer(a.env, payload["message_id"].(string), rendered)
		if err != nil {
			return err
		}
		publishAction = false

		break
	case model.ACTION_TYPE_CHAT_MESSAGE:
		actionChatMessage, err := model.ActionGetData[model.ActionTypeChatMessage](&action)
		if err != nil {
			return err
		}

		rendered, err := a.renderText(actionChatMessage.ChatMessage, payload)
		if err != nil {
			return err
		}

		err = twitchClient.SendChatMessage(a.env, rendered)
		if err != nil {
			return err
		}
		publishAction = false
		break
	case model.ACTION_TYPE_TEXT:
		actionText, err := model.ActionGetData[model.ActionTypeText](&action)
		if err != nil {
			return err
		}

		rendered, err := a.renderText(actionText.Text, payload)
		if err != nil {
			return err
		}
		actionText.Text = rendered
		action.SetData(actionText)
		break
	case model.ACTION_TYPE_GIF:
		actionGif, err := model.ActionGetData[model.ActionTypeGif](&action)
		if err != nil {
			return err
		}

		if forceDurationMs > 0 {
			actionGif.DurationMs = forceDurationMs
		}

		action.SetData(actionGif)
	case model.ACTION_TYPE_SOUND:
		actionSound, err := model.ActionGetData[model.ActionTypeSound](&action)
		if err != nil {
			return err
		}

		if forceDurationMs > 0 {
			actionSound.DurationMs = forceDurationMs
		}

		action.SetData(actionSound)
	}

	if publishAction {
		for id, actionChannel := range a.actionChannel {
			log.Printf("Action Worker>HandleAction>Publish: %s", id)
			actionChannel <- action
		}
	}
	return nil
}
