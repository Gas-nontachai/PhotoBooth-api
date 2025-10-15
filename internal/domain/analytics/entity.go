package analytics

import (
	"context"
	"time"
)

type Event struct {
	ID        string
	BoothID   string
	SessionID *string
	EventName string
	Payload   map[string]any
	CreatedAt time.Time
}

type Repository interface {
	Create(ctx context.Context, event *Event) error
	ListByBooth(ctx context.Context, boothID string, limit int) ([]Event, error)
}
