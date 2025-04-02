package handler

import (
	"io"
	"log"
	"net/http"

	"github.com/derkellernerd/dori/core"
	"github.com/derkellernerd/dori/model"
	"github.com/derkellernerd/dori/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Event struct {
	env          *core.Environment
	chatChannel  map[string]chan model.ChatEvent
	alertChannel map[string]chan model.Alert
	alertRepo    *repository.Alert
}

func NewEvent(env *core.Environment, alertRepo *repository.Alert) *Event {
	return &Event{
		env:          env,
		alertRepo:    alertRepo,
		chatChannel:  make(map[string]chan model.ChatEvent),
		alertChannel: make(map[string]chan model.Alert),
	}
}

func (h *Event) SendAlertEvent(alert *model.Alert) {
	for _, alertChan := range h.alertChannel {
		alertChan <- *alert
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

	h.alertChannel[id] = make(chan model.Alert)

	c.Stream(func(w io.Writer) bool {
		if msg, ok := <-h.alertChannel[id]; ok {
			c.SSEvent("message", msg)
			return true
		}
		return false
	})

	return
}

func (h *Event) AlertEventTest(c *gin.Context) {
	var alertEvent model.AlertEvent

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

	h.SendAlertEvent(&alert)

	c.Status(http.StatusNoContent)
}
