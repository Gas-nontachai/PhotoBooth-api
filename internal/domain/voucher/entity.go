package voucher

import (
	"context"
	"time"
)

type Type string

const (
	TypeDiscount Type = "discount"
	TypeFree     Type = "free"
)

type Unit string

const (
	UnitPercent Unit = "percent"
	UnitBaht    Unit = "baht"
	UnitSession Unit = "session"
)

type Voucher struct {
	ID        string
	Code      string
	Type      Type
	Value     float64
	Unit      Unit
	MaxUsage  int
	UsedCount int
	ValidFrom *time.Time
	ValidTo   *time.Time
	Active    bool
	CreatedAt time.Time
}

type Redemption struct {
	ID        string
	VoucherID string
	SessionID string
	Tel       *string
	Discount  *float64
	CreatedAt time.Time
}

type Repository interface {
	Create(ctx context.Context, voucher *Voucher) error
	Update(ctx context.Context, voucher *Voucher) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*Voucher, error)
	GetByCode(ctx context.Context, code string) (*Voucher, error)
	List(ctx context.Context, activeOnly bool) ([]Voucher, error)
}

type RedemptionRepository interface {
	Create(ctx context.Context, redemption *Redemption) error
	ListByVoucher(ctx context.Context, voucherID string) ([]Redemption, error)
}
