package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/derkellernerd/dori/core"
	"github.com/derkellernerd/dori/model"
	"github.com/derkellernerd/dori/repository"
	"github.com/gin-gonic/gin"
)

type Command struct {
	env         *core.Environment
	commandRepo *repository.Command
}

func NewCommand(env *core.Environment, commandRepo *repository.Command) *Command {
	return &Command{
		env:         env,
		commandRepo: commandRepo,
	}
}

func (h *Command) CommandGetAll(c *gin.Context) {
	commands, err := h.commandRepo.CommandFindAll()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse(commands))
}

func (h *Command) CommandCreate(c *gin.Context) {
	var commandCreateRequest model.CommandCreateRequest

	err := c.BindJSON(&commandCreateRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, NewErrorResponse(err))
		return
	}

	if model.CommandIsBlacklisted(strings.ToLower(commandCreateRequest.Command)) {
		c.AbortWithStatusJSON(http.StatusBadRequest, NewErrorResponse(fmt.Errorf("command %s is blacklisted", commandCreateRequest.Command)))
		return
	}

	command := model.Command{
		Command:          commandCreateRequest.Command,
		Type:             commandCreateRequest.Type,
		TimeoutInSeconds: commandCreateRequest.TimeoutInSeconds,
	}

	command.SetData(commandCreateRequest.Data)

	err = h.commandRepo.CommandInsert(&command)

	c.JSON(http.StatusCreated, NewSuccessResponse(command))
}

func (h *Command) CommandUpdate(c *gin.Context) {
	commandIdParam := c.Param("commandId")
	commandId, err := strconv.ParseUint(commandIdParam, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, NewErrorResponse(err))
		return
	}

	var commandUpdateRequest model.CommandUpdateRequest
	err = c.BindJSON(&commandUpdateRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, NewErrorResponse(err))
		return
	}

	command, err := h.commandRepo.CommandFindById(uint(commandId))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, NewErrorResponse(err))
		return
	}

	command.SetData(commandUpdateRequest.Data)
	command.TimeoutInSeconds = commandUpdateRequest.TimeoutInSeconds

	err = h.commandRepo.CommandUpdate(&command)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse(command))
}

func (h *Command) CommandDelete(c *gin.Context) {
	commandIdParam := c.Param("commandId")
	commandId, err := strconv.ParseUint(commandIdParam, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, NewErrorResponse(err))
		return
	}

	command, err := h.commandRepo.CommandFindById(uint(commandId))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, NewErrorResponse(err))
		return
	}

	err = h.commandRepo.CommandDelete(&command)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, NewErrorResponse(err))
		return
	}

	c.Status(http.StatusNoContent)
}
