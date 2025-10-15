package booth

import (
	"context"

	"go-ddd-clean/internal/domain/booth"

	"github.com/google/uuid"
)

type Service struct {
	repo booth.Repository
}

func NewService(repo booth.Repository) *Service {
	return &Service{repo: repo}
}

type CreateBoothInput struct {
	BranchID string
	Name     string
	Type     booth.BoothType
	Status   booth.BoothStatus
	Config   map[string]any
}

type UpdateBoothInput struct {
	ID       string
	BranchID string
	Name     string
	Type     booth.BoothType
	Status   booth.BoothStatus
	Config   map[string]any
}

func (s *Service) Create(ctx context.Context, input CreateBoothInput) (*booth.Booth, error) {
	status := input.Status
	if status == "" {
		status = booth.BoothStatusActive
	}
	entity := &booth.Booth{
		ID:       uuid.NewString(),
		BranchID: input.BranchID,
		Name:     input.Name,
		Type:     input.Type,
		Status:   status,
		Config:   input.Config,
	}
	if err := s.repo.Create(ctx, entity); err != nil {
		return nil, err
	}
	return entity, nil
}

func (s *Service) Update(ctx context.Context, input UpdateBoothInput) error {
	entity, err := s.repo.GetByID(ctx, input.ID)
	if err != nil {
		return err
	}
	entity.BranchID = input.BranchID
	entity.Name = input.Name
	if input.Type != "" {
		entity.Type = input.Type
	}
	if input.Status != "" {
		entity.Status = input.Status
	}
	entity.Config = input.Config
	return s.repo.Update(ctx, entity)
}

func (s *Service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *Service) Get(ctx context.Context, id string) (*booth.Booth, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) List(ctx context.Context, branchID *string) ([]booth.Booth, error) {
	return s.repo.List(ctx, branchID)
}
