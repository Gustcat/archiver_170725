package task

import (
	"archive/zip"
	"context"
	"errors"
	"fmt"
	"github.com/Gustcat/archiver_170725/internal/config"
	"github.com/Gustcat/archiver_170725/internal/model"
	"github.com/google/uuid"
	"net/url"
	"os"
	"path"
)

const ArchiveDirName = "archives"

var ErrOverTaskInProgressLimit = errors.New("over task limit")

const TaskInProgressLimit = 3

func (s *serv) Create(ctx context.Context) (string, error) {
	const op = "service.task.Create"

	tasksInProgress, err := s.taskRepo.CountByStatus(ctx, model.StatusInProgress)
	if err != nil {
		return "", err
	}
	if tasksInProgress >= TaskInProgressLimit {
		return "", fmt.Errorf("%s: %w", op, ErrOverTaskInProgressLimit)
	}

	id := uuid.New().String()

	if _, err := os.Stat(ArchiveDirName); os.IsNotExist(err) {
		err := os.Mkdir(ArchiveDirName, 0755)
		if err != nil {
			return "", fmt.Errorf("%s: failed to create archive dir: %w", op, err)
		}
	}

	filePath := path.Join(ArchiveDirName, id) + ".zip"

	archive, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("%s: failed to create zip file: %w", op, err)
	}

	zipWriter := zip.NewWriter(archive)

	task := &model.Task{
		Sources:      make([]model.FileUrl, 0, 3),
		Status:       model.StatusNew,
		ZipWriter:    zipWriter,
		ArchiveLink:  s.generateArchiveLink(id),
		DownloadDown: make(chan struct{}, 3),
	}

	err = s.taskRepo.Create(ctx, id, task)
	if err != nil {
		return "", fmt.Errorf("%s: failed to create task: %w", op, err)
	}

	return id, nil
}

func (s *serv) generateArchiveLink(id string) model.FileUrl {
	base := url.URL{
		Scheme: "http",
		Host:   s.config.Address,
		Path:   path.Join(config.TaskGroupUrl, id),
	}
	return model.FileUrl(base.String())
}
