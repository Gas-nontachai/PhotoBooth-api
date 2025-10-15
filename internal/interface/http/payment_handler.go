package http

import (
	"context"

	appPayment "go-ddd-clean/internal/application/payment"
	domainPayment "go-ddd-clean/internal/domain/payment"

	"github.com/gofiber/fiber/v2"
)

type paymentHandler struct {
	service *appPayment.Service
}

func newPaymentHandler(service *appPayment.Service) *paymentHandler {
	return &paymentHandler{service: service}
}

func (h *paymentHandler) register(router fiber.Router) {
	router.Post("/", h.create)
	router.Get("/:id", h.get)
	router.Put("/:id", h.update)
	router.Get("/session/:sessionID", h.getBySession)
}

func (h *paymentHandler) create(c *fiber.Ctx) error {
	var body struct {
		SessionID      string  `json:"session_id"`
		Method         string  `json:"method"`
		Amount         float64 `json:"amount"`
		Currency       string  `json:"currency"`
		Status         string  `json:"status"`
		TransactionRef *string `json:"transaction_ref"`
	}
	if err := c.BodyParser(&body); err != nil {
		return respondError(c, err)
	}
	if body.SessionID == "" || body.Method == "" {
		return respondError(c, fiber.NewError(fiber.StatusBadRequest, "session_id and method required"))
	}
	entity, err := h.service.Create(context.Background(), appPayment.CreatePaymentInput{
		SessionID:      body.SessionID,
		Method:         domainPayment.Method(body.Method),
		Amount:         body.Amount,
		Currency:       body.Currency,
		Status:         domainPayment.Status(body.Status),
		TransactionRef: body.TransactionRef,
	})
	if err != nil {
		return respondError(c, err)
	}
	return respondSuccess(c, fiber.StatusCreated, entity)
}

func (h *paymentHandler) get(c *fiber.Ctx) error {
	id := c.Params("id")
	entity, err := h.service.Get(context.Background(), id)
	if err != nil {
		return respondError(c, err)
	}
	return respondSuccess(c, fiber.StatusOK, entity)
}

func (h *paymentHandler) update(c *fiber.Ctx) error {
	id := c.Params("id")
	var body struct {
		Status         *string  `json:"status"`
		TransactionRef *string  `json:"transaction_ref"`
		Amount         *float64 `json:"amount"`
		Currency       *string  `json:"currency"`
		Method         *string  `json:"method"`
	}
	if err := c.BodyParser(&body); err != nil {
		return respondError(c, err)
	}
	var status domainPayment.Status
	if body.Status != nil {
		status = domainPayment.Status(*body.Status)
	}
	var methodPtr *domainPayment.Method
	if body.Method != nil {
		method := domainPayment.Method(*body.Method)
		methodPtr = &method
	}
	entity, err := h.service.Update(context.Background(), appPayment.UpdatePaymentStatusInput{
		ID:             id,
		Status:         status,
		TransactionRef: body.TransactionRef,
		Amount:         body.Amount,
		Currency:       body.Currency,
		Method:         methodPtr,
	})
	if err != nil {
		return respondError(c, err)
	}
	return respondSuccess(c, fiber.StatusOK, entity)
}

func (h *paymentHandler) getBySession(c *fiber.Ctx) error {
	sessionID := c.Params("sessionID")
	entity, err := h.service.GetBySession(context.Background(), sessionID)
	if err != nil {
		return respondError(c, err)
	}
	return respondSuccess(c, fiber.StatusOK, entity)
}
