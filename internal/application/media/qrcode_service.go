package media

import (
	"context"
	"time"

	domain "go-ddd-clean/internal/domain/media"

	"github.com/google/uuid"
)

type QRCodeService struct {
	repo domain.QRCodeRepository
}

func NewQRCodeService(repo domain.QRCodeRepository) *QRCodeService {
	return &QRCodeService{repo: repo}
}

type CreateQRCodeInput struct {
	PhotoID  string
	Hash     string
	ExpireAt *time.Time
}

func (s *QRCodeService) Create(ctx context.Context, input CreateQRCodeInput) (*domain.QRCode, error) {
	entity := &domain.QRCode{
		ID:       uuid.NewString(),
		PhotoID:  input.PhotoID,
		Hash:     input.Hash,
		ExpireAt: input.ExpireAt,
	}
	if err := s.repo.Create(ctx, entity); err != nil {
		return nil, err
	}
	return entity, nil
}

func (s *QRCodeService) GetByHash(ctx context.Context, hash string) (*domain.QRCode, error) {
	return s.repo.GetByHash(ctx, hash)
}

func (s *QRCodeService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
