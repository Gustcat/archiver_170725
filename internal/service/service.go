package service

import (
	"context"
	"github.com/Gustcat/archiver_170725/internal/model"
)

type TaskService interface {
	Create(ctx context.Context) (string, error)
	Update(ctx context.Context, id string, source model.FileUrl) error
	Get(ctx context.Context, id string) (*model.TaskResult, error)
}
