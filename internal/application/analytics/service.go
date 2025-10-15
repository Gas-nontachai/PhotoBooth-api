package analytics

import (
	"context"

	domain "go-ddd-clean/internal/domain/analytics"

	"github.com/google/uuid"
)

type Service struct {
	repo domain.Repository
}

func NewService(repo domain.Repository) *Service {
	return &Service{repo: repo}
}

type CreateEventInput struct {
	BoothID   string
	SessionID *string
	EventName string
	Payload   map[string]any
}

func (s *Service) Create(ctx context.Context, input CreateEventInput) (*domain.Event, error) {
	event := &domain.Event{
		ID:        uuid.NewString(),
		BoothID:   input.BoothID,
		SessionID: input.SessionID,
		EventName: input.EventName,
		Payload:   input.Payload,
	}
	if err := s.repo.Create(ctx, event); err != nil {
		return nil, err
	}
	return event, nil
}

func (s *Service) ListByBooth(ctx context.Context, boothID string, limit int) ([]domain.Event, error) {
	return s.repo.ListByBooth(ctx, boothID, limit)
}
