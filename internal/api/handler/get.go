package handler

import (
	"errors"
	"github.com/Gustcat/archiver_170725/internal/logger"
	"github.com/Gustcat/archiver_170725/internal/model"
	"github.com/Gustcat/archiver_170725/internal/repository/task"
	"github.com/Gustcat/archiver_170725/internal/response"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

func (h *Handler) Get(c *gin.Context) {
	const op = "handler.task.Get"

	ctx := c.Request.Context()
	log := logger.LogFromContextAddOP(ctx, op)

	log.Debug("Receive request for get task")

	var taskId model.TaskId
	err := c.ShouldBindUri(&taskId)
	if err != nil {
		log.Error("invalid url-parameter: id")
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error("invalid url-parameter: id"))
		return
	}

	taskResult, err := h.service.Get(ctx, taskId.ID)
	if errors.Is(err, task.ErrTaskNotFound) {
		log.Error("", slog.String("error", err.Error()), slog.Int64("id", taskId.ID))
		c.AbortWithStatusJSON(http.StatusNotFound, response.Error(err.Error()))
		return
	}
	if err != nil {
		log.Error("Failed to get task", slog.String("error", err.Error()), slog.Int64("id", taskId.ID))
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.Error("Failed to get task"))
		return
	}

	c.JSON(http.StatusOK, response.OK(taskResult))
}
