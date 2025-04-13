package handler

import (
	"bytes"
	"html/template"
	"io"
	"log"
	"net/http"

	"github.com/derkellernerd/kellerbot/core"
	"github.com/derkellernerd/kellerbot/model"
	"github.com/derkellernerd/kellerbot/repository"
	twitchClient "github.com/derkellernerd/kellerbot/twitch"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Event struct {
	env          *core.Environment
	chatChannel  map[string]chan model.ChatEvent
	alertChannel map[string]chan []model.Alert
	alertRepo    *repository.Alert
}

func NewEvent(env *core.Environment, alertRepo *repository.Alert) *Event {
	return &Event{
		env:          env,
		alertRepo:    alertRepo,
		chatChannel:  make(map[string]chan model.ChatEvent),
		alertChannel: make(map[string]chan []model.Alert),
	}
}

func (h *Event) SendAlertEvent(alerts []model.Alert) {
	for _, alertChan := range h.alertChannel {
		alertChan <- alerts
	}
}

func (h *Event) SendChatEvent(chatEvent *model.ChatEvent) {
	for _, chatChan := range h.chatChannel {
		chatChan <- *chatEvent
	}
}

func (h *Event) ChatEventHandler(c *gin.Context) {
	id := uuid.NewString()

	h.chatChannel[id] = make(chan model.ChatEvent)

	c.Stream(func(w io.Writer) bool {
		if msg, ok := <-h.chatChannel[id]; ok {
			c.SSEvent("message", msg)
			return true
		}
		return false
	})

	return
}

func (h *Event) ChatEventTest(c *gin.Context) {
	var message model.ChatEvent

	err := c.BindJSON(&message)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, NewErrorResponse(err))
		return
	}

	h.SendChatEvent(&message)

	c.Status(http.StatusNoContent)
}

func (h *Event) AlertEventHandler(c *gin.Context) {
	id := uuid.NewString()

	h.alertChannel[id] = make(chan []model.Alert)
	log.Printf("Alert>Event: new client with id %s", id)

	c.Stream(func(w io.Writer) bool {
		if msg, ok := <-h.alertChannel[id]; ok {
			c.SSEvent("message", msg)
			return true
		}
		return false
	})

	log.Printf("Alert>Event: disconnecting client with id %s", id)
	delete(h.alertChannel, id)

	return
}

func GetAlertFinal(env *core.Environment, alertRepo *repository.Alert, alert *model.Alert, payload map[string]any, durationOverwrite float64) ([]model.Alert, error) {
	alerts := []model.Alert{}

	if durationOverwrite > 0 {
		alert.DurationInSeconds = durationOverwrite
	}

	switch alert.Type {
	case model.ALERT_TYPE_COMPOSITION:
		alertComposition, err := alert.GetDataComposition()
		if err != nil {
			return nil, err
		}

		for _, alertName := range alertComposition.AlertNames {
			childAlert, err := alertRepo.AlertFindByName(alertName)
			if err != nil {
				return nil, err
			}

			data, err := GetAlertFinal(env, alertRepo, &childAlert, payload, alert.DurationInSeconds)
			if err != nil {
				return nil, err
			}
			alerts = append(alerts, data...)
		}
		break
	case model.ALERT_TYPE_TEXT:
		alertText, err := alert.GetDataText()
		if err != nil {
			return nil, err
		}
		tmpl, err := template.New("test").Parse(alertText.Text)
		if err != nil {
			return nil, err
		}
		var tpl bytes.Buffer
		err = tmpl.Execute(&tpl, payload)

		alertText.Text = tpl.String()
		alert.SetData(alertText)
		alerts = append(alerts, *alert)
		break
	case model.ALERT_TYPE_CHAT:
		alertText, err := alert.GetDataChatText()
		if err != nil {
			return nil, err
		}
		tmpl, err := template.New("test").Parse(alertText.Chat)
		if err != nil {
			return nil, err
		}
		var tpl bytes.Buffer
		err = tmpl.Execute(&tpl, payload)
		twitchClient.SendChatMessage(env, tpl.String())
		break
	default:
		alerts = append(alerts, *alert)
	}

	return alerts, nil
}

func (h *Event) AlertEventTest(c *gin.Context) {
	var alertEvent model.AlertEventRequest

	err := c.BindJSON(&alertEvent)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, NewErrorResponse(err))
		return
	}

	alert, err := h.alertRepo.AlertFindByName(alertEvent.Name)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, NewErrorResponse(err))
		return
	}
	log.Printf("Found Alert: %d", alert.ID)

	alerts, err := GetAlertFinal(h.env, h.alertRepo, &alert, alertEvent.Payload, 0)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, NewErrorResponse(err))
		return
	}

	h.SendAlertEvent(alerts)

	c.Status(http.StatusNoContent)
}

func (h *Event) Status(c *gin.Context) {
	type Status struct {
		ChatClients      []string
		AlertClients     []string
		ChatClientCount  int
		AlertClientCount int
	}

	status := Status{
		ChatClientCount:  len(h.chatChannel),
		AlertClientCount: len(h.alertChannel),
	}

	c.JSON(http.StatusOK, status)
}
