package handler

import (
	"errors"
	"github.com/Gustcat/archiver_170725/internal/logger"
	"github.com/Gustcat/archiver_170725/internal/model"
	"github.com/Gustcat/archiver_170725/internal/response"
	taskService "github.com/Gustcat/archiver_170725/internal/service/task"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

func (h *Handler) Update(c *gin.Context) {
	const op = "handler.task.Get"

	ctx := c.Request.Context()
	log := logger.LogFromContextAddOP(ctx, op)

	log.Debug("Receive request for add sources")

	var taskId model.TaskId
	err := c.ShouldBindUri(&taskId)
	if err != nil {
		log.Error("invalid url-parameter: id")
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error("invalid url-parameter: id"))
		return
	}

	var sourceRequest model.SourceRequest
	err = c.ShouldBindJSON(&sourceRequest)
	if err != nil {
		log.Error("invalid json body: ")
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error("invalid json body: "))
		return
	}

	err = h.service.Update(ctx, taskId.ID, sourceRequest.Source)
	if errors.Is(err, taskService.ErrOverSourcesLimit) {
		log.Error("", slog.String("error", err.Error()), slog.Int64("id", taskId.ID))
		c.AbortWithStatusJSON(http.StatusNotFound, response.Error(err.Error()))
		return
	}

	if err != nil {
		log.Error("Failed to add sources", slog.String("error", err.Error()), slog.Int64("id", taskId.ID))
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.Error("Failed to get task"))
		return
	}

	c.JSON(http.StatusOK, response.OK(nil))
}
