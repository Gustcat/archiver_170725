package task

import (
	"github.com/Gustcat/archiver_170725/internal/repository"
	"github.com/Gustcat/archiver_170725/internal/service"
)

type serv struct {
	taskRepo repository.TaskRepository
}

func NewServ(taskRepo repository.TaskRepository) service.TaskService {
	return &serv{
		taskRepo: taskRepo,
	}
}
