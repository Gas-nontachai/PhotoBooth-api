package http

import (
	"context"

	appBranch "go-ddd-clean/internal/application/branch"

	"github.com/gofiber/fiber/v2"
)

type branchHandler struct {
	service *appBranch.Service
}

func newBranchHandler(service *appBranch.Service) *branchHandler {
	return &branchHandler{service: service}
}

func (h *branchHandler) register(router fiber.Router) {
	router.Get("/", h.list)
	router.Post("/", h.create)
	router.Get("/:id", h.get)
	router.Put("/:id", h.update)
	router.Delete("/:id", h.delete)
}

func (h *branchHandler) list(c *fiber.Ctx) error {
	result, err := h.service.List(context.Background())
	if err != nil {
		return respondError(c, err)
	}
	return respondSuccess(c, fiber.StatusOK, result)
}

func (h *branchHandler) create(c *fiber.Ctx) error {
	var body struct {
		Name     string  `json:"name"`
		Location *string `json:"location"`
	}
	if err := c.BodyParser(&body); err != nil {
		return respondError(c, err)
	}
	if body.Name == "" {
		return respondError(c, fiber.NewError(fiber.StatusBadRequest, "name is required"))
	}
	branch, err := h.service.Create(context.Background(), appBranch.CreateBranchInput{
		Name:     body.Name,
		Location: body.Location,
	})
	if err != nil {
		return respondError(c, err)
	}
	return respondSuccess(c, fiber.StatusCreated, branch)
}

func (h *branchHandler) get(c *fiber.Ctx) error {
	id := c.Params("id")
	result, err := h.service.Get(context.Background(), id)
	if err != nil {
		return respondError(c, err)
	}
	return respondSuccess(c, fiber.StatusOK, result)
}

func (h *branchHandler) update(c *fiber.Ctx) error {
	id := c.Params("id")
	var body struct {
		Name     string  `json:"name"`
		Location *string `json:"location"`
	}
	if err := c.BodyParser(&body); err != nil {
		return respondError(c, err)
	}
	if body.Name == "" {
		return respondError(c, fiber.NewError(fiber.StatusBadRequest, "name is required"))
	}
	err := h.service.Update(context.Background(), appBranch.UpdateBranchInput{
		ID:       id,
		Name:     body.Name,
		Location: body.Location,
	})
	if err != nil {
		return respondError(c, err)
	}
	result, err := h.service.Get(context.Background(), id)
	if err != nil {
		return respondError(c, err)
	}
	return respondSuccess(c, fiber.StatusOK, result)
}

func (h *branchHandler) delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.service.Delete(context.Background(), id); err != nil {
		return respondError(c, err)
	}
	return respondSuccess(c, fiber.StatusNoContent, nil)
}
