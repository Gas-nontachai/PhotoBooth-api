package db

import (
	"context"

	"go-ddd-clean/internal/domain/media"

	"gorm.io/gorm"
)

type filterRepository struct {
	db *gorm.DB
}

func NewFilterRepository(db *gorm.DB) media.FilterRepository {
	return &filterRepository{db: db}
}

func (r *filterRepository) Create(ctx context.Context, f *media.Filter) error {
	model := FilterModel{
		ID:     f.ID,
		Name:   f.Name,
		Effect: toJSONMap(f.Effect),
		Active: f.Active,
	}
	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return err
	}
	f.CreatedAt = model.CreatedAt
	return nil
}

func (r *filterRepository) Update(ctx context.Context, f *media.Filter) error {
	return r.db.WithContext(ctx).
		Model(&FilterModel{ID: f.ID}).
		Updates(map[string]any{
			"name":   f.Name,
			"effect": toJSONMap(f.Effect),
			"active": f.Active,
		}).Error
}

func (r *filterRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&FilterModel{ID: id}).Error
}

func (r *filterRepository) GetByID(ctx context.Context, id string) (*media.Filter, error) {
	var model FilterModel
	if err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &media.Filter{
		ID:        model.ID,
		Name:      model.Name,
		Effect:    fromJSONMap(model.Effect),
		Active:    model.Active,
		CreatedAt: model.CreatedAt,
	}, nil
}

func (r *filterRepository) List(ctx context.Context, onlyActive bool) ([]media.Filter, error) {
	query := r.db.WithContext(ctx)
	if onlyActive {
		query = query.Where("active = ?", true)
	}
	var models []FilterModel
	if err := query.Order("created_at desc").Find(&models).Error; err != nil {
		return nil, err
	}
	result := make([]media.Filter, 0, len(models))
	for _, m := range models {
		result = append(result, media.Filter{
			ID:        m.ID,
			Name:      m.Name,
			Effect:    fromJSONMap(m.Effect),
			Active:    m.Active,
			CreatedAt: m.CreatedAt,
		})
	}
	return result, nil
}
