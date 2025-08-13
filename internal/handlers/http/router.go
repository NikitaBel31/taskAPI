package http

import (
	"net/http"
	"taskapi/internal/logger"
	"taskapi/internal/usecase"
)

type Router struct {
	svc usecase.TaskService
	log logger.Logger
}

func NewRouter(svc usecase.TaskService, log logger.Logger) *Router {
	return &Router{svc: svc, log: log}
}

func (rt *Router) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/tasks", rt.tasksCollection)
	mux.HandleFunc("/tasks/", rt.Get)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	//return requestIDMiddleware(mux)
	return requestIDMiddleware(rt.loggingMiddleware(mux))
}
