package media

import (
	"context"

	domain "go-ddd-clean/internal/domain/media"

	"github.com/google/uuid"
)

type PhotoService struct {
	repo domain.PhotoRepository
}

func NewPhotoService(repo domain.PhotoRepository) *PhotoService {
	return &PhotoService{repo: repo}
}

type CreatePhotoInput struct {
	SessionID   string
	FrameID     *string
	FilterID    *string
	StorageURL  string
	Composition map[string]any
	RenderedURL *string
}

type UpdatePhotoInput struct {
	ID          string
	FrameID     *string
	FilterID    *string
	StorageURL  *string
	Composition map[string]any
	RenderedURL *string
}

func (s *PhotoService) Create(ctx context.Context, input CreatePhotoInput) (*domain.Photo, error) {
	entity := &domain.Photo{
		ID:          uuid.NewString(),
		SessionID:   input.SessionID,
		FrameID:     input.FrameID,
		FilterID:    input.FilterID,
		StorageURL:  input.StorageURL,
		Composition: input.Composition,
		RenderedURL: input.RenderedURL,
	}
	if err := s.repo.Create(ctx, entity); err != nil {
		return nil, err
	}
	return entity, nil
}

func (s *PhotoService) Update(ctx context.Context, input UpdatePhotoInput) (*domain.Photo, error) {
	entity, err := s.repo.GetByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}
	if input.FrameID != nil {
		entity.FrameID = input.FrameID
	}
	if input.FilterID != nil {
		entity.FilterID = input.FilterID
	}
	if input.StorageURL != nil {
		entity.StorageURL = *input.StorageURL
	}
	if input.Composition != nil {
		entity.Composition = input.Composition
	}
	if input.RenderedURL != nil {
		entity.RenderedURL = input.RenderedURL
	}
	if err := s.repo.Update(ctx, entity); err != nil {
		return nil, err
	}
	return entity, nil
}

func (s *PhotoService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *PhotoService) Get(ctx context.Context, id string) (*domain.Photo, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *PhotoService) ListBySession(ctx context.Context, sessionID string) ([]domain.Photo, error) {
	return s.repo.ListBySession(ctx, sessionID)
}
