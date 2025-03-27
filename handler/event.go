package handler

import (
	"io"
	"net/http"

	"github.com/derkellernerd/dori/core"
	"github.com/derkellernerd/dori/model"
	"github.com/gin-gonic/gin"
)

type Event struct {
	env         *core.Environment
	chatChannel chan model.ChatEvent
}

func NewEvent(env *core.Environment, chatChannel chan model.ChatEvent) *Event {
	return &Event{
		env:         env,
		chatChannel: chatChannel,
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
