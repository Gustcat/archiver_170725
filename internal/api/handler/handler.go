package handler

import "github.com/Gustcat/archiver_170725/internal/service"

type Handler struct {
	service service.TaskService
}

func NewHandler(service service.TaskService) *Handler {
	return &Handler{service: service}
}
