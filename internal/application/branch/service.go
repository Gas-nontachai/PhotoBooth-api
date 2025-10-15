package branch

import (
	"context"

	"go-ddd-clean/internal/domain/branch"

	"github.com/google/uuid"
)

type Service struct {
	repo branch.Repository
}

func NewService(repo branch.Repository) *Service {
	return &Service{repo: repo}
}

type CreateBranchInput struct {
	Name     string
	Location *string
}

type UpdateBranchInput struct {
	ID       string
	Name     string
	Location *string
}

func (s *Service) Create(ctx context.Context, input CreateBranchInput) (*branch.Branch, error) {
	entity := &branch.Branch{
		ID:       uuid.NewString(),
		Name:     input.Name,
		Location: input.Location,
	}
	if err := s.repo.Create(ctx, entity); err != nil {
		return nil, err
	}
	return entity, nil
}

func (s *Service) Update(ctx context.Context, input UpdateBranchInput) error {
	entity, err := s.repo.GetByID(ctx, input.ID)
	if err != nil {
		return err
	}
	entity.Name = input.Name
	entity.Location = input.Location
	return s.repo.Update(ctx, entity)
}

func (s *Service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *Service) Get(ctx context.Context, id string) (*branch.Branch, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) List(ctx context.Context) ([]branch.Branch, error) {
	return s.repo.List(ctx)
}
