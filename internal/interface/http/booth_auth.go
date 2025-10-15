package http

import (
	"context"
	"errors"
	"strings"

	appBooth "go-ddd-clean/internal/application/booth"

	"github.com/gofiber/fiber/v2"
)

const boothTokenContextKey = "booth_token"

func newBoothAuthMiddleware(tokenService *appBooth.TokenService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		header := c.Get("Authorization")
		if header == "" {
			return fiber.ErrUnauthorized
		}
		token, err := parseBearerToken(header)
		if err != nil {
			return fiber.ErrUnauthorized
		}
		validated, err := tokenService.Validate(context.Background(), token)
		if err != nil {
			if errors.Is(err, appBooth.ErrInvalidBoothToken) || errors.Is(err, appBooth.ErrTokenMismatch) {
				return fiber.NewError(fiber.StatusUnauthorized, err.Error())
			}
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		c.Locals(boothTokenContextKey, validated)
		return c.Next()
	}
}

func parseBearerToken(header string) (string, error) {
	if header == "" {
		return "", fiber.ErrUnauthorized
	}
	if !strings.HasPrefix(header, "Bearer ") {
		return "", fiber.ErrUnauthorized
	}
	token := strings.TrimSpace(strings.TrimPrefix(header, "Bearer "))
	if token == "" {
		return "", fiber.ErrUnauthorized
	}
	return token, nil
}

func requireBoothToken(c *fiber.Ctx) (*appBooth.ValidatedToken, error) {
	raw := c.Locals(boothTokenContextKey)
	if raw == nil {
		return nil, fiber.ErrUnauthorized
	}
	token, ok := raw.(*appBooth.ValidatedToken)
	if !ok || token == nil {
		return nil, fiber.ErrUnauthorized
	}
	return token, nil
}
