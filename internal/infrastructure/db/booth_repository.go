package db

import (
	"context"

	"go-ddd-clean/internal/domain/booth"

	"gorm.io/gorm"
)

type boothRepository struct {
	db *gorm.DB
}

func NewBoothRepository(db *gorm.DB) booth.Repository {
	return &boothRepository{db: db}
}

func (r *boothRepository) Create(ctx context.Context, b *booth.Booth) error {
	model := BoothModel{
		ID:           b.ID,
		BranchID:     b.BranchID,
		Name:         b.Name,
		Type:         string(b.Type),
		Status:       string(b.Status),
		Config:       toJSONMap(b.Config),
		TokenVersion: b.TokenVersion,
	}
	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return err
	}
	b.CreatedAt = model.CreatedAt
	return nil
}

func (r *boothRepository) Update(ctx context.Context, b *booth.Booth) error {
	return r.db.WithContext(ctx).
		Model(&BoothModel{ID: b.ID}).
		Updates(map[string]any{
			"branch_id":     b.BranchID,
			"name":          b.Name,
			"type":          string(b.Type),
			"status":        string(b.Status),
			"config":        toJSONMap(b.Config),
			"token_version": b.TokenVersion,
		}).Error
}

func (r *boothRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&BoothModel{ID: id}).Error
}

func (r *boothRepository) GetByID(ctx context.Context, id string) (*booth.Booth, error) {
	var model BoothModel
	if err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return mapBoothModelToDomain(&model), nil
}

func (r *boothRepository) List(ctx context.Context, branchID *string) ([]booth.Booth, error) {
	query := r.db.WithContext(ctx).Model(&BoothModel{})
	if branchID != nil {
		query = query.Where("branch_id = ?", *branchID)
	}
	var models []BoothModel
	if err := query.Find(&models).Error; err != nil {
		return nil, err
	}
	result := make([]booth.Booth, 0, len(models))
	for _, m := range models {
		result = append(result, *mapBoothModelToDomain(&m))
	}
	return result, nil
}

func mapBoothModelToDomain(model *BoothModel) *booth.Booth {
	if model == nil {
		return nil
	}
	return &booth.Booth{
		ID:           model.ID,
		BranchID:     model.BranchID,
		Name:         model.Name,
		Type:         booth.BoothType(model.Type),
		Status:       booth.BoothStatus(model.Status),
		Config:       fromJSONMap(model.Config),
		TokenVersion: model.TokenVersion,
		CreatedAt:    model.CreatedAt,
	}
}

func (r *boothRepository) UpdateTokenVersion(ctx context.Context, id string, version int) error {
	return r.db.WithContext(ctx).
		Model(&BoothModel{ID: id}).
		Update("token_version", version).Error
}
