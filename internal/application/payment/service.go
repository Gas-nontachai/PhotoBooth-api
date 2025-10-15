package payment

import (
	"context"

	domain "go-ddd-clean/internal/domain/payment"

	"github.com/google/uuid"
)

type Service struct {
	repo domain.Repository
}

func NewService(repo domain.Repository) *Service {
	return &Service{repo: repo}
}

type CreatePaymentInput struct {
	SessionID      string
	Method         domain.Method
	Amount         float64
	Currency       string
	Status         domain.Status
	TransactionRef *string
}

type UpdatePaymentStatusInput struct {
	ID             string
	Status         domain.Status
	TransactionRef *string
	Amount         *float64
	Currency       *string
	Method         *domain.Method
}

func (s *Service) Create(ctx context.Context, input CreatePaymentInput) (*domain.Payment, error) {
	status := input.Status
	if status == "" {
		status = domain.StatusPending
	}
	currency := input.Currency
	if currency == "" {
		currency = "THB"
	}
	entity := &domain.Payment{
		ID:             uuid.NewString(),
		SessionID:      input.SessionID,
		Method:         input.Method,
		Amount:         input.Amount,
		Currency:       currency,
		Status:         status,
		TransactionRef: input.TransactionRef,
	}
	if err := s.repo.Create(ctx, entity); err != nil {
		return nil, err
	}
	return entity, nil
}

func (s *Service) Update(ctx context.Context, input UpdatePaymentStatusInput) (*domain.Payment, error) {
	entity, err := s.repo.GetByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}
	if input.Status != "" {
		entity.Status = input.Status
	}
	if input.TransactionRef != nil {
		entity.TransactionRef = input.TransactionRef
	}
	if input.Amount != nil {
		entity.Amount = *input.Amount
	}
	if input.Currency != nil {
		entity.Currency = *input.Currency
	}
	if input.Method != nil {
		entity.Method = *input.Method
	}
	if err := s.repo.Update(ctx, entity); err != nil {
		return nil, err
	}
	return entity, nil
}

func (s *Service) Get(ctx context.Context, id string) (*domain.Payment, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) GetBySession(ctx context.Context, sessionID string) (*domain.Payment, error) {
	return s.repo.GetBySessionID(ctx, sessionID)
}
