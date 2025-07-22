package task

import (
	"context"
	"errors"
	"fmt"
	"github.com/Gustcat/archiver_170725/internal/logger"
	"github.com/Gustcat/archiver_170725/internal/model"
	"log/slog"
	"os"
)

var (
	ErrArchiveNotFound = errors.New("archive not found")
)

func (s *serv) GetArchive(ctx context.Context, id string) (*os.File, error) {
	const op = "service.task.GetArchive"
	log := logger.LogFromContextAddOP(ctx, op)

	task, err := s.taskRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if task.Status == model.StatusFailed {
		return nil, fmt.Errorf("archiving failed")
	}
	if task.Status != model.StatusDone || !task.ArchiveClosed {
		log.Error("", slog.String("error", ErrArchiveNotFound.Error()))
		return nil, ErrArchiveNotFound
	}

	return task.Archive, nil
}
