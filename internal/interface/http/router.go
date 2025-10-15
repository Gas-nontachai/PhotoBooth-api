package http

import (
	"github.com/gofiber/fiber/v2"

	appAnalytics "go-ddd-clean/internal/application/analytics"
	appBooth "go-ddd-clean/internal/application/booth"
	appBranch "go-ddd-clean/internal/application/branch"
	appLogging "go-ddd-clean/internal/application/logging"
	appMedia "go-ddd-clean/internal/application/media"
	appPayment "go-ddd-clean/internal/application/payment"
	appSession "go-ddd-clean/internal/application/session"
	appUser "go-ddd-clean/internal/application/user"
	appVoucher "go-ddd-clean/internal/application/voucher"
)

type Router struct {
	branch      *appBranch.Service
	booth       *appBooth.Service
	boothTokens *appBooth.TokenService
	session     *appSession.Service
	photos      *appMedia.PhotoService
	frames      *appMedia.FrameService
	filters     *appMedia.FilterService
	qrcodes     *appMedia.QRCodeService
	user        *appUser.Service
	payment     *appPayment.Service
	voucher     *appVoucher.Service
	logging     *appLogging.Service
	analytics   *appAnalytics.Service
}

func NewRouter(
	branch *appBranch.Service,
	booth *appBooth.Service,
	boothTokens *appBooth.TokenService,
	session *appSession.Service,
	photos *appMedia.PhotoService,
	frames *appMedia.FrameService,
	filters *appMedia.FilterService,
	qrcodes *appMedia.QRCodeService,
	user *appUser.Service,
	payment *appPayment.Service,
	voucher *appVoucher.Service,
	logging *appLogging.Service,
	analytics *appAnalytics.Service,
) *Router {
	return &Router{
		branch:      branch,
		booth:       booth,
		boothTokens: boothTokens,
		session:     session,
		photos:      photos,
		frames:      frames,
		filters:     filters,
		qrcodes:     qrcodes,
		user:        user,
		payment:     payment,
		voucher:     voucher,
		logging:     logging,
		analytics:   analytics,
	}
}

func (r *Router) RegisterRoutes(router fiber.Router) {
	branchHandler := newBranchHandler(r.branch)
	boothHandler := newBoothHandler(r.booth, r.logging, r.analytics)
	sessionHandler := newSessionHandler(r.session, r.photos, r.payment)
	mediaHandler := newMediaHandler(r.session, r.photos, r.frames, r.filters, r.qrcodes)
	userHandler := newUserHandler(r.user)
	paymentHandler := newPaymentHandler(r.payment)
	voucherHandler := newVoucherHandler(r.voucher, r.session)
	boothTokenHandler := newBoothTokenHandler(r.boothTokens)
	boothAuth := newBoothAuthMiddleware(r.boothTokens)

	router.Post("/booth/register", boothTokenHandler.register)
	router.Post("/booth/regenerate-token", boothAuth, boothTokenHandler.regenerate)
	branchHandler.register(router.Group("/branches"))
	boothHandler.register(router.Group("/booths"))
	sessionHandler.register(router.Group("/sessions"), boothAuth)
	mediaHandler.register(router.Group("/media"), boothAuth)
	userHandler.register(router.Group("/users"))
	paymentHandler.register(router.Group("/payments"))
	voucherHandler.register(router.Group("/vouchers"), boothAuth)
}
