package branch

import (
	"context"
	"time"
)

type Branch struct {
	ID        string
	Name      string
	Location  *string
	CreatedAt time.Time
}

type Repository interface {
	Create(ctx context.Context, branch *Branch) error
	Update(ctx context.Context, branch *Branch) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*Branch, error)
	List(ctx context.Context) ([]Branch, error)
}
