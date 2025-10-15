package logging

import (
	"context"
	"time"
)

type Level string

const (
	LevelInfo  Level = "info"
	LevelWarn  Level = "warn"
	LevelError Level = "error"
)

type BoothLog struct {
	ID        string
	BoothID   string
	EventType string
	Level     Level
	Message   *string
	CreatedAt time.Time
}

type Repository interface {
	Create(ctx context.Context, log *BoothLog) error
	ListByBooth(ctx context.Context, boothID string, limit int) ([]BoothLog, error)
}
