package db

import (
	"context"

	"go-ddd-clean/internal/domain/branch"

	"gorm.io/gorm"
)

type branchRepository struct {
	db *gorm.DB
}

func NewBranchRepository(db *gorm.DB) branch.Repository {
	return &branchRepository{db: db}
}

func (r *branchRepository) Create(ctx context.Context, b *branch.Branch) error {
	model := BranchModel{
		ID:       b.ID,
		Name:     b.Name,
		Location: b.Location,
	}
	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return err
	}
	b.CreatedAt = model.CreatedAt
	return nil
}

func (r *branchRepository) Update(ctx context.Context, b *branch.Branch) error {
	return r.db.WithContext(ctx).
		Model(&BranchModel{ID: b.ID}).
		Updates(map[string]any{
			"name":     b.Name,
			"location": b.Location,
		}).Error
}

func (r *branchRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&BranchModel{ID: id}).Error
}

func (r *branchRepository) GetByID(ctx context.Context, id string) (*branch.Branch, error) {
	var model BranchModel
	if err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &branch.Branch{
		ID:        model.ID,
		Name:      model.Name,
		Location:  model.Location,
		CreatedAt: model.CreatedAt,
	}, nil
}

func (r *branchRepository) List(ctx context.Context) ([]branch.Branch, error) {
	var models []BranchModel
	if err := r.db.WithContext(ctx).Find(&models).Error; err != nil {
		return nil, err
	}
	result := make([]branch.Branch, 0, len(models))
	for _, m := range models {
		result = append(result, branch.Branch{
			ID:        m.ID,
			Name:      m.Name,
			Location:  m.Location,
			CreatedAt: m.CreatedAt,
		})
	}
	return result, nil
}
