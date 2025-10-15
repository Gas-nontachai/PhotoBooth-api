package payment

import (
	"context"
	"time"
)

type Method string

const (
	MethodCash   Method = "cash"
	MethodQR     Method = "qr"
	MethodStripe Method = "stripe"
	MethodPoints Method = "points"
)

type Status string

const (
	StatusPending Status = "pending"
	StatusSuccess Status = "success"
	StatusFailed  Status = "failed"
)

type Payment struct {
	ID             string
	SessionID      string
	Method         Method
	Amount         float64
	Currency       string
	Status         Status
	TransactionRef *string
	CreatedAt      time.Time
}

type Repository interface {
	Create(ctx context.Context, payment *Payment) error
	Update(ctx context.Context, payment *Payment) error
	GetByID(ctx context.Context, id string) (*Payment, error)
	GetBySessionID(ctx context.Context, sessionID string) (*Payment, error)
}
