package task

import (
	"context"
	"github.com/Gustcat/archiver_170725/internal/model"
)

func (s *serv) Get(ctx context.Context, id string) (*model.TaskResult, error) {
	task, err := s.taskRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	taskResult := &model.TaskResult{
		Status: task.Status,
	}
	if task.Status == model.StatusDone {
		taskResult.ArchiveLink = &task.ArchiveLink
	}

	return taskResult, nil
}
