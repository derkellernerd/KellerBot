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

type ChatCommand struct {
	env          *core.Environment
	chatCommands *repository.ChatCommand
	actionWorker *worker.Action
}

func NewChatCommand(env *core.Environment, chatCommandRepo *repository.ChatCommand, actionWorker *worker.Action) *ChatCommand {
	return &ChatCommand{
		env:          env,
		chatCommands: chatCommandRepo,
		actionWorker: actionWorker,
	}
}

func (h *ChatCommand) ChatCommandGetAll(c *gin.Context) {
	chatCommands, err := h.chatCommands.ChatCommandFindAll()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse(chatCommands))
}

func (h *ChatCommand) ChatCommandCreate(c *gin.Context) {
	var chatCommandCreateRequest model.ChatCommandCreateRequest

	err := c.BindJSON(&chatCommandCreateRequest)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, NewErrorResponse(err))
		return
	}

	chatCommand := model.ChatCommand{
		Command:          chatCommandCreateRequest.Command,
		Action:           chatCommandCreateRequest.Action,
		TimeoutInSeconds: chatCommandCreateRequest.TimeoutInSeconds,
	}

	err = h.chatCommands.ChatCommandInsert(&chatCommand)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, NewSuccessResponse(chatCommand))
}

func (h *ChatCommand) ChatCommandUpdate(c *gin.Context) {
	chatCommandIdParam := c.Param("chatCommandId")
	chatCommandId, err := strconv.ParseUint(chatCommandIdParam, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, NewErrorResponse(err))
		return
	}

	var chatCommandUpdateRequest model.ChatCommandUpdateRequest
	err = c.BindJSON(&chatCommandUpdateRequest)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, NewErrorResponse(err))
		return
	}

	chatCommand, err := h.chatCommands.ChatCommandFindById(uint(chatCommandId))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, NewErrorResponse(err))
		return
	}

	chatCommand.TimeoutInSeconds = chatCommandUpdateRequest.TimeoutInSeconds
	chatCommand.Action = chatCommandUpdateRequest.Action

	err = h.chatCommands.ChatCommandUpdate(&chatCommand)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, NewSuccessResponse(chatCommand))
}

func (h *ChatCommand) ChatCommandDelete(c *gin.Context) {
	chatCommandIdParam := c.Param("chatCommandId")
	chatCommandId, err := strconv.ParseUint(chatCommandIdParam, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, NewErrorResponse(err))
		return
	}

	chatCommand, err := h.chatCommands.ChatCommandFindById(uint(chatCommandId))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, NewErrorResponse(err))
		return
	}

	err = h.chatCommands.ChatCommandDelete(&chatCommand)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, NewErrorResponse(err))
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *ChatCommand) ChatCommandTest(c *gin.Context) {
	chatCommandIdParam := c.Param("chatCommandId")
	chatCommandId, err := strconv.ParseUint(chatCommandIdParam, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, NewErrorResponse(err))
		return
	}

	chatCommand, err := h.chatCommands.ChatCommandFindById(uint(chatCommandId))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, NewErrorResponse(err))
		return
	}

	err = h.actionWorker.HandleActionByName(chatCommand.Action, make(map[string]any))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, NewErrorResponse(err))
		return
	}

	c.Status(http.StatusOK)
}
