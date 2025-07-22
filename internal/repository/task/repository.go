package task

import (
	"context"
	"errors"
	"github.com/Gustcat/archiver_170725/internal/logger"
	"github.com/Gustcat/archiver_170725/internal/model"
	"github.com/Gustcat/archiver_170725/internal/repository"
	"log/slog"
	"sync"
)

const (
	SourcesLimit = 3
)

var (
	ErrOverSourcesLimit   = errors.New("over task limit")
	ErrTaskNotFound       = errors.New("task not found")
	ErrTaskAlreadyExists  = errors.New("task already exists")
	ErrSourceAlreadyAdded = errors.New("source already added")
)

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

	return nil
}

func (r *Repo) Update(ctx context.Context, id, sourceLink string) (*model.Task, int, error) {
	r.mu.RLock()
	task, ok := r.tasks[id]
	if !ok {
		r.mu.RUnlock()
		return nil, 0, ErrTaskNotFound
	}

	sourcesCount := len(task.Sources)
	if sourcesCount >= SourcesLimit {
		r.mu.RUnlock()
		return nil, 0, ErrOverSourcesLimit
	}
	r.mu.RUnlock()

	task.Mu.RLock()
	defer task.Mu.RUnlock()

	task.Status = model.StatusInProgress

	for _, sL := range task.Sources {
		if sL == sourceLink {
			return nil, 0, ErrSourceAlreadyAdded
		}
	}

	task.Sources = append(task.Sources, sourceLink)

	return task, len(task.Sources), nil
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

func (r *Repo) CountByStatus(ctx context.Context, status model.Status) int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var count int

	for _, task := range r.tasks {
		if task.Status == status {
			count++
		}
	}

	return count
}

func (r *Repo) ProcessDownloadSignal(ctx context.Context, task *model.Task, err error) {
	const op = "repo.task.ProcessDownloadSignal"

	log := logger.LogFromContextAddOP(ctx, op)

	task.Mu.RLock()
	defer task.Mu.RUnlock()

	task.DownloadDown = append(task.DownloadDown, err)

	if len(task.DownloadDown) == SourcesLimit {
		if !task.ZipWriterClosed {
			err := task.ZipWriter.Close()
			if err != nil {
				task.Status = model.StatusFailed
				log.Error("zip writer close failed", slog.String("error", err.Error()))
			}
			task.ZipWriterClosed = true
		}

		if !task.ArchiveClosed {
			err := task.Archive.Close()
			if err != nil {
				task.Status = model.StatusFailed
				log.Error("zip writer close failed", slog.String("error", err.Error()))
			}
			task.ArchiveClosed = true
		}

		task.Status = model.StatusDone
	}
}

func (r *Repo) DownloadSource(ctx context.Context, task *model.Task, fileName string, data []byte) error {
	const op = "repo.task.ProcessDownloadSignal"

	log := logger.LogFromContextAddOP(ctx, op)

	task.Mu.RLock()
	defer task.Mu.RUnlock()

	zipFile, err := task.ZipWriter.Create(fileName)
	if err != nil {
		log.Error("failed to create file in archive", slog.String("error", err.Error()))
		return err
	}

	if _, err := zipFile.Write(data); err != nil {
		log.Error("failed write to archive", slog.String("error", err.Error()))
		return err
	}

	return nil
}
