package task

import (
	"context"
	"github.com/Gustcat/archiver_170725/internal/model"
	"strings"
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

		var msgs []string
		for number, res := range task.DownloadDown {
			if res != nil {
				msg := task.Sources[number] + ": " + res.Error()
				msgs = append(msgs, msg)
			}
		}

		if len(msgs) > 0 {
			esmsg := strings.Join(msgs, "\n")
			taskResult.ErrorSourcesMessage = &esmsg
		}
	}

	return taskResult, nil
}
