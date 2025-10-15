package db

import (
	"context"

	"go-ddd-clean/internal/domain/analytics"

	"gorm.io/gorm"
)

type analyticsRepository struct {
	db *gorm.DB
}

func NewAnalyticsRepository(db *gorm.DB) analytics.Repository {
	return &analyticsRepository{db: db}
}

func (r *analyticsRepository) Create(ctx context.Context, event *analytics.Event) error {
	model := AnalyticsEventModel{
		ID:        event.ID,
		BoothID:   event.BoothID,
		SessionID: event.SessionID,
		EventName: event.EventName,
		Payload:   toJSONMap(event.Payload),
	}
	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return err
	}
	event.CreatedAt = model.CreatedAt
	return nil
}

func (r *analyticsRepository) ListByBooth(ctx context.Context, boothID string, limit int) ([]analytics.Event, error) {
	query := r.db.WithContext(ctx).Where("booth_id = ?", boothID).Order("created_at desc")
	if limit > 0 {
		query = query.Limit(limit)
	}
	var models []AnalyticsEventModel
	if err := query.Find(&models).Error; err != nil {
		return nil, err
	}
	result := make([]analytics.Event, 0, len(models))
	for _, m := range models {
		result = append(result, analytics.Event{
			ID:        m.ID,
			BoothID:   m.BoothID,
			SessionID: m.SessionID,
			EventName: m.EventName,
			Payload:   fromJSONMap(m.Payload),
			CreatedAt: m.CreatedAt,
		})
	}
	return result, nil
}
