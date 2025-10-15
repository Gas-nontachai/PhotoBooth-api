package booth

import (
	"context"
	"time"
)

type BoothType string

const (
	BoothTypePhysical BoothType = "physical"
	BoothTypeVirtual  BoothType = "virtual"
)

type BoothStatus string

const (
	BoothStatusActive   BoothStatus = "active"
	BoothStatusInactive BoothStatus = "inactive"
)

type Booth struct {
	ID           string
	BranchID     string
	Name         string
	Type         BoothType
	Status       BoothStatus
	Config       map[string]any
	TokenVersion int
	CreatedAt    time.Time
}

type Repository interface {
	Create(ctx context.Context, booth *Booth) error
	Update(ctx context.Context, booth *Booth) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*Booth, error)
	List(ctx context.Context, branchID *string) ([]Booth, error)
	UpdateTokenVersion(ctx context.Context, id string, version int) error
}
