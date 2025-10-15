package db

import (
	"time"

	"gorm.io/datatypes"
)

type BranchModel struct {
	ID        string `gorm:"type:uuid;primaryKey"`
	Name      string
	Location  *string
	CreatedAt time.Time `gorm:"autoCreateTime"`

	Booths []BoothModel `gorm:"foreignKey:BranchID"`
}

type BoothModel struct {
	ID           string `gorm:"type:uuid;primaryKey"`
	BranchID     string `gorm:"type:uuid;index"`
	Name         string
	Type         string
	Status       string            `gorm:"default:active"`
	Config       datatypes.JSONMap `gorm:"type:jsonb"`
	TokenVersion int               `gorm:"default:0"`
	CreatedAt    time.Time         `gorm:"autoCreateTime"`

	Branch    BranchModel
	Sessions  []SessionModel        `gorm:"foreignKey:BoothID"`
	Logs      []BoothLogModel       `gorm:"foreignKey:BoothID"`
	Analytics []AnalyticsEventModel `gorm:"foreignKey:BoothID"`
}

type SessionModel struct {
	ID            string     `gorm:"type:uuid;primaryKey"`
	BoothID       string     `gorm:"type:uuid;index"`
	UserID        *string    `gorm:"type:uuid"`
	VoucherID     *string    `gorm:"type:uuid"`
	PaymentID     *string    `gorm:"type:uuid"`
	StartedAt     *time.Time `gorm:"autoCreateTime"`
	FinishedAt    *time.Time
	Status        string `gorm:"default:started"`
	TotalPrice    *float64
	BoothSnapshot datatypes.JSONMap `gorm:"type:jsonb"`
	PhoneTemp     *string

	Photos      []PhotoModel             `gorm:"foreignKey:SessionID"`
	Analytics   []AnalyticsEventModel    `gorm:"foreignKey:SessionID"`
	Redemptions []VoucherRedemptionModel `gorm:"foreignKey:SessionID"`
	Payment     PaymentModel
}

type PhotoModel struct {
	ID          string  `gorm:"type:uuid;primaryKey"`
	SessionID   string  `gorm:"type:uuid;index"`
	FrameID     *string `gorm:"type:uuid"`
	FilterID    *string `gorm:"type:uuid"`
	StorageURL  string
	Composition datatypes.JSONMap `gorm:"type:jsonb"`
	RenderedURL *string
	CreatedAt   time.Time `gorm:"autoCreateTime"`

	QRCodes []QRCodeModel `gorm:"foreignKey:PhotoID"`
}

type FrameModel struct {
	ID        string `gorm:"type:uuid;primaryKey"`
	Name      string
	Theme     *string
	FileURL   string
	Active    bool      `gorm:"default:true"`
	CreatedAt time.Time `gorm:"autoCreateTime"`

	Photos []PhotoModel `gorm:"foreignKey:FrameID"`
}

type FilterModel struct {
	ID        string `gorm:"type:uuid;primaryKey"`
	Name      string
	Effect    datatypes.JSONMap `gorm:"type:jsonb"`
	Active    bool              `gorm:"default:true"`
	CreatedAt time.Time         `gorm:"autoCreateTime"`

	Photos []PhotoModel `gorm:"foreignKey:FilterID"`
}

type QRCodeModel struct {
	ID        string `gorm:"type:uuid;primaryKey"`
	PhotoID   string `gorm:"type:uuid;index"`
	Hash      string `gorm:"uniqueIndex"`
	ExpireAt  *time.Time
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

type UserModel struct {
	ID        string  `gorm:"type:uuid;primaryKey"`
	Tel       *string `gorm:"unique"`
	Email     *string `gorm:"unique"`
	Password  *string
	Role      string    `gorm:"default:customer"`
	Points    int       `gorm:"default:0"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	Sessions []SessionModel `gorm:"foreignKey:UserID"`
}

type PaymentModel struct {
	ID             string `gorm:"type:uuid;primaryKey"`
	SessionID      string `gorm:"type:uuid;uniqueIndex"`
	Method         string
	Amount         float64
	Currency       string `gorm:"default:THB"`
	Status         string `gorm:"default:pending"`
	TransactionRef *string
	CreatedAt      time.Time `gorm:"autoCreateTime"`
}

type VoucherModel struct {
	ID        string `gorm:"type:uuid;primaryKey"`
	Code      string `gorm:"uniqueIndex"`
	Type      string
	Value     float64
	Unit      string
	MaxUsage  int `gorm:"default:1"`
	UsedCount int `gorm:"default:0"`
	ValidFrom *time.Time
	ValidTo   *time.Time
	Active    bool      `gorm:"default:true"`
	CreatedAt time.Time `gorm:"autoCreateTime"`

	Redemptions []VoucherRedemptionModel `gorm:"foreignKey:VoucherID"`
}

type VoucherRedemptionModel struct {
	ID        string `gorm:"type:uuid;primaryKey"`
	VoucherID string `gorm:"type:uuid;index"`
	SessionID string `gorm:"type:uuid;index"`
	Tel       *string
	Discount  *float64
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

type BoothLogModel struct {
	ID        string `gorm:"type:uuid;primaryKey"`
	BoothID   string `gorm:"type:uuid;index"`
	EventType string
	Level     string `gorm:"default:info"`
	Message   *string
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

type AnalyticsEventModel struct {
	ID        string  `gorm:"type:uuid;primaryKey"`
	BoothID   string  `gorm:"type:uuid;index"`
	SessionID *string `gorm:"type:uuid"`
	EventName string
	Payload   datatypes.JSONMap `gorm:"type:jsonb"`
	CreatedAt time.Time         `gorm:"autoCreateTime"`
}
