package todo

import (
	"context"
	"database/sql"
	"fmt"

	e "github.com/KirillLich/todoapi/internal/errors"
)

type Repo interface {
	Create(ctx context.Context, todo Todo) error
	GetAll(ctx context.Context) ([]Todo, error)
	GetById(ctx context.Context, id int) (Todo, error)
	Update(ctx context.Context, todo Todo) error
	Delete(ctx context.Context, id int) error
}

type Service struct {
	Repo Repo
}

func NewService(repo *Repository) *Service {
	return &Service{Repo: repo}
}

func (s *Service) Create(ctx context.Context, todo Todo) error {
	if todo.Title == "" {
		return fmt.Errorf("service.Create: %w", e.ErrEmptyTitle)
	}
	return s.Repo.Create(ctx, todo)
}

func (s *Service) GetAll(ctx context.Context) ([]Todo, error) {
	return s.Repo.GetAll(ctx)
}

func (s *Service) GetById(ctx context.Context, id int) (Todo, error) {
	T, err := s.Repo.GetById(ctx, id)
	if err == sql.ErrNoRows {
		return Todo{}, e.ErrNotFound
	}
	return T, err
}

func (s *Service) Update(ctx context.Context, todo Todo) error {
	if todo.Title == "" {
		return fmt.Errorf("service.Update: %w", e.ErrEmptyTitle)
	}
	err := s.Repo.Update(ctx, todo)
	if err == sql.ErrNoRows {
		return e.ErrNotFound
	}
	return err
}

func (s *Service) Delete(ctx context.Context, id int) error {
	return s.Repo.Delete(ctx, id)
}
