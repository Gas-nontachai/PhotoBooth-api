package db

import (
	"context"

	"go-ddd-clean/internal/domain/media"

	"gorm.io/gorm"
)

type frameRepository struct {
	db *gorm.DB
}

func NewFrameRepository(db *gorm.DB) media.FrameRepository {
	return &frameRepository{db: db}
}

func (r *frameRepository) Create(ctx context.Context, f *media.Frame) error {
	model := FrameModel{
		ID:      f.ID,
		Name:    f.Name,
		Theme:   f.Theme,
		FileURL: f.FileURL,
		Active:  f.Active,
	}
	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return err
	}
	f.CreatedAt = model.CreatedAt
	return nil
}

func (r *frameRepository) Update(ctx context.Context, f *media.Frame) error {
	return r.db.WithContext(ctx).
		Model(&FrameModel{ID: f.ID}).
		Updates(map[string]any{
			"name":     f.Name,
			"theme":    f.Theme,
			"file_url": f.FileURL,
			"active":   f.Active,
		}).Error
}

func (r *frameRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&FrameModel{ID: id}).Error
}

func (r *frameRepository) GetByID(ctx context.Context, id string) (*media.Frame, error) {
	var model FrameModel
	if err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &media.Frame{
		ID:        model.ID,
		Name:      model.Name,
		Theme:     model.Theme,
		FileURL:   model.FileURL,
		Active:    model.Active,
		CreatedAt: model.CreatedAt,
	}, nil
}

func (r *frameRepository) List(ctx context.Context, onlyActive bool) ([]media.Frame, error) {
	query := r.db.WithContext(ctx)
	if onlyActive {
		query = query.Where("active = ?", true)
	}
	var models []FrameModel
	if err := query.Order("created_at desc").Find(&models).Error; err != nil {
		return nil, err
	}
	result := make([]media.Frame, 0, len(models))
	for _, m := range models {
		result = append(result, media.Frame{
			ID:        m.ID,
			Name:      m.Name,
			Theme:     m.Theme,
			FileURL:   m.FileURL,
			Active:    m.Active,
			CreatedAt: m.CreatedAt,
		})
	}
	return result, nil
}
