package task

import (
	"context"
	"errors"
	"fmt"
	"github.com/Gustcat/archiver_170725/internal/model"
)

var ErrOverTaskInProgressLimit = errors.New("over task limit")

const TaskInProgressLimit = 3

func (s *serv) Create(ctx context.Context) (int64, error) {
	const op = "service.task.Create"

	tasksInProgress, err := s.taskRepo.CountByStatus(ctx, model.StatusInProgress)
	if err != nil {
		return 0, err
	}
	if tasksInProgress >= TaskInProgressLimit {
		return 0, fmt.Errorf("%s: %w", op, ErrOverTaskInProgressLimit)
	}

	return s.taskRepo.Create(ctx)
}
