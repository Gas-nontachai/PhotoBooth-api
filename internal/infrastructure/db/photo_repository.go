package db

import (
	"context"

	"go-ddd-clean/internal/domain/media"

	"gorm.io/gorm"
)

type photoRepository struct {
	db *gorm.DB
}

func NewPhotoRepository(db *gorm.DB) media.PhotoRepository {
	return &photoRepository{db: db}
}

func (r *photoRepository) Create(ctx context.Context, p *media.Photo) error {
	model := PhotoModel{
		ID:          p.ID,
		SessionID:   p.SessionID,
		FrameID:     p.FrameID,
		FilterID:    p.FilterID,
		StorageURL:  p.StorageURL,
		Composition: toJSONMap(p.Composition),
		RenderedURL: p.RenderedURL,
	}
	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return err
	}
	p.CreatedAt = model.CreatedAt
	return nil
}

func (r *photoRepository) Update(ctx context.Context, p *media.Photo) error {
	return r.db.WithContext(ctx).
		Model(&PhotoModel{ID: p.ID}).
		Updates(map[string]any{
			"session_id":   p.SessionID,
			"frame_id":     p.FrameID,
			"filter_id":    p.FilterID,
			"storage_url":  p.StorageURL,
			"composition":  toJSONMap(p.Composition),
			"rendered_url": p.RenderedURL,
		}).Error
}

func (r *photoRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&PhotoModel{ID: id}).Error
}

func (r *photoRepository) GetByID(ctx context.Context, id string) (*media.Photo, error) {
	var model PhotoModel
	if err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return mapPhotoModelToDomain(&model), nil
}

func (r *photoRepository) ListBySession(ctx context.Context, sessionID string) ([]media.Photo, error) {
	var models []PhotoModel
	if err := r.db.WithContext(ctx).
		Where("session_id = ?", sessionID).
		Order("created_at asc").
		Find(&models).Error; err != nil {
		return nil, err
	}
	result := make([]media.Photo, 0, len(models))
	for _, m := range models {
		result = append(result, *mapPhotoModelToDomain(&m))
	}
	return result, nil
}

func mapPhotoModelToDomain(model *PhotoModel) *media.Photo {
	if model == nil {
		return nil
	}
	return &media.Photo{
		ID:          model.ID,
		SessionID:   model.SessionID,
		FrameID:     model.FrameID,
		FilterID:    model.FilterID,
		StorageURL:  model.StorageURL,
		Composition: fromJSONMap(model.Composition),
		RenderedURL: model.RenderedURL,
		CreatedAt:   model.CreatedAt,
	}
}
