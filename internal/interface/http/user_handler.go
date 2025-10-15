package http

import (
	"context"

	appUser "go-ddd-clean/internal/application/user"
	domainUser "go-ddd-clean/internal/domain/user"

	"github.com/gofiber/fiber/v2"
)

type userHandler struct {
	service *appUser.Service
}

func newUserHandler(service *appUser.Service) *userHandler {
	return &userHandler{service: service}
}

func (h *userHandler) register(router fiber.Router) {
	router.Get("/", h.list)
	router.Post("/", h.create)
	router.Get("/:id", h.get)
	router.Put("/:id", h.update)
	router.Delete("/:id", h.delete)
	router.Post("/:id/points", h.adjustPoints)
}

func (h *userHandler) list(c *fiber.Ctx) error {
	result, err := h.service.List(context.Background())
	if err != nil {
		return respondError(c, err)
	}
	return respondSuccess(c, fiber.StatusOK, result)
}

func (h *userHandler) create(c *fiber.Ctx) error {
	var body struct {
		Tel      *string `json:"tel"`
		Email    *string `json:"email"`
		Password *string `json:"password"`
		Role     *string `json:"role"`
	}
	if err := c.BodyParser(&body); err != nil {
		return respondError(c, err)
	}
	role := domainUser.RoleCustomer
	if body.Role != nil {
		role = domainUser.Role(*body.Role)
	}
	entity, err := h.service.Create(context.Background(), appUser.CreateUserInput{
		Tel:      body.Tel,
		Email:    body.Email,
		Password: body.Password,
		Role:     role,
	})
	if err != nil {
		return respondError(c, err)
	}
	return respondSuccess(c, fiber.StatusCreated, entity)
}

func (h *userHandler) get(c *fiber.Ctx) error {
	id := c.Params("id")
	entity, err := h.service.Get(context.Background(), id)
	if err != nil {
		return respondError(c, err)
	}
	return respondSuccess(c, fiber.StatusOK, entity)
}

func (h *userHandler) update(c *fiber.Ctx) error {
	id := c.Params("id")
	var body struct {
		Tel      *string `json:"tel"`
		Email    *string `json:"email"`
		Password *string `json:"password"`
		Role     *string `json:"role"`
		Points   *int    `json:"points"`
	}
	if err := c.BodyParser(&body); err != nil {
		return respondError(c, err)
	}
	var rolePtr *domainUser.Role
	if body.Role != nil {
		r := domainUser.Role(*body.Role)
		rolePtr = &r
	}
	entity, err := h.service.Update(context.Background(), appUser.UpdateUserInput{
		ID:       id,
		Tel:      body.Tel,
		Email:    body.Email,
		Password: body.Password,
		Role:     rolePtr,
		Points:   body.Points,
	})
	if err != nil {
		return respondError(c, err)
	}
	return respondSuccess(c, fiber.StatusOK, entity)
}

func (h *userHandler) delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.service.Delete(context.Background(), id); err != nil {
		return respondError(c, err)
	}
	return respondSuccess(c, fiber.StatusNoContent, nil)
}

func (h *userHandler) adjustPoints(c *fiber.Ctx) error {
	id := c.Params("id")
	var body struct {
		Delta int `json:"delta"`
	}
	if err := c.BodyParser(&body); err != nil {
		return respondError(c, err)
	}
	entity, err := h.service.AdjustPoints(context.Background(), id, body.Delta)
	if err != nil {
		return respondError(c, err)
	}
	return respondSuccess(c, fiber.StatusOK, entity)
}
