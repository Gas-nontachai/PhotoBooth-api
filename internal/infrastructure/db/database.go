package db

import (
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectDB(dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		log.Fatal("❌ Failed to connect database:", err)
	}
	if err := db.AutoMigrate(
		&BranchModel{},
		&UserModel{},
		&PaymentModel{},
		&VoucherModel{},
		&FrameModel{},
		&FilterModel{},
	); err != nil {
		log.Fatal("❌ Failed to run base migrations:", err)
	}
	if err := db.AutoMigrate(
		&BoothModel{},
		&SessionModel{},
		&PhotoModel{},
		&QRCodeModel{},
		&VoucherRedemptionModel{},
		&BoothLogModel{},
		&AnalyticsEventModel{},
	); err != nil {
		log.Fatal("❌ Failed to run migrations:", err)
	}
	log.Println("✅ Connected to database and ran migrations")
	return db
}
