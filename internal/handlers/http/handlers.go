package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"taskapi/internal/domain"
	"taskapi/internal/dto"
	"taskapi/internal/usecase"
)

func (rt *Router) tasksCollection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		rt.GetList(w, r)
	case http.MethodPost:
		rt.Create(w, r)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (rt *Router) Get(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/tasks/"), "/")
	id := parts[0]
	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}
	reqID := requestIDFromCtx(r.Context())

	t, err := rt.svc.Get(r.Context(), reqID, id)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, usecase.ErrNotFound) {
			status = http.StatusNotFound
		}
		writeJSON(w, status, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, t)
}

func (rt *Router) GetList(w http.ResponseWriter, r *http.Request) {
	reqID := requestIDFromCtx(r.Context())
	var st *domain.Status
	if q := strings.TrimSpace(r.URL.Query().Get("status")); q != "" {
		s := domain.Status(q)
		if s != domain.StatusTodo && s != domain.StatusInProgress && s != domain.StatusDone {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid status"})
			return
		}
		st = &s
	}
	list, err := rt.svc.List(r.Context(), reqID, st)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, list)
}

func (rt *Router) Create(w http.ResponseWriter, r *http.Request) {
	defer func() {
		_ = r.Body.Close()
	}()
	var in dto.CreateInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON body"})
		return
	}
	reqID := requestIDFromCtx(r.Context())
	t, err := rt.svc.Create(r.Context(), reqID, in)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, usecase.ErrBadRequest) {
			status = http.StatusBadRequest
		}
		writeJSON(w, status, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, t)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
