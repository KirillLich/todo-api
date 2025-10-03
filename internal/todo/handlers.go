package todo

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	e "github.com/KirillLich/todoapi/internal/errors"
	"github.com/KirillLich/todoapi/internal/parse"
)

type TodoService interface {
	Create(ctx context.Context, todo Todo) error
	GetAll(ctx context.Context) ([]Todo, error)
	GetById(ctx context.Context, id int) (Todo, error)
	Update(ctx context.Context, todo Todo) error
	Delete(ctx context.Context, id int) error
}

type Handler struct {
	service TodoService
	logger  *slog.Logger
}

func NewHandler(s *Service, l *slog.Logger) *Handler {
	return &Handler{service: s, logger: l}
}

func (h *Handler) Todos(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.post(w, r)
	case http.MethodGet:
		h.getAll(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) TodosId(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.get(w, r)
	case http.MethodPut:
		h.put(w, r)
	case http.MethodDelete:
		h.delete(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) post(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	maxBytes := int64(1 << 20)
	var newTodo Todo
	dec := json.NewDecoder(io.LimitReader(r.Body, maxBytes))

	if err := dec.Decode(&newTodo); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	//context.WithTimeout(ctx, time.Second*5)
	if err := h.service.Create(ctx, newTodo); err != nil {
		code, msg := e.MapErrorHttpStatus(err)
		http.Error(w, msg, code)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) getAll(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := r.Context()
	//TODO: put repeated code into one func
	var buf bytes.Buffer
	content, err := h.service.GetAll(ctx)
	if err != nil {
		code, msg := e.MapErrorHttpStatus(err)
		http.Error(w, msg, code)
		return
	}
	if err := json.NewEncoder(&buf).Encode(content); err != nil {
		http.Error(w, "failed to encode json", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	buf.WriteTo(w)

}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	id, err := parse.ParseId(r.URL.Path, "todos")
	if err != nil {
		code, msg := e.MapErrorHttpStatus(err)
		http.Error(w, msg, code)
		return
	}
	ctx := r.Context()
	//TODO: put repeated code into one func
	var buf bytes.Buffer
	content, err := h.service.GetById(ctx, id)
	if err != nil {
		code, msg := e.MapErrorHttpStatus(err)
		http.Error(w, msg, code)
		return
	}
	if err := json.NewEncoder(&buf).Encode(content); err != nil {
		http.Error(w, "failed to encode json", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	buf.WriteTo(w)
}

func (h *Handler) put(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	id, err := parse.ParseId(r.URL.Path, "todos")
	if err != nil {
		code, msg := e.MapErrorHttpStatus(err)
		http.Error(w, msg, code)
		return
	}
	maxBytes := int64(1 << 20)
	var newTodo Todo
	dec := json.NewDecoder(io.LimitReader(r.Body, maxBytes))
	if err := dec.Decode(&newTodo); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if id != newTodo.Id {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	//context.WithTimeout(ctx, time.Second*5)
	if err := h.service.Update(ctx, newTodo); err != nil {
		code, msg := e.MapErrorHttpStatus(err)
		http.Error(w, msg, code)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	id, err := parse.ParseId(r.URL.Path, "todos")
	if err != nil {
		code, msg := e.MapErrorHttpStatus(err)
		http.Error(w, msg, code)
		return
	}
	ctx := r.Context()
	if err := h.service.Delete(ctx, id); err != nil {
		code, msg := e.MapErrorHttpStatus(err)
		http.Error(w, msg, code)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
