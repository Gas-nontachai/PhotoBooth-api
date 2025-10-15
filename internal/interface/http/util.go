package http

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func respondError(c *fiber.Ctx, err error) error {
	if err == nil {
		return nil
	}
	status := fiber.StatusBadRequest
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		status = fiber.StatusNotFound
	case errors.Is(err, fiber.ErrUnauthorized):
		status = fiber.StatusUnauthorized
	case errors.Is(err, fiber.ErrForbidden):
		status = fiber.StatusForbidden
	}
	return c.Status(status).JSON(fiber.Map{"error": err.Error()})
}

func respondSuccess(c *fiber.Ctx, status int, data any) error {
	if data == nil {
		return c.SendStatus(status)
	}
	return c.Status(status).JSON(data)
}

func parseBoolQuery(c *fiber.Ctx, key string) (bool, bool) {
	raw := c.Query(key, "")
	if raw == "" {
		return false, false
	}
	if raw == "true" || raw == "1" {
		return true, true
	}
	if raw == "false" || raw == "0" {
		return false, true
	}
	return false, false
}

func statusFromHTTP(status int) int {
	if status == 0 {
		return http.StatusOK
	}
	return status
}
