package task

import (
	"context"
	"errors"
	"github.com/Gustcat/archiver_170725/internal/model"
	"github.com/Gustcat/archiver_170725/internal/repository"
	"sync"
)

var ErrTaskNotFound = errors.New("task not found")

type Repo struct {
	mu     sync.RWMutex
	tasks  map[int64]*model.Task
	lastID int64
}

func NewRepo() repository.TaskRepository {
	return &Repo{
		tasks: make(map[int64]*model.Task),
	}
}

func (r *Repo) Create(ctx context.Context) (int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.lastID++
	r.tasks[r.lastID] = &model.Task{
		Sources: make([]model.FileUrl, 0, 3),
		Status:  model.StatusNew,
	}

	return 0, nil
}

func (r *Repo) Update(ctx context.Context, link string) error {
	return nil
}

func (r *Repo) Get(ctx context.Context, id int64) (*model.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	task, ok := r.tasks[id]
	if !ok {
		return nil, ErrTaskNotFound
	}

	return task, nil
}

func (r *Repo) CountByStatus(ctx context.Context, status model.Status) (int, error) {
	return 0, nil
}
