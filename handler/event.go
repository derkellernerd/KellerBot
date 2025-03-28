package handler

import (
	"io"
	"log"
	"net/http"

	"github.com/derkellernerd/dori/core"
	"github.com/derkellernerd/dori/model"
	"github.com/derkellernerd/dori/repository"
	"github.com/gin-gonic/gin"
)

type Event struct {
	env          *core.Environment
	chatChannel  chan model.ChatEvent
	alertChannel chan model.Alert
	alertRepo    *repository.Alert
}

func NewEvent(env *core.Environment, chatChannel chan model.ChatEvent, alertChannel chan model.Alert, alertRepo *repository.Alert) *Event {
	return &Event{
		env:          env,
		chatChannel:  chatChannel,
		alertChannel: alertChannel,
		alertRepo:    alertRepo,
	}
}

func (h *Event) ChatEventHandler(c *gin.Context) {
	c.Stream(func(w io.Writer) bool {
		if msg, ok := <-h.chatChannel; ok {
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

	h.chatChannel <- message

	c.Status(http.StatusNoContent)
}

func (h *Event) AlertEventHandler(c *gin.Context) {
	c.Stream(func(w io.Writer) bool {
		if msg, ok := <-h.alertChannel; ok {
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

	h.alertChannel <- alert

	c.Status(http.StatusNoContent)
}
