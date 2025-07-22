package task

import (
	"github.com/Gustcat/archiver_170725/internal/config"
	"github.com/Gustcat/archiver_170725/internal/repository"
	"github.com/Gustcat/archiver_170725/internal/service"
	"net/http"
)

type serv struct {
	taskRepo repository.TaskRepository
	config   *config.HTTPServer
	client   *http.Client
}

func NewServ(taskRepo repository.TaskRepository, config *config.HTTPServer, client *http.Client) service.TaskService {
	return &serv{
		taskRepo: taskRepo,
		config:   config,
		client:   client,
	}
}
