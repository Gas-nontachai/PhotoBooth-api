package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/gofiber/swagger"

	docs "go-ddd-clean/docs"
	appAnalytics "go-ddd-clean/internal/application/analytics"
	appBooth "go-ddd-clean/internal/application/booth"
	appBranch "go-ddd-clean/internal/application/branch"
	appLogging "go-ddd-clean/internal/application/logging"
	appMedia "go-ddd-clean/internal/application/media"
	appPayment "go-ddd-clean/internal/application/payment"
	appSession "go-ddd-clean/internal/application/session"
	appUser "go-ddd-clean/internal/application/user"
	appVoucher "go-ddd-clean/internal/application/voucher"
	"go-ddd-clean/internal/infrastructure/config"
	infraDB "go-ddd-clean/internal/infrastructure/db"
	httpTransport "go-ddd-clean/internal/interface/http"
)

// @title Photobooth Platforms API
// @version 1.0
// @description API documentation for the go-ddd-clean service.
// @BasePath /api
// @schemes http
// @securityDefinitions.apikey BoothTokenAuth
// @in header
// @name Authorization
func main() {
	cfg := config.LoadConfig()
	database := infraDB.ConnectDB(cfg.DB_DSN)

	docs.SwaggerInfo.Title = "Photobooth Platforms API"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Description = "API documentation for the go-ddd-clean service."
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Schemes = []string{"http"}
	docs.SwaggerInfo.Host = fmt.Sprintf("localhost:%s", cfg.AppPort)

	branchRepo := infraDB.NewBranchRepository(database)
	boothRepo := infraDB.NewBoothRepository(database)
	sessionRepo := infraDB.NewSessionRepository(database)
	photoRepo := infraDB.NewPhotoRepository(database)
	frameRepo := infraDB.NewFrameRepository(database)
	filterRepo := infraDB.NewFilterRepository(database)
	qrRepo := infraDB.NewQRCodeRepository(database)
	userRepo := infraDB.NewUserRepository(database)
	paymentRepo := infraDB.NewPaymentRepository(database)
	voucherRepo := infraDB.NewVoucherRepository(database)
	voucherRedemptionRepo := infraDB.NewVoucherRedemptionRepository(database)
	logRepository := infraDB.NewLogRepository(database)
	analyticsRepo := infraDB.NewAnalyticsRepository(database)

	branchService := appBranch.NewService(branchRepo)
	boothService := appBooth.NewService(boothRepo)
	boothTokenService := appBooth.NewTokenService(boothRepo, cfg.BoothTokenSecret)
	sessionService := appSession.NewService(sessionRepo)
	photoService := appMedia.NewPhotoService(photoRepo)
	frameService := appMedia.NewFrameService(frameRepo)
	filterService := appMedia.NewFilterService(filterRepo)
	qrService := appMedia.NewQRCodeService(qrRepo)
	userService := appUser.NewService(userRepo)
	paymentService := appPayment.NewService(paymentRepo)
	voucherService := appVoucher.NewService(voucherRepo, voucherRedemptionRepo)
	logService := appLogging.NewService(logRepository)
	analyticsService := appAnalytics.NewService(analyticsRepo)

	router := httpTransport.NewRouter(
		branchService,
		boothService,
		boothTokenService,
		sessionService,
		photoService,
		frameService,
		filterService,
		qrService,
		userService,
		paymentService,
		voucherService,
		logService,
		analyticsService,
	)

	app := fiber.New()
	app.Get("/swagger/*", fiberSwagger.HandlerDefault)
	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})
	api := app.Group("/api")
	router.RegisterRoutes(api)

	if err := app.Listen(":" + cfg.AppPort); err != nil {
		log.Fatal(err)
	}
}
