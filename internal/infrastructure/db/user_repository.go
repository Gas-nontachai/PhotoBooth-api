package db

import (
	"context"

	"go-ddd-clean/internal/domain/user"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) user.Repository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, u *user.User) error {
	model := UserModel{
		ID:       u.ID,
		Tel:      u.Tel,
		Email:    u.Email,
		Password: u.Password,
		Role:     string(u.Role),
		Points:   u.Points,
	}
	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return err
	}
	u.CreatedAt = model.CreatedAt
	u.UpdatedAt = model.UpdatedAt
	return nil
}

func (r *userRepository) Update(ctx context.Context, u *user.User) error {
	return r.db.WithContext(ctx).
		Model(&UserModel{ID: u.ID}).
		Updates(map[string]any{
			"tel":      u.Tel,
			"email":    u.Email,
			"password": u.Password,
			"role":     string(u.Role),
			"points":   u.Points,
		}).Error
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&UserModel{ID: id}).Error
}

func (r *userRepository) GetByID(ctx context.Context, id string) (*user.User, error) {
	var model UserModel
	if err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return mapUserModelToDomain(&model), nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	var model UserModel
	if err := r.db.WithContext(ctx).First(&model, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return mapUserModelToDomain(&model), nil
}

func (r *userRepository) List(ctx context.Context) ([]user.User, error) {
	var models []UserModel
	if err := r.db.WithContext(ctx).Find(&models).Error; err != nil {
		return nil, err
	}
	result := make([]user.User, 0, len(models))
	for _, m := range models {
		result = append(result, *mapUserModelToDomain(&m))
	}
	return result, nil
}

func mapUserModelToDomain(model *UserModel) *user.User {
	if model == nil {
		return nil
	}
	return &user.User{
		ID:        model.ID,
		Tel:       model.Tel,
		Email:     model.Email,
		Password:  model.Password,
		Role:      user.Role(model.Role),
		Points:    model.Points,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}
