package seeder

import (
	"context"
	"errors"
	"time"

	"go-ddd-clean/internal/infrastructure/db"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func Run(ctx context.Context, database *gorm.DB) error {
	return database.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		branchIDs, err := seedBranches(tx)
		if err != nil {
			return err
		}
		if err := seedBooths(tx, branchIDs); err != nil {
			return err
		}
		if err := seedFrames(tx); err != nil {
			return err
		}
		if err := seedFilters(tx); err != nil {
			return err
		}
		if err := seedUsers(tx); err != nil {
			return err
		}
		if err := seedVouchers(tx); err != nil {
			return err
		}
		return nil
	})
}

func seedBranches(tx *gorm.DB) (map[string]string, error) {
	definitions := []struct {
		Name     string
		Location string
	}{
		{Name: "Central Plaza", Location: "Bangkok, Thailand"},
		{Name: "Phuket Marina", Location: "Phuket, Thailand"},
	}

	ids := make(map[string]string, len(definitions))
	for _, def := range definitions {
		var model db.BranchModel
		err := tx.Where("name = ?", def.Name).First(&model).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			model = db.BranchModel{
				ID:       uuid.NewString(),
				Name:     def.Name,
				Location: optionalString(def.Location),
			}
			if err := tx.Create(&model).Error; err != nil {
				return nil, err
			}
		} else if err != nil {
			return nil, err
		}
		ids[def.Name] = model.ID
	}
	return ids, nil
}

func seedBooths(tx *gorm.DB, branchIDs map[string]string) error {
	definitions := []struct {
		Name       string
		BranchName string
		Type       string
		Status     string
		Config     datatypes.JSONMap
	}{
		{
			Name:       "Central Lobby Booth",
			BranchName: "Central Plaza",
			Type:       "physical",
			Status:     "active",
			Config: datatypes.JSONMap{
				"camera":     "sony-a7",
				"background": "neon",
			},
		},
		{
			Name:       "Phuket Virtual Booth",
			BranchName: "Phuket Marina",
			Type:       "virtual",
			Status:     "inactive",
			Config: datatypes.JSONMap{
				"background": "sunset",
				"overlay":    "tropical",
			},
		},
	}

	for _, def := range definitions {
		var model db.BoothModel
		err := tx.Where("name = ?", def.Name).First(&model).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			branchID := branchIDs[def.BranchName]
			if branchID == "" {
				return errors.New("missing branch for booth seed: " + def.BranchName)
			}
			model = db.BoothModel{
				ID:       uuid.NewString(),
				BranchID: branchID,
				Name:     def.Name,
				Type:     def.Type,
				Status:   def.Status,
				Config:   def.Config,
			}
			if err := tx.Create(&model).Error; err != nil {
				return err
			}
		} else if err != nil {
			return err
		}
	}
	return nil
}

func seedFrames(tx *gorm.DB) error {
	definitions := []struct {
		Name    string
		Theme   string
		FileURL string
		Active  bool
	}{
		{
			Name:    "Classic Gold",
			Theme:   "classic",
			FileURL: "https://cdn.example.com/frames/classic-gold.png",
			Active:  true,
		},
		{
			Name:    "Festival Neon",
			Theme:   "festival",
			FileURL: "https://cdn.example.com/frames/festival-neon.png",
			Active:  true,
		},
	}

	for _, def := range definitions {
		var model db.FrameModel
		err := tx.Where("name = ?", def.Name).First(&model).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			model = db.FrameModel{
				ID:      uuid.NewString(),
				Name:    def.Name,
				Theme:   optionalString(def.Theme),
				FileURL: def.FileURL,
				Active:  def.Active,
			}
			if err := tx.Create(&model).Error; err != nil {
				return err
			}
		} else if err != nil {
			return err
		}
	}
	return nil
}

func seedFilters(tx *gorm.DB) error {
	definitions := []struct {
		Name   string
		Effect datatypes.JSONMap
		Active bool
	}{
		{
			Name: "Warm Glow",
			Effect: datatypes.JSONMap{
				"brightness":  10,
				"saturation":  8,
				"temperature": 4,
			},
			Active: true,
		},
		{
			Name: "Black & White",
			Effect: datatypes.JSONMap{
				"desaturate": true,
				"contrast":   6,
			},
			Active: true,
		},
	}

	for _, def := range definitions {
		var model db.FilterModel
		err := tx.Where("name = ?", def.Name).First(&model).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			model = db.FilterModel{
				ID:     uuid.NewString(),
				Name:   def.Name,
				Effect: def.Effect,
				Active: def.Active,
			}
			if err := tx.Create(&model).Error; err != nil {
				return err
			}
		} else if err != nil {
			return err
		}
	}
	return nil
}

func seedUsers(tx *gorm.DB) error {
	definitions := []struct {
		Email  *string
		Tel    *string
		Role   string
		Points int
	}{
		{
			Email:  optionalString("customer@example.com"),
			Tel:    optionalString("+66999999999"),
			Role:   "customer",
			Points: 120,
		},
		{
			Email:  optionalString("admin@example.com"),
			Role:   "admin",
			Points: 0,
		},
	}

	for _, def := range definitions {
		query := tx.Model(&db.UserModel{})
		switch {
		case def.Email != nil:
			query = query.Where("email = ?", *def.Email)
		case def.Tel != nil:
			query = query.Where("tel = ?", *def.Tel)
		}

		var model db.UserModel
		err := query.First(&model).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			model = db.UserModel{
				ID:     uuid.NewString(),
				Email:  def.Email,
				Tel:    def.Tel,
				Role:   def.Role,
				Points: def.Points,
			}
			if err := tx.Create(&model).Error; err != nil {
				return err
			}
		} else if err != nil {
			return err
		}
	}
	return nil
}

func seedVouchers(tx *gorm.DB) error {
	now := time.Now()
	nextMonth := now.AddDate(0, 1, 0)

	definitions := []struct {
		Code      string
		Type      string
		Value     float64
		Unit      string
		MaxUsage  int
		ValidFrom *time.Time
		ValidTo   *time.Time
		Active    bool
	}{
		{
			Code:      "WELCOME50",
			Type:      "discount",
			Value:     50,
			Unit:      "baht",
			MaxUsage:  100,
			ValidFrom: &now,
			ValidTo:   &nextMonth,
			Active:    true,
		},
		{
			Code:     "FREESESSION",
			Type:     "free",
			Value:    1,
			Unit:     "session",
			MaxUsage: 50,
			Active:   true,
		},
	}

	for _, def := range definitions {
		var model db.VoucherModel
		err := tx.Where("code = ?", def.Code).First(&model).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			model = db.VoucherModel{
				ID:        uuid.NewString(),
				Code:      def.Code,
				Type:      def.Type,
				Value:     def.Value,
				Unit:      def.Unit,
				MaxUsage:  def.MaxUsage,
				ValidFrom: def.ValidFrom,
				ValidTo:   def.ValidTo,
				Active:    def.Active,
			}
			if err := tx.Create(&model).Error; err != nil {
				return err
			}
		} else if err != nil {
			return err
		}
	}
	return nil
}

func optionalString(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}
