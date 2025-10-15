package db

import (
	"context"

	"go-ddd-clean/internal/domain/session"

	"gorm.io/gorm"
)

type sessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) session.Repository {
	return &sessionRepository{db: db}
}

func (r *sessionRepository) Create(ctx context.Context, s *session.Session) error {
	model := SessionModel{
		ID:            s.ID,
		BoothID:       s.BoothID,
		UserID:        s.UserID,
		VoucherID:     s.VoucherID,
		PaymentID:     s.PaymentID,
		StartedAt:     s.StartedAt,
		FinishedAt:    s.FinishedAt,
		Status:        string(s.Status),
		TotalPrice:    s.TotalPrice,
		BoothSnapshot: toJSONMap(s.BoothSnapshot),
		PhoneTemp:     s.PhoneTemp,
	}
	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return err
	}
	return nil
}

func (r *sessionRepository) Update(ctx context.Context, s *session.Session) error {
	return r.db.WithContext(ctx).
		Model(&SessionModel{ID: s.ID}).
		Updates(map[string]any{
			"booth_id":       s.BoothID,
			"user_id":        s.UserID,
			"voucher_id":     s.VoucherID,
			"payment_id":     s.PaymentID,
			"started_at":     s.StartedAt,
			"finished_at":    s.FinishedAt,
			"status":         string(s.Status),
			"total_price":    s.TotalPrice,
			"booth_snapshot": toJSONMap(s.BoothSnapshot),
			"phone_temp":     s.PhoneTemp,
		}).Error
}

func (r *sessionRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&SessionModel{ID: id}).Error
}

func (r *sessionRepository) GetByID(ctx context.Context, id string) (*session.Session, error) {
	var model SessionModel
	if err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return mapSessionModelToDomain(&model), nil
}

func (r *sessionRepository) List(ctx context.Context, boothID *string, status *session.Status) ([]session.Session, error) {
	query := r.db.WithContext(ctx).Model(&SessionModel{})
	if boothID != nil {
		query = query.Where("booth_id = ?", *boothID)
	}
	if status != nil {
		query = query.Where("status = ?", string(*status))
	}
	var models []SessionModel
	if err := query.Order("started_at desc").Find(&models).Error; err != nil {
		return nil, err
	}
	result := make([]session.Session, 0, len(models))
	for _, m := range models {
		result = append(result, *mapSessionModelToDomain(&m))
	}
	return result, nil
}

func mapSessionModelToDomain(model *SessionModel) *session.Session {
	if model == nil {
		return nil
	}
	return &session.Session{
		ID:            model.ID,
		BoothID:       model.BoothID,
		UserID:        model.UserID,
		VoucherID:     model.VoucherID,
		PaymentID:     model.PaymentID,
		StartedAt:     model.StartedAt,
		FinishedAt:    model.FinishedAt,
		Status:        session.Status(model.Status),
		TotalPrice:    model.TotalPrice,
		BoothSnapshot: fromJSONMap(model.BoothSnapshot),
		PhoneTemp:     model.PhoneTemp,
	}
}
