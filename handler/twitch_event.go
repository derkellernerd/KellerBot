package handler

import (
	"net/http"
	"strconv"

	"github.com/derkellernerd/kellerbot/core"
	"github.com/derkellernerd/kellerbot/model"
	"github.com/derkellernerd/kellerbot/repository"
	"github.com/derkellernerd/kellerbot/worker"
	"github.com/gin-gonic/gin"
)

type TwitchEvent struct {
	env             *core.Environment
	twitchEventRepo *repository.TwitchEvent
	actionWorker    *worker.Action
}

func NewTwitchEvent(env *core.Environment, twitchEventRepo *repository.TwitchEvent, actionWorker *worker.Action) *TwitchEvent {
	return &TwitchEvent{
		env:             env,
		twitchEventRepo: twitchEventRepo,
		actionWorker:    actionWorker,
	}
}

func (h *TwitchEvent) TwitchEventGetAll(c *gin.Context) {
	twitchEvents, err := h.twitchEventRepo.TwitchEventFindAll()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse(twitchEvents))
}

func (h *TwitchEvent) TwitchEventCreate(c *gin.Context) {
	var twitchEventCreateRequest model.TwitchEventCreateRequest

	err := c.BindJSON(&twitchEventCreateRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, NewErrorResponse(err))
		return
	}

	twitchEvent := model.TwitchEvent{
		TwitchEventSubscription: twitchEventCreateRequest.TwitchEventSubscription,
		ActionName:              twitchEventCreateRequest.ActionName,
	}

	err = h.twitchEventRepo.TwitchEventInsert(&twitchEvent)

	c.JSON(http.StatusCreated, NewSuccessResponse(twitchEvent))
}

func (h *TwitchEvent) TwitchEventUpdate(c *gin.Context) {
	twitchEventIdParam := c.Param("twitchEventId")
	twitchEventId, err := strconv.ParseUint(twitchEventIdParam, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, NewErrorResponse(err))
		return
	}

	var twitchEventUpdateRequest model.TwitchEventUpdateRequest
	err = c.BindJSON(&twitchEventUpdateRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, NewErrorResponse(err))
		return
	}

	twitchEvent, err := h.twitchEventRepo.TwitchEventFindById(uint(twitchEventId))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, NewErrorResponse(err))
		return
	}

	twitchEvent.ActionName = twitchEventUpdateRequest.ActionName

	err = h.twitchEventRepo.TwitchEventUpdate(&twitchEvent)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse(twitchEvent))
}

func (h *TwitchEvent) TwitchEventDelete(c *gin.Context) {
	twitchEventIdParam := c.Param("twitchEventId")
	twitchEventId, err := strconv.ParseUint(twitchEventIdParam, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, NewErrorResponse(err))
		return
	}

	twitchEvent, err := h.twitchEventRepo.TwitchEventFindById(uint(twitchEventId))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, NewErrorResponse(err))
		return
	}

	err = h.twitchEventRepo.TwitchEventDelete(&twitchEvent)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, NewErrorResponse(err))
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *TwitchEvent) TwitchEventTest(c *gin.Context) {
	twitchEventIdParam := c.Param("twitchEventId")
	twitchEventId, err := strconv.ParseUint(twitchEventIdParam, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, NewErrorResponse(err))
		return
	}

	twitchEvent, err := h.twitchEventRepo.TwitchEventFindById(uint(twitchEventId))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, NewErrorResponse(err))
		return
	}

	payload := map[string]any{}
	switch twitchEvent.TwitchEventSubscription {
	case "channel.raid":
		payload["from_broadcaster_user_name"] = "derkellerbot"
		payload["viewers"] = 9001
		break
	case "channel.follow":
		payload["user_name"] = "derkellerbot"
		break
	}

	err = h.actionWorker.HandleActionByName(twitchEvent.ActionName, payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, NewErrorResponse(err))
		return
	}

	c.Status(http.StatusNoContent)
}
