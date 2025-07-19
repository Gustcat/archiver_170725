package repository

import (
	"context"
	"github.com/Gustcat/archiver_170725/internal/model"
)

type TaskRepository interface {
	Create(ctx context.Context) (int64, error)
	Update(ctx context.Context, link string) error
	Get(ctx context.Context, id int64) (*model.Task, error)
	CountByStatus(ctx context.Context, status model.Status) (int, error)
}
