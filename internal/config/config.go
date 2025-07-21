package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v10"
)

const TaskGroupUrl = "/api/v1/tasks"

type HTTPServer struct {
	Host        string `env:"HTTP_HOST" envDefault:"localhost"`
	Port        string `env:"HTTP_PORT" envDefault:"8080"`
	Address     string
	Timeout     time.Duration `env:"HTTP_TIMEOUT" envDefault:"5s"`
	IdleTimeout time.Duration `env:"HTTP_IDLE_TIMEOUT" envDefault:"60s"`
	User        string        `env:"HTTP_USER" envDefault:"user"`
	Password    string        `env:"HTTP_PASSWORD" envDefault:"password"`
}

func New() (*HTTPServer, error) {
	cfg := &HTTPServer{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("loading config from env is failed: %w", err)
	}
	buildHTTPAddress(cfg)

	return cfg, nil
}

func buildHTTPAddress(httpserver *HTTPServer) {
	httpserver.Address = fmt.Sprintf("%s:%s", httpserver.Host, httpserver.Port)
}
