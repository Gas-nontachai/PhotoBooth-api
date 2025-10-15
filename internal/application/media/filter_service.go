package media

import (
	"context"

	domain "go-ddd-clean/internal/domain/media"

	"github.com/google/uuid"
)

type FilterService struct {
	repo domain.FilterRepository
}

func NewFilterService(repo domain.FilterRepository) *FilterService {
	return &FilterService{repo: repo}
}

type CreateFilterInput struct {
	Name   string
	Effect map[string]any
	Active bool
}

type UpdateFilterInput struct {
	ID     string
	Name   string
	Effect map[string]any
	Active bool
}

func (s *FilterService) Create(ctx context.Context, input CreateFilterInput) (*domain.Filter, error) {
	entity := &domain.Filter{
		ID:     uuid.NewString(),
		Name:   input.Name,
		Effect: input.Effect,
		Active: input.Active,
	}
	if err := s.repo.Create(ctx, entity); err != nil {
		return nil, err
	}
	return entity, nil
}

func (s *FilterService) Update(ctx context.Context, input UpdateFilterInput) (*domain.Filter, error) {
	entity, err := s.repo.GetByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}
	entity.Name = input.Name
	entity.Effect = input.Effect
	entity.Active = input.Active
	if err := s.repo.Update(ctx, entity); err != nil {
		return nil, err
	}
	return entity, nil
}

func (s *FilterService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *FilterService) Get(ctx context.Context, id string) (*domain.Filter, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *FilterService) List(ctx context.Context, onlyActive bool) ([]domain.Filter, error) {
	return s.repo.List(ctx, onlyActive)
}
