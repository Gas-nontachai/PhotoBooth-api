package session

import (
	"context"
	"time"
)

type Status string

const (
	StatusStarted   Status = "started"
	StatusSuccess   Status = "success"
	StatusFailed    Status = "failed"
	StatusCancelled Status = "cancelled"
)

type Session struct {
	ID            string
	BoothID       string
	UserID        *string
	VoucherID     *string
	PaymentID     *string
	StartedAt     *time.Time
	FinishedAt    *time.Time
	Status        Status
	TotalPrice    *float64
	BoothSnapshot map[string]any
	PhoneTemp     *string
}

type Repository interface {
	Create(ctx context.Context, session *Session) error
	Update(ctx context.Context, session *Session) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*Session, error)
	List(ctx context.Context, boothID *string, status *Status) ([]Session, error)
}
