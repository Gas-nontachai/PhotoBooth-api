package db

import (
	"context"

	"go-ddd-clean/internal/domain/logging"

	"gorm.io/gorm"
)

type logRepository struct {
	db *gorm.DB
}

func NewLogRepository(db *gorm.DB) logging.Repository {
	return &logRepository{db: db}
}

func (r *logRepository) Create(ctx context.Context, logEntry *logging.BoothLog) error {
	model := BoothLogModel{
		ID:        logEntry.ID,
		BoothID:   logEntry.BoothID,
		EventType: logEntry.EventType,
		Level:     string(logEntry.Level),
		Message:   logEntry.Message,
	}
	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return err
	}
	logEntry.CreatedAt = model.CreatedAt
	return nil
}

func (r *logRepository) ListByBooth(ctx context.Context, boothID string, limit int) ([]logging.BoothLog, error) {
	query := r.db.WithContext(ctx).Where("booth_id = ?", boothID).Order("created_at desc")
	if limit > 0 {
		query = query.Limit(limit)
	}
	var models []BoothLogModel
	if err := query.Find(&models).Error; err != nil {
		return nil, err
	}
	result := make([]logging.BoothLog, 0, len(models))
	for _, m := range models {
		result = append(result, logging.BoothLog{
			ID:        m.ID,
			BoothID:   m.BoothID,
			EventType: m.EventType,
			Level:     logging.Level(m.Level),
			Message:   m.Message,
			CreatedAt: m.CreatedAt,
		})
	}
	return result, nil
}
