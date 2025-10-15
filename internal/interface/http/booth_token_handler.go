package http

import (
	"context"
	"errors"

	appBooth "go-ddd-clean/internal/application/booth"

	"github.com/gofiber/fiber/v2"
)

type boothTokenHandler struct {
	tokenService *appBooth.TokenService
}

func newBoothTokenHandler(tokenService *appBooth.TokenService) *boothTokenHandler {
	return &boothTokenHandler{
		tokenService: tokenService,
	}
}

func (h *boothTokenHandler) register(c *fiber.Ctx) error {
	var body struct {
		BoothID  string `json:"booth_id"`
		BranchID string `json:"branch_id"`
	}
	if err := c.BodyParser(&body); err != nil {
		return respondError(c, err)
	}
	if body.BoothID == "" || body.BranchID == "" {
		return respondError(c, fiber.NewError(fiber.StatusBadRequest, "booth_id and branch_id are required"))
	}
	token, err := h.tokenService.Register(context.Background(), body.BoothID, body.BranchID)
	if err != nil {
		if errors.Is(err, appBooth.ErrInvalidBoothToken) || errors.Is(err, appBooth.ErrTokenMismatch) {
			return respondError(c, fiber.NewError(fiber.StatusUnauthorized, err.Error()))
		}
		return respondError(c, err)
	}
	return respondSuccess(c, fiber.StatusOK, fiber.Map{
		"token": token,
	})
}

func (h *boothTokenHandler) regenerate(c *fiber.Ctx) error {
	tokenInfo, err := requireBoothToken(c)
	if err != nil {
		return respondError(c, err)
	}
	token, err := h.tokenService.Regenerate(context.Background(), tokenInfo.BoothID)
	if err != nil {
		return respondError(c, err)
	}
	return respondSuccess(c, fiber.StatusOK, fiber.Map{
		"token": token,
	})
}
