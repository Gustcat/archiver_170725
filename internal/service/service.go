package service

import (
	"context"
	"github.com/Gustcat/archiver_170725/internal/model"
	"os"
)

type TaskService interface {
	Create(ctx context.Context) (string, error)
	Update(ctx context.Context, id, source string) error
	Get(ctx context.Context, id string) (*model.TaskResult, error)
	GetArchive(ctx context.Context, id string) (*os.File, error)
}
