package handler

import (
	"net/http"
	"strconv"

	"github.com/derkellernerd/dori/core"
	"github.com/derkellernerd/dori/model"
	"github.com/derkellernerd/dori/repository"
	"github.com/gin-gonic/gin"
)

type TwitchEvent struct {
	env             *core.Environment
	twitchEventRepo *repository.TwitchEvent
}

func NewTwitchEvent(env *core.Environment, twitchEventRepo *repository.TwitchEvent) *TwitchEvent {
	return &TwitchEvent{
		env:             env,
		twitchEventRepo: twitchEventRepo,
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
		AlertName:               twitchEventCreateRequest.AlertName,
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

	twitchEvent.AlertName = twitchEventUpdateRequest.AlertName

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
