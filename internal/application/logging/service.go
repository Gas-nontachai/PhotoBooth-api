package logging

import (
	"context"

	domain "go-ddd-clean/internal/domain/logging"

	"github.com/google/uuid"
)

type Service struct {
	repo domain.Repository
}

func NewService(repo domain.Repository) *Service {
	return &Service{repo: repo}
}

type CreateLogInput struct {
	BoothID   string
	EventType string
	Level     domain.Level
	Message   *string
}

func (s *Service) Write(ctx context.Context, input CreateLogInput) (*domain.BoothLog, error) {
	level := input.Level
	if level == "" {
		level = domain.LevelInfo
	}
	entry := &domain.BoothLog{
		ID:        uuid.NewString(),
		BoothID:   input.BoothID,
		EventType: input.EventType,
		Level:     level,
		Message:   input.Message,
	}
	if err := s.repo.Create(ctx, entry); err != nil {
		return nil, err
	}
	return entry, nil
}

func (s *Service) List(ctx context.Context, boothID string, limit int) ([]domain.BoothLog, error) {
	return s.repo.ListByBooth(ctx, boothID, limit)
}
