package handler

import (
	"fmt"
	"net/http"
	"path"
	"strconv"

	"github.com/derkellernerd/kellerbot/core"
	"github.com/derkellernerd/kellerbot/model"
	"github.com/derkellernerd/kellerbot/repository"
	"github.com/gin-gonic/gin"
)

type Action struct {
	env     *core.Environment
	actions *repository.Action
}

func NewAction(env *core.Environment, actionRepo *repository.Action) *Action {
	return &Action{
		env:     env,
		actions: actionRepo,
	}
}

func (h *Action) ActionGetAll(c *gin.Context) {
	actions, err := h.actions.ActionFindAll()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse(actions))
}

func (h *Action) ActionGetById(c *gin.Context) {
	actionIdParam := c.Param("actionId")
	actionId, err := strconv.ParseUint(actionIdParam, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, NewErrorResponse(err))
		return
	}

	action, err := h.actions.ActionFindById(uint(actionId))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse(action))
}

func (h *Action) ActionCreate(c *gin.Context) {
	var actionCreateRequest model.ActionCreateRequest

	err := c.BindJSON(&actionCreateRequest)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, NewErrorResponse(err))
		return
	}

	action := model.Action{
		ActionName: actionCreateRequest.ActionName,
		ActionType: actionCreateRequest.ActionType,
	}

	if actionCreateRequest.Data != nil {
		err = action.SetData(actionCreateRequest.Data)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, NewErrorResponse(err))
			return
		}
	}

	err = h.actions.ActionInsert(&action)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, NewSuccessResponse(action))
}

func (h *Action) ActionUpdate(c *gin.Context) {
	actionIdParam := c.Param("actionId")
	actionId, err := strconv.ParseUint(actionIdParam, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, NewErrorResponse(err))
		return
	}

	var actionUpdateRequest model.ActionUpdateRequest
	err = c.BindJSON(&actionUpdateRequest)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, NewErrorResponse(err))
		return
	}

	action, err := h.actions.ActionFindById(uint(actionId))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, NewErrorResponse(err))
		return
	}

	err = action.SetData(actionUpdateRequest.Data)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, NewErrorResponse(err))
		return
	}

	switch action.ActionType {
	case model.ACTION_TYPE_COMPOSITION:
		composition, err := model.ActionGetData[model.ActionTypeComposition](&action)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, NewErrorResponse(err))
			return
		}

		newActions := []string{}
		for _, action := range composition.Actions {
			if action != "" {
				newActions = append(newActions, action)
			}
		}

		composition.Actions = newActions
		action.SetData(composition)
		break
	}

	err = h.actions.ActionUpdate(&action)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, NewSuccessResponse(action))
}

func (h *Action) ActionDelete(c *gin.Context) {
	actionIdParam := c.Param("actionId")
	actionId, err := strconv.ParseUint(actionIdParam, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, NewErrorResponse(err))
		return
	}

	action, err := h.actions.ActionFindById(uint(actionId))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, NewErrorResponse(err))
		return
	}

	err = h.actions.ActionDelete(&action)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, NewErrorResponse(err))
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Action) ActionUploadFile(c *gin.Context) {
	actionIdParam := c.Param("actionId")
	actionId, err := strconv.ParseUint(actionIdParam, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, NewErrorResponse(err))
		return
	}

	action, err := h.actions.ActionFindById(uint(actionId))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, NewErrorResponse(err))
		return
	}

	file, _ := c.FormFile("file")
	destinationFilePath := path.Join("data", "actions", fmt.Sprintf("%d", action.ID))
	c.SaveUploadedFile(file, destinationFilePath)

	err = nil
	switch action.ActionType {
	case model.ACTION_TYPE_GIF:
		actionData, err := model.ActionGetData[model.ActionTypeGif](&action)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, NewErrorResponse(err))
			return
		}
		actionData.FileName = destinationFilePath
		action.SetData(actionData)
	case model.ACTION_TYPE_SOUND:
		actionData, err := model.ActionGetData[model.ActionTypeSound](&action)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, NewErrorResponse(err))
			return
		}
		actionData.FileName = destinationFilePath
		action.SetData(actionData)
	default:
		panic(fmt.Sprintf("unexpected model.ActionType: %#v", action.ActionType))
	}
	err = h.actions.ActionUpdate(&action)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, action)
}

func (h *Action) ActionGetFile(c *gin.Context) {
	actionIdParam := c.Param("actionId")
	actionId, err := strconv.ParseUint(actionIdParam, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, NewErrorResponse(err))
		return
	}

	action, err := h.actions.ActionFindById(uint(actionId))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, NewErrorResponse(err))
		return
	}

	filePath := ""

	switch action.ActionType {
	case model.ACTION_TYPE_GIF:
		actionData, err := model.ActionGetData[model.ActionTypeGif](&action)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, NewErrorResponse(err))
			return
		}
		filePath = actionData.FileName
	case model.ACTION_TYPE_SOUND:
		actionData, err := model.ActionGetData[model.ActionTypeSound](&action)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, NewErrorResponse(err))
			return
		}
		filePath = actionData.FileName
	default:
		panic(fmt.Sprintf("unexpected model.ActionType: %#v", action.ActionType))
	}

	if filePath != "" {
		c.File(filePath)
		return
	}

	c.Status(http.StatusNotFound)
}
