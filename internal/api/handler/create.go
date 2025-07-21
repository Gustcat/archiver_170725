package handler

import (
	"errors"
	"github.com/Gustcat/archiver_170725/internal/logger"
	"github.com/Gustcat/archiver_170725/internal/model"
	"github.com/Gustcat/archiver_170725/internal/response"
	"github.com/Gustcat/archiver_170725/internal/service/task"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

func (h *Handler) Create(c *gin.Context) {
	const op = "handler.task.Create"

	ctx := c.Request.Context()
	log := logger.LogFromContextAddOP(ctx, op)

	log.Debug("Receive request for create task")

	id, err := h.service.Create(ctx)
	if errors.Is(err, task.ErrOverTaskInProgressLimit) {
		log.Error("", slog.String("error", err.Error()))
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error(err.Error()))
	}
	if err != nil {
		log.Error("Failed to create task", slog.String("error", err.Error()))
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.Error("Failed to create task"))
		return
	}

	log.Info("Task created", slog.String("id", id))

	resp := &model.TaskId{ID: id}
	c.JSON(http.StatusCreated, response.OK(resp))
}
