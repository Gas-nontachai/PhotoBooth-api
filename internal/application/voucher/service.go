package voucher

import (
	"context"
	"errors"
	"time"

	domain "go-ddd-clean/internal/domain/voucher"

	"github.com/google/uuid"
)

type Service struct {
	voucherRepo    domain.Repository
	redemptionRepo domain.RedemptionRepository
}

func NewService(voucherRepo domain.Repository, redemptionRepo domain.RedemptionRepository) *Service {
	return &Service{
		voucherRepo:    voucherRepo,
		redemptionRepo: redemptionRepo,
	}
}

type CreateVoucherInput struct {
	Code      string
	Type      domain.Type
	Value     float64
	Unit      domain.Unit
	MaxUsage  int
	ValidFrom *time.Time
	ValidTo   *time.Time
	Active    bool
}

type UpdateVoucherInput struct {
	ID        string
	Value     *float64
	Unit      *domain.Unit
	MaxUsage  *int
	ValidFrom *time.Time
	ValidTo   *time.Time
	Active    *bool
}

type RedeemVoucherInput struct {
	Code      string
	SessionID string
	Tel       *string
	Discount  *float64
}

func (s *Service) Create(ctx context.Context, input CreateVoucherInput) (*domain.Voucher, error) {
	if input.MaxUsage == 0 {
		input.MaxUsage = 1
	}
	entity := &domain.Voucher{
		ID:        uuid.NewString(),
		Code:      input.Code,
		Type:      input.Type,
		Value:     input.Value,
		Unit:      input.Unit,
		MaxUsage:  input.MaxUsage,
		ValidFrom: input.ValidFrom,
		ValidTo:   input.ValidTo,
		Active:    input.Active,
	}
	if err := s.voucherRepo.Create(ctx, entity); err != nil {
		return nil, err
	}
	return entity, nil
}

func (s *Service) Update(ctx context.Context, input UpdateVoucherInput) (*domain.Voucher, error) {
	entity, err := s.voucherRepo.GetByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}
	if input.Value != nil {
		entity.Value = *input.Value
	}
	if input.Unit != nil {
		entity.Unit = *input.Unit
	}
	if input.MaxUsage != nil {
		entity.MaxUsage = *input.MaxUsage
	}
	if input.ValidFrom != nil {
		entity.ValidFrom = input.ValidFrom
	}
	if input.ValidTo != nil {
		entity.ValidTo = input.ValidTo
	}
	if input.Active != nil {
		entity.Active = *input.Active
	}
	if err := s.voucherRepo.Update(ctx, entity); err != nil {
		return nil, err
	}
	return entity, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	return s.voucherRepo.Delete(ctx, id)
}

func (s *Service) Get(ctx context.Context, id string) (*domain.Voucher, error) {
	return s.voucherRepo.GetByID(ctx, id)
}

func (s *Service) GetByCode(ctx context.Context, code string) (*domain.Voucher, error) {
	return s.voucherRepo.GetByCode(ctx, code)
}

func (s *Service) List(ctx context.Context, activeOnly bool) ([]domain.Voucher, error) {
	return s.voucherRepo.List(ctx, activeOnly)
}

func (s *Service) Redeem(ctx context.Context, input RedeemVoucherInput) (*domain.Redemption, *domain.Voucher, error) {
	v, err := s.voucherRepo.GetByCode(ctx, input.Code)
	if err != nil {
		return nil, nil, err
	}
	if !v.Active {
		return nil, nil, errors.New("voucher inactive")
	}
	now := time.Now()
	if v.ValidFrom != nil && now.Before(*v.ValidFrom) {
		return nil, nil, errors.New("voucher not yet valid")
	}
	if v.ValidTo != nil && now.After(*v.ValidTo) {
		return nil, nil, errors.New("voucher expired")
	}
	if v.MaxUsage > 0 && v.UsedCount >= v.MaxUsage {
		return nil, nil, errors.New("voucher usage limit reached")
	}

	redemption := &domain.Redemption{
		ID:        uuid.NewString(),
		VoucherID: v.ID,
		SessionID: input.SessionID,
		Tel:       input.Tel,
		Discount:  input.Discount,
	}
	if err := s.redemptionRepo.Create(ctx, redemption); err != nil {
		return nil, nil, err
	}
	v.UsedCount++
	if err := s.voucherRepo.Update(ctx, v); err != nil {
		return nil, nil, err
	}
	return redemption, v, nil
}
