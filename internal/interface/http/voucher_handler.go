package http

import (
	"context"
	"time"

	appSession "go-ddd-clean/internal/application/session"
	appVoucher "go-ddd-clean/internal/application/voucher"
	domainVoucher "go-ddd-clean/internal/domain/voucher"

	"github.com/gofiber/fiber/v2"
)

type voucherHandler struct {
	service        *appVoucher.Service
	sessionService *appSession.Service
}

func newVoucherHandler(service *appVoucher.Service, sessionService *appSession.Service) *voucherHandler {
	return &voucherHandler{
		service:        service,
		sessionService: sessionService,
	}
}

func (h *voucherHandler) register(router fiber.Router, boothAuth fiber.Handler) {
	router.Get("/", h.list)
	router.Post("/", h.create)
	router.Get("/:id", h.get)
	router.Put("/:id", h.update)
	router.Delete("/:id", h.delete)
	router.Get("/code/:code", h.getByCode)
	router.Post("/redeem", boothAuth, h.redeem)
}

func (h *voucherHandler) list(c *fiber.Ctx) error {
	active, ok := parseBoolQuery(c, "active")
	result, err := h.service.List(context.Background(), ok && active)
	if err != nil {
		return respondError(c, err)
	}
	return respondSuccess(c, fiber.StatusOK, result)
}

func (h *voucherHandler) create(c *fiber.Ctx) error {
	var body struct {
		Code      string  `json:"code"`
		Type      string  `json:"type"`
		Value     float64 `json:"value"`
		Unit      string  `json:"unit"`
		MaxUsage  int     `json:"max_usage"`
		ValidFrom *int64  `json:"valid_from"`
		ValidTo   *int64  `json:"valid_to"`
		Active    *bool   `json:"active"`
	}
	if err := c.BodyParser(&body); err != nil {
		return respondError(c, err)
	}
	if body.Code == "" || body.Type == "" || body.Unit == "" {
		return respondError(c, fiber.NewError(fiber.StatusBadRequest, "code, type and unit required"))
	}
	var (
		validFrom *time.Time
		validTo   *time.Time
	)
	if body.ValidFrom != nil {
		t := time.Unix(*body.ValidFrom, 0)
		validFrom = &t
	}
	if body.ValidTo != nil {
		t := time.Unix(*body.ValidTo, 0)
		validTo = &t
	}
	active := true
	if body.Active != nil {
		active = *body.Active
	}
	entity, err := h.service.Create(context.Background(), appVoucher.CreateVoucherInput{
		Code:      body.Code,
		Type:      domainVoucher.Type(body.Type),
		Value:     body.Value,
		Unit:      domainVoucher.Unit(body.Unit),
		MaxUsage:  body.MaxUsage,
		ValidFrom: validFrom,
		ValidTo:   validTo,
		Active:    active,
	})
	if err != nil {
		return respondError(c, err)
	}
	return respondSuccess(c, fiber.StatusCreated, entity)
}

func (h *voucherHandler) get(c *fiber.Ctx) error {
	id := c.Params("id")
	entity, err := h.service.Get(context.Background(), id)
	if err != nil {
		return respondError(c, err)
	}
	return respondSuccess(c, fiber.StatusOK, entity)
}

func (h *voucherHandler) update(c *fiber.Ctx) error {
	id := c.Params("id")
	var body struct {
		Value     *float64 `json:"value"`
		Unit      *string  `json:"unit"`
		MaxUsage  *int     `json:"max_usage"`
		ValidFrom *int64   `json:"valid_from"`
		ValidTo   *int64   `json:"valid_to"`
		Active    *bool    `json:"active"`
	}
	if err := c.BodyParser(&body); err != nil {
		return respondError(c, err)
	}
	var (
		unitPtr   *domainVoucher.Unit
		validFrom *time.Time
		validTo   *time.Time
	)
	if body.Unit != nil {
		unit := domainVoucher.Unit(*body.Unit)
		unitPtr = &unit
	}
	if body.ValidFrom != nil {
		t := time.Unix(*body.ValidFrom, 0)
		validFrom = &t
	}
	if body.ValidTo != nil {
		t := time.Unix(*body.ValidTo, 0)
		validTo = &t
	}
	entity, err := h.service.Update(context.Background(), appVoucher.UpdateVoucherInput{
		ID:        id,
		Value:     body.Value,
		Unit:      unitPtr,
		MaxUsage:  body.MaxUsage,
		ValidFrom: validFrom,
		ValidTo:   validTo,
		Active:    body.Active,
	})
	if err != nil {
		return respondError(c, err)
	}
	return respondSuccess(c, fiber.StatusOK, entity)
}

func (h *voucherHandler) delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.service.Delete(context.Background(), id); err != nil {
		return respondError(c, err)
	}
	return respondSuccess(c, fiber.StatusNoContent, nil)
}

func (h *voucherHandler) getByCode(c *fiber.Ctx) error {
	code := c.Params("code")
	entity, err := h.service.GetByCode(context.Background(), code)
	if err != nil {
		return respondError(c, err)
	}
	return respondSuccess(c, fiber.StatusOK, entity)
}

func (h *voucherHandler) redeem(c *fiber.Ctx) error {
	token, err := requireBoothToken(c)
	if err != nil {
		return respondError(c, err)
	}
	var body struct {
		Code      string   `json:"code"`
		SessionID string   `json:"session_id"`
		Tel       *string  `json:"tel"`
		Discount  *float64 `json:"discount"`
	}
	if err := c.BodyParser(&body); err != nil {
		return respondError(c, err)
	}
	if body.Code == "" || body.SessionID == "" {
		return respondError(c, fiber.NewError(fiber.StatusBadRequest, "code and session_id required"))
	}
	session, err := h.sessionService.Get(context.Background(), body.SessionID)
	if err != nil {
		return respondError(c, err)
	}
	if session.BoothID != token.BoothID {
		return respondError(c, fiber.ErrForbidden)
	}
	redemption, voucher, err := h.service.Redeem(context.Background(), appVoucher.RedeemVoucherInput{
		Code:      body.Code,
		SessionID: body.SessionID,
		Tel:       body.Tel,
		Discount:  body.Discount,
	})
	if err != nil {
		return respondError(c, err)
	}
	return respondSuccess(c, fiber.StatusOK, fiber.Map{
		"redemption": redemption,
		"voucher":    voucher,
	})
}
