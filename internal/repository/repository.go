package repository

import (
	"context"
	"github.com/Gustcat/archiver_170725/internal/model"
)

type TaskRepository interface {
	Create(ctx context.Context, id string, task *model.Task) error
	Update(ctx context.Context, id, sourceLink string) (*model.Task, int, error)
	Get(ctx context.Context, id string) (*model.Task, error)
	CountByStatus(ctx context.Context, status model.Status) int
	ProcessDownloadSignal(ctx context.Context, task *model.Task, err error)
	DownloadSource(ctx context.Context, task *model.Task, fileName string, data []byte) error
}
