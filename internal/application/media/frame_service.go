package media

import (
	"context"

	domain "go-ddd-clean/internal/domain/media"

	"github.com/google/uuid"
)

type FrameService struct {
	repo domain.FrameRepository
}

func NewFrameService(repo domain.FrameRepository) *FrameService {
	return &FrameService{repo: repo}
}

type CreateFrameInput struct {
	Name    string
	Theme   *string
	FileURL string
	Active  bool
}

type UpdateFrameInput struct {
	ID      string
	Name    string
	Theme   *string
	FileURL string
	Active  bool
}

func (s *FrameService) Create(ctx context.Context, input CreateFrameInput) (*domain.Frame, error) {
	entity := &domain.Frame{
		ID:      uuid.NewString(),
		Name:    input.Name,
		Theme:   input.Theme,
		FileURL: input.FileURL,
		Active:  input.Active,
	}
	if err := s.repo.Create(ctx, entity); err != nil {
		return nil, err
	}
	return entity, nil
}

func (s *FrameService) Update(ctx context.Context, input UpdateFrameInput) (*domain.Frame, error) {
	entity, err := s.repo.GetByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}
	entity.Name = input.Name
	entity.Theme = input.Theme
	entity.FileURL = input.FileURL
	entity.Active = input.Active
	if err := s.repo.Update(ctx, entity); err != nil {
		return nil, err
	}
	return entity, nil
}

func (s *FrameService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *FrameService) Get(ctx context.Context, id string) (*domain.Frame, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *FrameService) List(ctx context.Context, onlyActive bool) ([]domain.Frame, error) {
	return s.repo.List(ctx, onlyActive)
}
