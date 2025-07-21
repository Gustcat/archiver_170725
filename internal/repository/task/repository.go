package task

import (
	"context"
	"errors"
	"fmt"
	"github.com/Gustcat/archiver_170725/internal/model"
	"github.com/Gustcat/archiver_170725/internal/repository"
	"sync"
)

var ErrTaskNotFound = errors.New("task not found")
var ErrTaskAlreadyExists = errors.New("task already exists")

type Repo struct {
	mu    sync.RWMutex
	tasks map[string]*model.Task
}

func NewRepo() repository.TaskRepository {
	return &Repo{
		tasks: make(map[string]*model.Task),
	}
}

func (r *Repo) Create(ctx context.Context, id string, task *model.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.tasks[id]; ok {
		return ErrTaskAlreadyExists
	}

	r.tasks[id] = task
	for key, value := range r.tasks {
		fmt.Printf("%+v", key)
		fmt.Println()
		fmt.Printf("%+v", value)
	}
	fmt.Printf("%+v", r.tasks)

	return nil
}

func (r *Repo) Update(ctx context.Context, link string) error {
	return nil
}

func (r *Repo) Get(ctx context.Context, id string) (*model.Task, error) {
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
