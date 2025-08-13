package app

import (
	"os"
	"taskapi/internal/config"
	httpHandler "taskapi/internal/handlers/http"
	"taskapi/internal/logger"
	"taskapi/internal/repository/memory"
	"taskapi/internal/usecase"
)

type Container struct {
	Config *config.Config
	Logger *logger.Async
	Repo   *memory.Repo
	Svc    usecase.TaskService
	Router httpHandler.Router
}

func NewContainer() *Container {
	cfg := config.Load()
	log := logger.NewAsync(cfg.LogBuffer, os.Stdout)
	repo := memory.New()
	svc := usecase.NewService(repo, log)
	router := httpHandler.NewRouter(svc, log)

	return &Container{
		Config: cfg,
		Logger: log,
		Repo:   repo,
		Svc:    svc,
		Router: *router,
	}
}
