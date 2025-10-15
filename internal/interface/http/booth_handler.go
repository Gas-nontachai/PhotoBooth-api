package http

import (
	"context"
	"strconv"

	appAnalytics "go-ddd-clean/internal/application/analytics"
	appBooth "go-ddd-clean/internal/application/booth"
	appLogging "go-ddd-clean/internal/application/logging"
	domainBooth "go-ddd-clean/internal/domain/booth"
	domainLogging "go-ddd-clean/internal/domain/logging"

	"github.com/gofiber/fiber/v2"
)

type boothHandler struct {
	boothService     *appBooth.Service
	loggingService   *appLogging.Service
	analyticsService *appAnalytics.Service
}

func newBoothHandler(
	boothService *appBooth.Service,
	loggingService *appLogging.Service,
	analyticsService *appAnalytics.Service,
) *boothHandler {
	return &boothHandler{
		boothService:     boothService,
		loggingService:   loggingService,
		analyticsService: analyticsService,
	}
}

func (h *boothHandler) register(router fiber.Router) {
	router.Get("/", h.list)
	router.Post("/", h.create)
	router.Get("/:id", h.get)
	router.Put("/:id", h.update)
	router.Delete("/:id", h.delete)

	router.Get("/:id/logs", h.listLogs)
	router.Post("/:id/logs", h.createLog)

	router.Get("/:id/analytics", h.listAnalytics)
	router.Post("/:id/analytics", h.createAnalyticsEvent)
}

func (h *boothHandler) list(c *fiber.Ctx) error {
	branchID := c.Query("branch_id", "")
	var filter *string
	if branchID != "" {
		filter = &branchID
	}
	result, err := h.boothService.List(context.Background(), filter)
	if err != nil {
		return respondError(c, err)
	}
	return respondSuccess(c, fiber.StatusOK, result)
}

func (h *boothHandler) create(c *fiber.Ctx) error {
	var body struct {
		BranchID string         `json:"branch_id"`
		Name     string         `json:"name"`
		Type     string         `json:"type"`
		Status   string         `json:"status"`
		Config   map[string]any `json:"config"`
	}
	if err := c.BodyParser(&body); err != nil {
		return respondError(c, err)
	}
	if body.BranchID == "" || body.Name == "" {
		return respondError(c, fiber.NewError(fiber.StatusBadRequest, "branch_id and name are required"))
	}
	boothType := domainBooth.BoothType(body.Type)
	if boothType == "" {
		boothType = domainBooth.BoothTypePhysical
	}
	status := domainBooth.BoothStatus(body.Status)
	entity, err := h.boothService.Create(context.Background(), appBooth.CreateBoothInput{
		BranchID: body.BranchID,
		Name:     body.Name,
		Type:     boothType,
		Status:   status,
		Config:   body.Config,
	})
	if err != nil {
		return respondError(c, err)
	}
	return respondSuccess(c, fiber.StatusCreated, entity)
}

func (h *boothHandler) get(c *fiber.Ctx) error {
	id := c.Params("id")
	entity, err := h.boothService.Get(context.Background(), id)
	if err != nil {
		return respondError(c, err)
	}
	return respondSuccess(c, fiber.StatusOK, entity)
}

func (h *boothHandler) update(c *fiber.Ctx) error {
	id := c.Params("id")
	var body struct {
		BranchID string         `json:"branch_id"`
		Name     string         `json:"name"`
		Type     *string        `json:"type"`
		Status   *string        `json:"status"`
		Config   map[string]any `json:"config"`
	}
	if err := c.BodyParser(&body); err != nil {
		return respondError(c, err)
	}
	currentType := domainBooth.BoothType("")
	if body.Type != nil {
		currentType = domainBooth.BoothType(*body.Type)
	}
	currentStatus := domainBooth.BoothStatus("")
	if body.Status != nil {
		currentStatus = domainBooth.BoothStatus(*body.Status)
	}
	err := h.boothService.Update(context.Background(), appBooth.UpdateBoothInput{
		ID:       id,
		BranchID: body.BranchID,
		Name:     body.Name,
		Type:     currentType,
		Status:   currentStatus,
		Config:   body.Config,
	})
	if err != nil {
		return respondError(c, err)
	}
	entity, err := h.boothService.Get(context.Background(), id)
	if err != nil {
		return respondError(c, err)
	}
	return respondSuccess(c, fiber.StatusOK, entity)
}

func (h *boothHandler) delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.boothService.Delete(context.Background(), id); err != nil {
		return respondError(c, err)
	}
	return respondSuccess(c, fiber.StatusNoContent, nil)
}

func (h *boothHandler) listLogs(c *fiber.Ctx) error {
	boothID := c.Params("id")
	limitRaw := c.Query("limit", "")
	limit := 50
	if limitRaw != "" {
		if parsed, err := strconv.Atoi(limitRaw); err == nil && parsed > 0 {
			limit = parsed
		}
	}
	result, err := h.loggingService.List(context.Background(), boothID, limit)
	if err != nil {
		return respondError(c, err)
	}
	return respondSuccess(c, fiber.StatusOK, result)
}

func (h *boothHandler) createLog(c *fiber.Ctx) error {
	boothID := c.Params("id")
	var body struct {
		EventType string  `json:"event_type"`
		Level     string  `json:"level"`
		Message   *string `json:"message"`
	}
	if err := c.BodyParser(&body); err != nil {
		return respondError(c, err)
	}
	if body.EventType == "" {
		return respondError(c, fiber.NewError(fiber.StatusBadRequest, "event_type required"))
	}
	entry, err := h.loggingService.Write(context.Background(), appLogging.CreateLogInput{
		BoothID:   boothID,
		EventType: body.EventType,
		Level:     parseLogLevel(body.Level),
		Message:   body.Message,
	})
	if err != nil {
		return respondError(c, err)
	}
	return respondSuccess(c, fiber.StatusCreated, entry)
}

func (h *boothHandler) listAnalytics(c *fiber.Ctx) error {
	boothID := c.Params("id")
	limit := 100
	if raw := c.Query("limit", ""); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 {
			limit = parsed
		}
	}
	result, err := h.analyticsService.ListByBooth(context.Background(), boothID, limit)
	if err != nil {
		return respondError(c, err)
	}
	return respondSuccess(c, fiber.StatusOK, result)
}

func (h *boothHandler) createAnalyticsEvent(c *fiber.Ctx) error {
	boothID := c.Params("id")
	var body struct {
		SessionID *string        `json:"session_id"`
		EventName string         `json:"event_name"`
		Payload   map[string]any `json:"payload"`
	}
	if err := c.BodyParser(&body); err != nil {
		return respondError(c, err)
	}
	if body.EventName == "" {
		return respondError(c, fiber.NewError(fiber.StatusBadRequest, "event_name required"))
	}
	event, err := h.analyticsService.Create(context.Background(), appAnalytics.CreateEventInput{
		BoothID:   boothID,
		SessionID: body.SessionID,
		EventName: body.EventName,
		Payload:   body.Payload,
	})
	if err != nil {
		return respondError(c, err)
	}
	return respondSuccess(c, fiber.StatusCreated, event)
}

func parseLogLevel(level string) domainLogging.Level {
	switch level {
	case string(domainLogging.LevelWarn):
		return domainLogging.LevelWarn
	case string(domainLogging.LevelError):
		return domainLogging.LevelError
	default:
		return domainLogging.LevelInfo
	}
}
