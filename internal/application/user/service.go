package user

import (
	"context"
	"errors"

	domain "go-ddd-clean/internal/domain/user"

	"github.com/google/uuid"
)

type Service struct {
	repo domain.Repository
}

func NewService(repo domain.Repository) *Service {
	return &Service{repo: repo}
}

type CreateUserInput struct {
	Tel      *string
	Email    *string
	Password *string
	Role     domain.Role
}

type UpdateUserInput struct {
	ID       string
	Tel      *string
	Email    *string
	Password *string
	Role     *domain.Role
	Points   *int
}

func (s *Service) Create(ctx context.Context, input CreateUserInput) (*domain.User, error) {
	role := input.Role
	if role == "" {
		role = domain.RoleCustomer
	}
	entity := &domain.User{
		ID:       uuid.NewString(),
		Tel:      input.Tel,
		Email:    input.Email,
		Password: input.Password,
		Role:     role,
		Points:   0,
	}
	if err := s.repo.Create(ctx, entity); err != nil {
		return nil, err
	}
	return entity, nil
}

func (s *Service) Update(ctx context.Context, input UpdateUserInput) (*domain.User, error) {
	entity, err := s.repo.GetByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}
	if input.Tel != nil {
		entity.Tel = input.Tel
	}
	if input.Email != nil {
		entity.Email = input.Email
	}
	if input.Password != nil {
		entity.Password = input.Password
	}
	if input.Role != nil {
		entity.Role = *input.Role
	}
	if input.Points != nil {
		entity.Points = *input.Points
	}
	if err := s.repo.Update(ctx, entity); err != nil {
		return nil, err
	}
	return entity, nil
}

func (s *Service) AdjustPoints(ctx context.Context, userID string, delta int) (*domain.User, error) {
	entity, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	newPoints := entity.Points + delta
	if newPoints < 0 {
		return nil, errors.New("points cannot be negative")
	}
	entity.Points = newPoints
	if err := s.repo.Update(ctx, entity); err != nil {
		return nil, err
	}
	return entity, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *Service) Get(ctx context.Context, id string) (*domain.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	return s.repo.GetByEmail(ctx, email)
}

func (s *Service) List(ctx context.Context) ([]domain.User, error) {
	return s.repo.List(ctx)
}
