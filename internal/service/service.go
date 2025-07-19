package service

import (
	"context"
	"github.com/Gustcat/archiver_170725/internal/model"
)

type TaskService interface {
	Create(ctx context.Context) (int64, error)
	Update(ctx context.Context, id int64, source model.FileUrl) error
	Get(ctx context.Context, id int64) (*model.TaskResult, error)
}
