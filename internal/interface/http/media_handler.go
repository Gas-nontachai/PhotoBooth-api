package http

import (
	"context"
	"time"

	appMedia "go-ddd-clean/internal/application/media"
	appSession "go-ddd-clean/internal/application/session"
	domainMedia "go-ddd-clean/internal/domain/media"

	"github.com/gofiber/fiber/v2"
)

type mediaHandler struct {
	sessionService *appSession.Service
	photoService   *appMedia.PhotoService
	frameService   *appMedia.FrameService
	filterService  *appMedia.FilterService
	qrService      *appMedia.QRCodeService
}

func newMediaHandler(
	sessionService *appSession.Service,
	photoService *appMedia.PhotoService,
	frameService *appMedia.FrameService,
	filterService *appMedia.FilterService,
	qrService *appMedia.QRCodeService,
) *mediaHandler {
	return &mediaHandler{
		sessionService: sessionService,
		photoService:   photoService,
		frameService:   frameService,
		filterService:  filterService,
		qrService:      qrService,
	}
}

func (h *mediaHandler) register(router fiber.Router, boothAuth fiber.Handler) {
	photos := router.Group("/photos", boothAuth)
	photos.Get("/", h.listPhotos)
	photos.Post("/", h.createPhoto)
	photos.Get("/:id", h.getPhoto)
	photos.Put("/:id", h.updatePhoto)
	photos.Delete("/:id", h.deletePhoto)

	frames := router.Group("/frames")
	frames.Get("/", h.listFrames)
	frames.Post("/", h.createFrame)
	frames.Get("/:id", h.getFrame)
	frames.Put("/:id", h.updateFrame)
	frames.Delete("/:id", h.deleteFrame)

	filters := router.Group("/filters")
	filters.Get("/", h.listFilters)
	filters.Post("/", h.createFilter)
	filters.Get("/:id", h.getFilter)
	filters.Put("/:id", h.updateFilter)
	filters.Delete("/:id", h.deleteFilter)

	qrcodes := router.Group("/qrcodes")
	qrcodes.Post("/", h.createQRCode)
	qrcodes.Get("/:hash", h.getQRCode)
	qrcodes.Delete("/:id", h.deleteQRCode)
}

func (h *mediaHandler) listPhotos(c *fiber.Ctx) error {
	token, err := requireBoothToken(c)
	if err != nil {
		return respondError(c, err)
	}
	sessionID := c.Query("session_id", "")
	if sessionID == "" {
		return respondError(c, fiber.NewError(fiber.StatusBadRequest, "session_id query param required"))
	}
	if err := h.ensureSessionBelongs(context.Background(), sessionID, token.BoothID); err != nil {
		return respondError(c, err)
	}
	result, err := h.photoService.ListBySession(context.Background(), sessionID)
	if err != nil {
		return respondError(c, err)
	}
	return respondSuccess(c, fiber.StatusOK, result)
}

func (h *mediaHandler) createPhoto(c *fiber.Ctx) error {
	token, err := requireBoothToken(c)
	if err != nil {
		return respondError(c, err)
	}
	var body struct {
		SessionID   string         `json:"session_id"`
		FrameID     *string        `json:"frame_id"`
		FilterID    *string        `json:"filter_id"`
		StorageURL  string         `json:"storage_url"`
		Composition map[string]any `json:"composition"`
		RenderedURL *string        `json:"rendered_url"`
	}
	if err := c.BodyParser(&body); err != nil {
		return respondError(c, err)
	}
	if body.SessionID == "" || body.StorageURL == "" {
		return respondError(c, fiber.NewError(fiber.StatusBadRequest, "session_id and storage_url are required"))
	}
	if err := h.ensureSessionBelongs(context.Background(), body.SessionID, token.BoothID); err != nil {
		return respondError(c, err)
	}
	entity, err := h.photoService.Create(context.Background(), appMedia.CreatePhotoInput{
		SessionID:   body.SessionID,
		FrameID:     body.FrameID,
		FilterID:    body.FilterID,
		StorageURL:  body.StorageURL,
		Composition: body.Composition,
		RenderedURL: body.RenderedURL,
	})
	if err != nil {
		return respondError(c, err)
	}
	return respondSuccess(c, fiber.StatusCreated, entity)
}

func (h *mediaHandler) getPhoto(c *fiber.Ctx) error {
	token, err := requireBoothToken(c)
	if err != nil {
		return respondError(c, err)
	}
	id := c.Params("id")
	entity, err := h.ensurePhotoBelongs(context.Background(), id, token.BoothID)
	if err != nil {
		return respondError(c, err)
	}
	return respondSuccess(c, fiber.StatusOK, entity)
}

func (h *mediaHandler) updatePhoto(c *fiber.Ctx) error {
	token, err := requireBoothToken(c)
	if err != nil {
		return respondError(c, err)
	}
	id := c.Params("id")
	var body struct {
		FrameID     *string        `json:"frame_id"`
		FilterID    *string        `json:"filter_id"`
		StorageURL  *string        `json:"storage_url"`
		Composition map[string]any `json:"composition"`
		RenderedURL *string        `json:"rendered_url"`
	}
	if err := c.BodyParser(&body); err != nil {
		return respondError(c, err)
	}
	if _, err := h.ensurePhotoBelongs(context.Background(), id, token.BoothID); err != nil {
		return respondError(c, err)
	}
	entity, err := h.photoService.Update(context.Background(), appMedia.UpdatePhotoInput{
		ID:          id,
		FrameID:     body.FrameID,
		FilterID:    body.FilterID,
		StorageURL:  body.StorageURL,
		Composition: body.Composition,
		RenderedURL: body.RenderedURL,
	})
	if err != nil {
		return respondError(c, err)
	}
	return respondSuccess(c, fiber.StatusOK, entity)
}

func (h *mediaHandler) deletePhoto(c *fiber.Ctx) error {
	token, err := requireBoothToken(c)
	if err != nil {
		return respondError(c, err)
	}
	id := c.Params("id")
	if _, err := h.ensurePhotoBelongs(context.Background(), id, token.BoothID); err != nil {
		return respondError(c, err)
	}
	if err := h.photoService.Delete(context.Background(), id); err != nil {
		return respondError(c, err)
	}
	return respondSuccess(c, fiber.StatusNoContent, nil)
}

func (h *mediaHandler) listFrames(c *fiber.Ctx) error {
	active, ok := parseBoolQuery(c, "active")
	result, err := h.frameService.List(context.Background(), ok && active)
	if err != nil {
		return respondError(c, err)
	}
	return respondSuccess(c, fiber.StatusOK, result)
}

func (h *mediaHandler) createFrame(c *fiber.Ctx) error {
	var body struct {
		Name    string  `json:"name"`
		Theme   *string `json:"theme"`
		FileURL string  `json:"file_url"`
		Active  *bool   `json:"active"`
	}
	if err := c.BodyParser(&body); err != nil {
		return respondError(c, err)
	}
	if body.Name == "" || body.FileURL == "" {
		return respondError(c, fiber.NewError(fiber.StatusBadRequest, "name and file_url required"))
	}
	active := true
	if body.Active != nil {
		active = *body.Active
	}
	entity, err := h.frameService.Create(context.Background(), appMedia.CreateFrameInput{
		Name:    body.Name,
		Theme:   body.Theme,
		FileURL: body.FileURL,
		Active:  active,
	})
	if err != nil {
		return respondError(c, err)
	}
	return respondSuccess(c, fiber.StatusCreated, entity)
}

func (h *mediaHandler) getFrame(c *fiber.Ctx) error {
	id := c.Params("id")
	entity, err := h.frameService.Get(context.Background(), id)
	if err != nil {
		return respondError(c, err)
	}
	return respondSuccess(c, fiber.StatusOK, entity)
}

func (h *mediaHandler) updateFrame(c *fiber.Ctx) error {
	id := c.Params("id")
	var body struct {
		Name    string  `json:"name"`
		Theme   *string `json:"theme"`
		FileURL string  `json:"file_url"`
		Active  bool    `json:"active"`
	}
	if err := c.BodyParser(&body); err != nil {
		return respondError(c, err)
	}
	entity, err := h.frameService.Update(context.Background(), appMedia.UpdateFrameInput{
		ID:      id,
		Name:    body.Name,
		Theme:   body.Theme,
		FileURL: body.FileURL,
		Active:  body.Active,
	})
	if err != nil {
		return respondError(c, err)
	}
	return respondSuccess(c, fiber.StatusOK, entity)
}

func (h *mediaHandler) deleteFrame(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.frameService.Delete(context.Background(), id); err != nil {
		return respondError(c, err)
	}
	return respondSuccess(c, fiber.StatusNoContent, nil)
}

func (h *mediaHandler) listFilters(c *fiber.Ctx) error {
	active, ok := parseBoolQuery(c, "active")
	result, err := h.filterService.List(context.Background(), ok && active)
	if err != nil {
		return respondError(c, err)
	}
	return respondSuccess(c, fiber.StatusOK, result)
}

func (h *mediaHandler) createFilter(c *fiber.Ctx) error {
	var body struct {
		Name   string         `json:"name"`
		Effect map[string]any `json:"effect"`
		Active *bool          `json:"active"`
	}
	if err := c.BodyParser(&body); err != nil {
		return respondError(c, err)
	}
	if body.Name == "" {
		return respondError(c, fiber.NewError(fiber.StatusBadRequest, "name required"))
	}
	active := true
	if body.Active != nil {
		active = *body.Active
	}
	entity, err := h.filterService.Create(context.Background(), appMedia.CreateFilterInput{
		Name:   body.Name,
		Effect: body.Effect,
		Active: active,
	})
	if err != nil {
		return respondError(c, err)
	}
	return respondSuccess(c, fiber.StatusCreated, entity)
}

func (h *mediaHandler) getFilter(c *fiber.Ctx) error {
	id := c.Params("id")
	entity, err := h.filterService.Get(context.Background(), id)
	if err != nil {
		return respondError(c, err)
	}
	return respondSuccess(c, fiber.StatusOK, entity)
}

func (h *mediaHandler) updateFilter(c *fiber.Ctx) error {
	id := c.Params("id")
	var body struct {
		Name   string         `json:"name"`
		Effect map[string]any `json:"effect"`
		Active bool           `json:"active"`
	}
	if err := c.BodyParser(&body); err != nil {
		return respondError(c, err)
	}
	entity, err := h.filterService.Update(context.Background(), appMedia.UpdateFilterInput{
		ID:     id,
		Name:   body.Name,
		Effect: body.Effect,
		Active: body.Active,
	})
	if err != nil {
		return respondError(c, err)
	}
	return respondSuccess(c, fiber.StatusOK, entity)
}

func (h *mediaHandler) deleteFilter(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.filterService.Delete(context.Background(), id); err != nil {
		return respondError(c, err)
	}
	return respondSuccess(c, fiber.StatusNoContent, nil)
}

func (h *mediaHandler) createQRCode(c *fiber.Ctx) error {
	var body struct {
		PhotoID  string `json:"photo_id"`
		Hash     string `json:"hash"`
		ExpireAt *int64 `json:"expire_at"`
	}
	if err := c.BodyParser(&body); err != nil {
		return respondError(c, err)
	}
	if body.PhotoID == "" || body.Hash == "" {
		return respondError(c, fiber.NewError(fiber.StatusBadRequest, "photo_id and hash required"))
	}
	var expire *time.Time
	if body.ExpireAt != nil {
		ts := time.Unix(*body.ExpireAt, 0)
		expire = &ts
	}
	entity, err := h.qrService.Create(context.Background(), appMedia.CreateQRCodeInput{
		PhotoID:  body.PhotoID,
		Hash:     body.Hash,
		ExpireAt: expire,
	})
	if err != nil {
		return respondError(c, err)
	}
	return respondSuccess(c, fiber.StatusCreated, entity)
}

func (h *mediaHandler) getQRCode(c *fiber.Ctx) error {
	hash := c.Params("hash")
	entity, err := h.qrService.GetByHash(context.Background(), hash)
	if err != nil {
		return respondError(c, err)
	}
	return respondSuccess(c, fiber.StatusOK, entity)
}

func (h *mediaHandler) deleteQRCode(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.qrService.Delete(context.Background(), id); err != nil {
		return respondError(c, err)
	}
	return respondSuccess(c, fiber.StatusNoContent, nil)
}

func (h *mediaHandler) ensureSessionBelongs(ctx context.Context, sessionID string, boothID string) error {
	session, err := h.sessionService.Get(ctx, sessionID)
	if err != nil {
		return err
	}
	if session.BoothID != boothID {
		return fiber.ErrForbidden
	}
	return nil
}

func (h *mediaHandler) ensurePhotoBelongs(ctx context.Context, photoID string, boothID string) (*domainMedia.Photo, error) {
	photo, err := h.photoService.Get(ctx, photoID)
	if err != nil {
		return nil, err
	}
	if err := h.ensureSessionBelongs(ctx, photo.SessionID, boothID); err != nil {
		return nil, err
	}
	return photo, nil
}
