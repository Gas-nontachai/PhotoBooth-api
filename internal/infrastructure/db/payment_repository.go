package db

import (
	"context"

	"go-ddd-clean/internal/domain/payment"

	"gorm.io/gorm"
)

type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) payment.Repository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) Create(ctx context.Context, p *payment.Payment) error {
	model := PaymentModel{
		ID:             p.ID,
		SessionID:      p.SessionID,
		Method:         string(p.Method),
		Amount:         p.Amount,
		Currency:       p.Currency,
		Status:         string(p.Status),
		TransactionRef: p.TransactionRef,
	}
	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return err
	}
	p.CreatedAt = model.CreatedAt
	return nil
}

func (r *paymentRepository) Update(ctx context.Context, p *payment.Payment) error {
	return r.db.WithContext(ctx).
		Model(&PaymentModel{ID: p.ID}).
		Updates(map[string]any{
			"status":          string(p.Status),
			"transaction_ref": p.TransactionRef,
			"amount":          p.Amount,
			"currency":        p.Currency,
			"method":          string(p.Method),
		}).Error
}

func (r *paymentRepository) GetByID(ctx context.Context, id string) (*payment.Payment, error) {
	var model PaymentModel
	if err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return mapPaymentModelToDomain(&model), nil
}

func (r *paymentRepository) GetBySessionID(ctx context.Context, sessionID string) (*payment.Payment, error) {
	var model PaymentModel
	if err := r.db.WithContext(ctx).First(&model, "session_id = ?", sessionID).Error; err != nil {
		return nil, err
	}
	return mapPaymentModelToDomain(&model), nil
}

func mapPaymentModelToDomain(model *PaymentModel) *payment.Payment {
	if model == nil {
		return nil
	}
	return &payment.Payment{
		ID:             model.ID,
		SessionID:      model.SessionID,
		Method:         payment.Method(model.Method),
		Amount:         model.Amount,
		Currency:       model.Currency,
		Status:         payment.Status(model.Status),
		TransactionRef: model.TransactionRef,
		CreatedAt:      model.CreatedAt,
	}
}
