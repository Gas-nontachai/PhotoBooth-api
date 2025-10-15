package http

import (
	"context"
	"time"

	appMedia "go-ddd-clean/internal/application/media"
	appPayment "go-ddd-clean/internal/application/payment"
	appSession "go-ddd-clean/internal/application/session"
	domainSession "go-ddd-clean/internal/domain/session"

	"github.com/gofiber/fiber/v2"
)

type sessionHandler struct {
	sessionService *appSession.Service
	photoService   *appMedia.PhotoService
	paymentService *appPayment.Service
}

func newSessionHandler(
	sessionService *appSession.Service,
	photoService *appMedia.PhotoService,
	paymentService *appPayment.Service,
) *sessionHandler {
	return &sessionHandler{
		sessionService: sessionService,
		photoService:   photoService,
		paymentService: paymentService,
	}
}

func (h *sessionHandler) register(router fiber.Router, boothAuth fiber.Handler) {
	protected := router.Group("/", boothAuth)
	protected.Get("/", h.list)
	protected.Post("/", h.create)
	protected.Get("/:id", h.get)
	protected.Put("/:id", h.update)
	protected.Delete("/:id", h.delete)

	protected.Get("/:id/photos", h.listPhotos)
	protected.Get("/:id/payment", h.getPayment)
}

func (h *sessionHandler) list(c *fiber.Ctx) error {
	token, err := requireBoothToken(c)
	if err != nil {
		return respondError(c, err)
	}
	if boothID := c.Query("booth_id", ""); boothID != "" && boothID != token.BoothID {
		return respondError(c, fiber.ErrForbidden)
	}
	var statusFilter *domainSession.Status
	if st := c.Query("status", ""); st != "" {
		status := domainSession.Status(st)
		statusFilter = &status
	}
	boothFilter := token.BoothID
	result, err := h.sessionService.List(context.Background(), &boothFilter, statusFilter)
	if err != nil {
		return respondError(c, err)
	}
	return respondSuccess(c, fiber.StatusOK, result)
}

func (h *sessionHandler) create(c *fiber.Ctx) error {
	token, err := requireBoothToken(c)
	if err != nil {
		return respondError(c, err)
	}
	var body struct {
		BoothID       string         `json:"booth_id"`
		UserID        *string        `json:"user_id"`
		VoucherID     *string        `json:"voucher_id"`
		PaymentID     *string        `json:"payment_id"`
		Status        *string        `json:"status"`
		TotalPrice    *float64       `json:"total_price"`
		BoothSnapshot map[string]any `json:"booth_snapshot"`
		PhoneTemp     *string        `json:"phone_temp"`
	}
	if err := c.BodyParser(&body); err != nil {
		return respondError(c, err)
	}
	if body.BoothID != "" && body.BoothID != token.BoothID {
		return respondError(c, fiber.ErrForbidden)
	}
	var status domainSession.Status
	if body.Status != nil {
		status = domainSession.Status(*body.Status)
	}
	entity, err := h.sessionService.Create(context.Background(), appSession.CreateSessionInput{
		BoothID:       token.BoothID,
		UserID:        body.UserID,
		VoucherID:     body.VoucherID,
		PaymentID:     body.PaymentID,
		Status:        status,
		TotalPrice:    body.TotalPrice,
		BoothSnapshot: body.BoothSnapshot,
		PhoneTemp:     body.PhoneTemp,
	})
	if err != nil {
		return respondError(c, err)
	}
	return respondSuccess(c, fiber.StatusCreated, entity)
}

func (h *sessionHandler) get(c *fiber.Ctx) error {
	token, err := requireBoothToken(c)
	if err != nil {
		return respondError(c, err)
	}
	id := c.Params("id")
	entity, err := h.sessionService.Get(context.Background(), id)
	if err != nil {
		return respondError(c, err)
	}
	if entity.BoothID != token.BoothID {
		return respondError(c, fiber.ErrForbidden)
	}
	return respondSuccess(c, fiber.StatusOK, entity)
}

func (h *sessionHandler) update(c *fiber.Ctx) error {
	token, err := requireBoothToken(c)
	if err != nil {
		return respondError(c, err)
	}
	id := c.Params("id")
	var body struct {
		UserID        *string        `json:"user_id"`
		VoucherID     *string        `json:"voucher_id"`
		PaymentID     *string        `json:"payment_id"`
		Status        *string        `json:"status"`
		TotalPrice    *float64       `json:"total_price"`
		FinishedAt    *int64         `json:"finished_at"`
		BoothSnapshot map[string]any `json:"booth_snapshot"`
		PhoneTemp     *string        `json:"phone_temp"`
	}
	if err := c.BodyParser(&body); err != nil {
		return respondError(c, err)
	}
	var (
		statusPtr     *domainSession.Status
		finishedAtPtr *time.Time
	)
	if body.Status != nil {
		status := domainSession.Status(*body.Status)
		statusPtr = &status
	}
	if body.FinishedAt != nil {
		ts := time.Unix(*body.FinishedAt, 0)
		finishedAtPtr = &ts
	}
	entity, err := h.sessionService.Update(context.Background(), appSession.UpdateSessionInput{
		ID:            id,
		UserID:        body.UserID,
		VoucherID:     body.VoucherID,
		PaymentID:     body.PaymentID,
		Status:        statusPtr,
		TotalPrice:    body.TotalPrice,
		FinishedAt:    finishedAtPtr,
		BoothSnapshot: body.BoothSnapshot,
		PhoneTemp:     body.PhoneTemp,
	})
	if err != nil {
		return respondError(c, err)
	}
	if entity.BoothID != token.BoothID {
		return respondError(c, fiber.ErrForbidden)
	}
	return respondSuccess(c, fiber.StatusOK, entity)
}

func (h *sessionHandler) delete(c *fiber.Ctx) error {
	token, err := requireBoothToken(c)
	if err != nil {
		return respondError(c, err)
	}
	id := c.Params("id")
	session, err := h.sessionService.Get(context.Background(), id)
	if err != nil {
		return respondError(c, err)
	}
	if session.BoothID != token.BoothID {
		return respondError(c, fiber.ErrForbidden)
	}
	if err := h.sessionService.Delete(context.Background(), id); err != nil {
		return respondError(c, err)
	}
	return respondSuccess(c, fiber.StatusNoContent, nil)
}

func (h *sessionHandler) listPhotos(c *fiber.Ctx) error {
	token, err := requireBoothToken(c)
	if err != nil {
		return respondError(c, err)
	}
	sessionID := c.Params("id")
	session, err := h.sessionService.Get(context.Background(), sessionID)
	if err != nil {
		return respondError(c, err)
	}
	if session.BoothID != token.BoothID {
		return respondError(c, fiber.ErrForbidden)
	}
	result, err := h.photoService.ListBySession(context.Background(), sessionID)
	if err != nil {
		return respondError(c, err)
	}
	return respondSuccess(c, fiber.StatusOK, result)
}

func (h *sessionHandler) getPayment(c *fiber.Ctx) error {
	token, err := requireBoothToken(c)
	if err != nil {
		return respondError(c, err)
	}
	sessionID := c.Params("id")
	session, err := h.sessionService.Get(context.Background(), sessionID)
	if err != nil {
		return respondError(c, err)
	}
	if session.BoothID != token.BoothID {
		return respondError(c, fiber.ErrForbidden)
	}
	result, err := h.paymentService.GetBySession(context.Background(), sessionID)
	if err != nil {
		return respondError(c, err)
	}
	return respondSuccess(c, fiber.StatusOK, result)
}
