package handler

import (
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/derkellernerd/kellerbot/core"
	"github.com/derkellernerd/kellerbot/repository"
	"github.com/derkellernerd/kellerbot/worker"
	"github.com/gin-gonic/gin"
)

type Event struct {
	env          *core.Environment
	actionWorker *worker.Action
	events       *repository.Event
}

func NewEvent(env *core.Environment, actionWorker *worker.Action, eventRepo *repository.Event) *Event {
	return &Event{
		env:          env,
		actionWorker: actionWorker,
		events:       eventRepo,
	}
}

func (h *Event) EventGetAll(c *gin.Context) {
	events, err := h.events.EventFindAll()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse(events))
}

func (h *Event) EventActionReplay(c *gin.Context) {
	eventIdParam := c.Param("eventId")
	eventId, err := strconv.ParseUint(eventIdParam, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, NewErrorResponse(err))
		return
	}

	event, err := h.events.EventFindById(uint(eventId))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, NewErrorResponse(err))
		return
	}

	payload, err := event.GetPayload()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, NewErrorResponse(err))
		return
	}

	err = h.actionWorker.HandleActionByName(event.ExecutingActionName, payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, NewErrorResponse(err))
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Event) AlertEventHandler(c *gin.Context) {
	id, eventChannel := h.actionWorker.RegisterListener()
	c.Writer.CloseNotify()
	w := c.Writer
	clientGone := w.CloseNotify()

	c.Stream(func(w io.Writer) bool {
		select {
		case <-clientGone:
			log.Printf("Event Handler>Alert>Client Gone: %s", id)
			return true
		default:
			if msg, ok := <-eventChannel; ok {
				c.SSEvent("message", msg)
				return true
			}
			return false
		}
	})

	h.actionWorker.UnregisterListener(id)

	c.Status(http.StatusGone)
}
