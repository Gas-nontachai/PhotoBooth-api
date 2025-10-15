package db

import (
	"context"

	"go-ddd-clean/internal/domain/media"

	"gorm.io/gorm"
)

type qrCodeRepository struct {
	db *gorm.DB
}

func NewQRCodeRepository(db *gorm.DB) media.QRCodeRepository {
	return &qrCodeRepository{db: db}
}

func (r *qrCodeRepository) Create(ctx context.Context, code *media.QRCode) error {
	model := QRCodeModel{
		ID:       code.ID,
		PhotoID:  code.PhotoID,
		Hash:     code.Hash,
		ExpireAt: code.ExpireAt,
	}
	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return err
	}
	code.CreatedAt = model.CreatedAt
	return nil
}

func (r *qrCodeRepository) GetByHash(ctx context.Context, hash string) (*media.QRCode, error) {
	var model QRCodeModel
	if err := r.db.WithContext(ctx).First(&model, "hash = ?", hash).Error; err != nil {
		return nil, err
	}
	return &media.QRCode{
		ID:        model.ID,
		PhotoID:   model.PhotoID,
		Hash:      model.Hash,
		ExpireAt:  model.ExpireAt,
		CreatedAt: model.CreatedAt,
	}, nil
}

func (r *qrCodeRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&QRCodeModel{ID: id}).Error
}
