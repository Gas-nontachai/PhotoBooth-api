package session

import (
	"context"
	"time"

	"go-ddd-clean/internal/domain/session"

	"github.com/google/uuid"
)

type Service struct {
	repo session.Repository
}

func NewService(repo session.Repository) *Service {
	return &Service{repo: repo}
}

type CreateSessionInput struct {
	BoothID       string
	UserID        *string
	VoucherID     *string
	PaymentID     *string
	Status        session.Status
	TotalPrice    *float64
	BoothSnapshot map[string]any
	PhoneTemp     *string
}

type UpdateSessionInput struct {
	ID            string
	UserID        *string
	VoucherID     *string
	PaymentID     *string
	Status        *session.Status
	TotalPrice    *float64
	FinishedAt    *time.Time
	BoothSnapshot map[string]any
	PhoneTemp     *string
}

func (s *Service) Create(ctx context.Context, input CreateSessionInput) (*session.Session, error) {
	now := time.Now()
	status := input.Status
	if status == "" {
		status = session.StatusStarted
	}
	entity := &session.Session{
		ID:            uuid.NewString(),
		BoothID:       input.BoothID,
		UserID:        input.UserID,
		VoucherID:     input.VoucherID,
		PaymentID:     input.PaymentID,
		StartedAt:     &now,
		Status:        status,
		TotalPrice:    input.TotalPrice,
		BoothSnapshot: input.BoothSnapshot,
		PhoneTemp:     input.PhoneTemp,
	}
	if err := s.repo.Create(ctx, entity); err != nil {
		return nil, err
	}
	return entity, nil
}

func (s *Service) Update(ctx context.Context, input UpdateSessionInput) (*session.Session, error) {
	entity, err := s.repo.GetByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}
	if input.UserID != nil {
		entity.UserID = input.UserID
	}
	if input.VoucherID != nil {
		entity.VoucherID = input.VoucherID
	}
	if input.PaymentID != nil {
		entity.PaymentID = input.PaymentID
	}
	if input.Status != nil {
		entity.Status = *input.Status
	}
	if input.TotalPrice != nil {
		entity.TotalPrice = input.TotalPrice
	}
	if input.FinishedAt != nil {
		entity.FinishedAt = input.FinishedAt
	}
	if input.BoothSnapshot != nil {
		entity.BoothSnapshot = input.BoothSnapshot
	}
	if input.PhoneTemp != nil {
		entity.PhoneTemp = input.PhoneTemp
	}
	if err := s.repo.Update(ctx, entity); err != nil {
		return nil, err
	}
	return entity, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *Service) Get(ctx context.Context, id string) (*session.Session, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) List(ctx context.Context, boothID *string, status *session.Status) ([]session.Session, error) {
	return s.repo.List(ctx, boothID, status)
}
