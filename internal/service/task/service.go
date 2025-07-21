package task

import (
	"github.com/Gustcat/archiver_170725/internal/config"
	"github.com/Gustcat/archiver_170725/internal/repository"
	"github.com/Gustcat/archiver_170725/internal/service"
)

type serv struct {
	taskRepo repository.TaskRepository
	config   *config.HTTPServer
}

func NewServ(taskRepo repository.TaskRepository, config *config.HTTPServer) service.TaskService {
	return &serv{
		taskRepo: taskRepo,
		config:   config,
	}
}
