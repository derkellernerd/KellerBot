package handler

import (
	"fmt"
	"net/http"
	"path"
	"strconv"

	"github.com/derkellernerd/dori/core"
	"github.com/derkellernerd/dori/model"
	"github.com/derkellernerd/dori/repository"
	"github.com/gin-gonic/gin"
)

type Alert struct {
	env       *core.Environment
	alertRepo *repository.Alert
}

func NewAlert(env *core.Environment, alertRepo *repository.Alert) *Alert {
	return &Alert{
		env:       env,
		alertRepo: alertRepo,
	}
}

func (h *Alert) AlertGetAll(c *gin.Context) {
	alerts, err := h.alertRepo.AlertFindAll()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse(alerts))
}

func (h *Alert) AlertCreate(c *gin.Context) {
	var alertCreateRequest model.AlertCreateRequest

	err := c.BindJSON(&alertCreateRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, NewErrorResponse(err))
		return
	}

	alert := model.Alert{
		Name: alertCreateRequest.Name,
		Type: alertCreateRequest.Type,
	}

	if alertCreateRequest.Data != nil {
		alert.SetData(alertCreateRequest.Data)
	}

	err = h.alertRepo.AlertInsert(&alert)

	c.JSON(http.StatusCreated, NewSuccessResponse(alert))
}

func (h *Alert) AlertUploadFile(c *gin.Context) {
	alertIdParam := c.Param("alertId")
	alertId, err := strconv.ParseUint(alertIdParam, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, NewErrorResponse(err))
		return
	}

	alert, err := h.alertRepo.AlertFindById(uint(alertId))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, NewErrorResponse(err))
		return
	}

	file, _ := c.FormFile("file")
	destinationFilePath := path.Join("data", "alerts", fmt.Sprintf("%d", alert.ID))
	c.SaveUploadedFile(file, destinationFilePath)

	err = nil
	switch alert.Type {
	case model.ALERT_TYPE_GIF:
		alertGif := model.AlertTypeGif{
			GifPath: destinationFilePath,
		}
		alert.SetData(alertGif)
		break
	case model.ALERT_TYPE_GIF_SOUND:
	case model.ALERT_TYPE_SOUND:
	case model.ALERT_TYPE_VIDEO:
	default:
		err = fmt.Errorf("not implemented right now")
	}

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, NewErrorResponse(err))
		return
	}

	err = h.alertRepo.AlertUpdate(&alert)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, alert)
}

func (h *Alert) AlertGetFile(c *gin.Context) {
	alertIdParam := c.Param("alertId")
	alertId, err := strconv.ParseUint(alertIdParam, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, NewErrorResponse(err))
		return
	}

	slot := c.Query("slot")

	alert, err := h.alertRepo.AlertFindById(uint(alertId))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, NewErrorResponse(err))
		return
	}

	filePath := ""

	switch alert.Type {
	case model.ALERT_TYPE_SOUND:
		alertSound, err := alert.GetDataSound()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, NewErrorResponse(err))
			return
		}
		filePath = alertSound.SoundPath
		break
	case model.ALERT_TYPE_VIDEO:
		alertVideo, err := alert.GetDataVideo()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, NewErrorResponse(err))
			return
		}
		filePath = fmt.Sprintf("./data/alerts/%s", alertVideo.VideoPath)
		break
	case model.ALERT_TYPE_GIF:
		alertGif, err := alert.GetDataGif()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, NewErrorResponse(err))
			return
		}
		filePath = alertGif.GifPath
	case model.ALERT_TYPE_GIF_SOUND:
		alertGifSound, err := alert.GetDataGifSound()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, NewErrorResponse(err))
			return
		}

		switch slot {
		case "gif":
			filePath = fmt.Sprintf("./data/alerts/%s", alertGifSound.GifPath)
			break
		case "sound":
			filePath = fmt.Sprintf("./data/alerts/%s", alertGifSound.SoundPath)
			break
		}
	}

	if filePath != "" {
		c.File(filePath)
		return
	}

	c.Status(http.StatusNotFound)
}
