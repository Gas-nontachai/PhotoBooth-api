package media

import (
	"context"
	"time"
)

type Photo struct {
	ID          string
	SessionID   string
	FrameID     *string
	FilterID    *string
	StorageURL  string
	Composition map[string]any
	RenderedURL *string
	CreatedAt   time.Time
}

type Frame struct {
	ID        string
	Name      string
	Theme     *string
	FileURL   string
	Active    bool
	CreatedAt time.Time
}

type Filter struct {
	ID        string
	Name      string
	Effect    map[string]any
	Active    bool
	CreatedAt time.Time
}

type QRCode struct {
	ID        string
	PhotoID   string
	Hash      string
	ExpireAt  *time.Time
	CreatedAt time.Time
}

type PhotoRepository interface {
	Create(ctx context.Context, photo *Photo) error
	Update(ctx context.Context, photo *Photo) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*Photo, error)
	ListBySession(ctx context.Context, sessionID string) ([]Photo, error)
}

type FrameRepository interface {
	Create(ctx context.Context, frame *Frame) error
	Update(ctx context.Context, frame *Frame) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*Frame, error)
	List(ctx context.Context, onlyActive bool) ([]Frame, error)
}

type FilterRepository interface {
	Create(ctx context.Context, filter *Filter) error
	Update(ctx context.Context, filter *Filter) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*Filter, error)
	List(ctx context.Context, onlyActive bool) ([]Filter, error)
}

type QRCodeRepository interface {
	Create(ctx context.Context, code *QRCode) error
	GetByHash(ctx context.Context, hash string) (*QRCode, error)
	Delete(ctx context.Context, id string) error
}
