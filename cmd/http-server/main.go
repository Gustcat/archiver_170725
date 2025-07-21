package main

import (
	"context"
	"github.com/Gustcat/archiver_170725/internal/api/handler"
	"github.com/Gustcat/archiver_170725/internal/config"
	"github.com/Gustcat/archiver_170725/internal/logger"
	taskRepo "github.com/Gustcat/archiver_170725/internal/repository/task"
	taskService "github.com/Gustcat/archiver_170725/internal/service/task"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log := logger.SetupLogger(slog.LevelDebug)

	ctx := context.Background()

	err := godotenv.Load(".env")
	if err != nil {
		log.Warn("doesn't load env file: %s", slog.String("error", err.Error()))
	}

	conf, err := config.New()
	if err != nil {
		log.Error("doesn't set config: %s", slog.String("error", err.Error()))
		os.Exit(1)
	}

	repo := taskRepo.NewRepo()

	service := taskService.NewServ(repo, conf)
	h := handler.NewHandler(service)

	log.Debug("Try to setup router")
	router := gin.New()

	router.Use(gin.Recovery())

	r := router.Group(config.TaskGroupUrl)
	{
		r.POST("/", h.Create)
		r.GET("/:id", h.Get)
		r.PATCH("/:id", h.Update)
	}

	srv := &http.Server{
		Addr:         conf.Address,
		Handler:      router,
		ReadTimeout:  conf.Timeout,
		WriteTimeout: conf.Timeout,
		IdleTimeout:  conf.IdleTimeout,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error("failed to start http server", slog.String("error", err.Error()))
		}
	}()

	log.Info("Server started", slog.String("address", conf.Address))

	<-quit

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("failed to shutdown http server", slog.String("error", err.Error()))
		return
	}

	log.Info("Server stopped", slog.String("address", conf.Address))
}
