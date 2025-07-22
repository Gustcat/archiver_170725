package handler

import (
	"errors"
	"fmt"
	"github.com/Gustcat/archiver_170725/internal/logger"
	"github.com/Gustcat/archiver_170725/internal/model"
	taskRepo "github.com/Gustcat/archiver_170725/internal/repository/task"
	"github.com/Gustcat/archiver_170725/internal/response"
	taskService "github.com/Gustcat/archiver_170725/internal/service/task"
	"github.com/gin-gonic/gin"
	"io"
	"log/slog"
	"net/http"
	"path/filepath"
)

func (h *Handler) GetArchive(c *gin.Context) {
	const op = "handler.task.GetArchive"

	ctx := c.Request.Context()
	log := logger.LogFromContextAddOP(ctx, op)

	log.Debug("Receive archive")

	var taskId model.TaskId
	err := c.ShouldBindUri(&taskId)
	if err != nil {
		log.Error("invalid url-parameter: id")
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error("invalid url-parameter: id"))
		return
	}

	archive, err := h.service.GetArchive(ctx, taskId.ID)
	if errors.Is(err, taskRepo.ErrTaskNotFound) || errors.Is(err, taskService.ErrArchiveNotFound) {
		log.Error("", slog.String("error", err.Error()), slog.String("id", taskId.ID))
		c.AbortWithStatusJSON(http.StatusNotFound, response.Error(err.Error()))
		return
	}
	if err != nil {
		log.Error("Failed to get task", slog.String("error", err.Error()), slog.String("id", taskId.ID))
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.Error("Failed to get task"))
		return
	}

	if _, err := archive.Seek(0, io.SeekStart); err != nil {
		log.Error("Failed to move pointer to the beginning of archive",
			slog.String("error", err.Error()),
			slog.String("id", taskId.ID))
		c.JSON(http.StatusInternalServerError, response.Error("failed to prepare archive for download"))
		return
	}

	fileName := filepath.Base(archive.Name())

	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, fileName))
	c.Header("Content-Type", "application/zip")

	if _, err := io.Copy(c.Writer, archive); err != nil {
		log.Error("failed archive transfer",
			slog.String("error", err.Error()),
			slog.String("id", taskId.ID))
		return
	}
}
