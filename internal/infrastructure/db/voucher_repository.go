package db

import (
	"context"

	"go-ddd-clean/internal/domain/voucher"

	"gorm.io/gorm"
)

type voucherRepository struct {
	db *gorm.DB
}

func NewVoucherRepository(db *gorm.DB) voucher.Repository {
	return &voucherRepository{db: db}
}

func (r *voucherRepository) Create(ctx context.Context, v *voucher.Voucher) error {
	model := VoucherModel{
		ID:        v.ID,
		Code:      v.Code,
		Type:      string(v.Type),
		Value:     v.Value,
		Unit:      string(v.Unit),
		MaxUsage:  v.MaxUsage,
		UsedCount: v.UsedCount,
		ValidFrom: v.ValidFrom,
		ValidTo:   v.ValidTo,
		Active:    v.Active,
	}
	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return err
	}
	v.CreatedAt = model.CreatedAt
	return nil
}

func (r *voucherRepository) Update(ctx context.Context, v *voucher.Voucher) error {
	return r.db.WithContext(ctx).
		Model(&VoucherModel{ID: v.ID}).
		Updates(map[string]any{
			"code":       v.Code,
			"type":       string(v.Type),
			"value":      v.Value,
			"unit":       string(v.Unit),
			"max_usage":  v.MaxUsage,
			"used_count": v.UsedCount,
			"valid_from": v.ValidFrom,
			"valid_to":   v.ValidTo,
			"active":     v.Active,
		}).Error
}

func (r *voucherRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&VoucherModel{ID: id}).Error
}

func (r *voucherRepository) GetByID(ctx context.Context, id string) (*voucher.Voucher, error) {
	var model VoucherModel
	if err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return mapVoucherModelToDomain(&model), nil
}

func (r *voucherRepository) GetByCode(ctx context.Context, code string) (*voucher.Voucher, error) {
	var model VoucherModel
	if err := r.db.WithContext(ctx).First(&model, "code = ?", code).Error; err != nil {
		return nil, err
	}
	return mapVoucherModelToDomain(&model), nil
}

func (r *voucherRepository) List(ctx context.Context, activeOnly bool) ([]voucher.Voucher, error) {
	query := r.db.WithContext(ctx)
	if activeOnly {
		query = query.Where("active = ?", true)
	}
	var models []VoucherModel
	if err := query.Find(&models).Error; err != nil {
		return nil, err
	}
	result := make([]voucher.Voucher, 0, len(models))
	for _, m := range models {
		result = append(result, *mapVoucherModelToDomain(&m))
	}
	return result, nil
}

type voucherRedemptionRepository struct {
	db *gorm.DB
}

func NewVoucherRedemptionRepository(db *gorm.DB) voucher.RedemptionRepository {
	return &voucherRedemptionRepository{db: db}
}

func (r *voucherRedemptionRepository) Create(ctx context.Context, red *voucher.Redemption) error {
	model := VoucherRedemptionModel{
		ID:        red.ID,
		VoucherID: red.VoucherID,
		SessionID: red.SessionID,
		Tel:       red.Tel,
		Discount:  red.Discount,
	}
	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return err
	}
	red.CreatedAt = model.CreatedAt
	return nil
}

func (r *voucherRedemptionRepository) ListByVoucher(ctx context.Context, voucherID string) ([]voucher.Redemption, error) {
	var models []VoucherRedemptionModel
	if err := r.db.WithContext(ctx).Where("voucher_id = ?", voucherID).Find(&models).Error; err != nil {
		return nil, err
	}
	result := make([]voucher.Redemption, 0, len(models))
	for _, m := range models {
		result = append(result, voucher.Redemption{
			ID:        m.ID,
			VoucherID: m.VoucherID,
			SessionID: m.SessionID,
			Tel:       m.Tel,
			Discount:  m.Discount,
			CreatedAt: m.CreatedAt,
		})
	}
	return result, nil
}

func mapVoucherModelToDomain(model *VoucherModel) *voucher.Voucher {
	if model == nil {
		return nil
	}
	return &voucher.Voucher{
		ID:        model.ID,
		Code:      model.Code,
		Type:      voucher.Type(model.Type),
		Value:     model.Value,
		Unit:      voucher.Unit(model.Unit),
		MaxUsage:  model.MaxUsage,
		UsedCount: model.UsedCount,
		ValidFrom: model.ValidFrom,
		ValidTo:   model.ValidTo,
		Active:    model.Active,
		CreatedAt: model.CreatedAt,
	}
}
