package main

import (
	domainAnalytics "go-ddd-clean/internal/domain/analytics"
	domainBooth "go-ddd-clean/internal/domain/booth"
	domainBranch "go-ddd-clean/internal/domain/branch"
	domainLogging "go-ddd-clean/internal/domain/logging"
	domainMedia "go-ddd-clean/internal/domain/media"
	domainPayment "go-ddd-clean/internal/domain/payment"
	domainSession "go-ddd-clean/internal/domain/session"
	domainUser "go-ddd-clean/internal/domain/user"
	domainVoucher "go-ddd-clean/internal/domain/voucher"
)

type Branch = domainBranch.Branch
type Booth = domainBooth.Booth
type BoothLog = domainLogging.BoothLog
type AnalyticsEvent = domainAnalytics.Event
type Session = domainSession.Session
type Photo = domainMedia.Photo
type Frame = domainMedia.Frame
type Filter = domainMedia.Filter
type QRCode = domainMedia.QRCode
type User = domainUser.User
type Payment = domainPayment.Payment
type Voucher = domainVoucher.Voucher
type VoucherRedemption = domainVoucher.Redemption

type ErrorResponse struct {
	Error string `json:"error"`
}

type HealthResponse struct {
	Status string `json:"status"`
}

type BoothTokenRegisterRequest struct {
	BoothID  string `json:"booth_id"`
	BranchID string `json:"branch_id"`
}

type BoothTokenResponse struct {
	Token string `json:"token"`
}

type BranchCreateRequest struct {
	Name     string  `json:"name"`
	Location *string `json:"location"`
}

type BranchUpdateRequest struct {
	Name     string  `json:"name"`
	Location *string `json:"location"`
}

type BoothCreateRequest struct {
	BranchID string         `json:"branch_id"`
	Name     string         `json:"name"`
	Type     string         `json:"type"`
	Status   string         `json:"status"`
	Config   map[string]any `json:"config"`
}

type BoothUpdateRequest struct {
	BranchID string         `json:"branch_id"`
	Name     string         `json:"name"`
	Type     *string        `json:"type"`
	Status   *string        `json:"status"`
	Config   map[string]any `json:"config"`
}

type BoothLogCreateRequest struct {
	EventType string  `json:"event_type"`
	Level     string  `json:"level"`
	Message   *string `json:"message"`
}

type BoothAnalyticsEventRequest struct {
	SessionID *string        `json:"session_id"`
	EventName string         `json:"event_name"`
	Payload   map[string]any `json:"payload"`
}

type SessionCreateRequest struct {
	BoothID       string         `json:"booth_id"`
	UserID        *string        `json:"user_id"`
	VoucherID     *string        `json:"voucher_id"`
	PaymentID     *string        `json:"payment_id"`
	Status        *string        `json:"status"`
	TotalPrice    *float64       `json:"total_price"`
	BoothSnapshot map[string]any `json:"booth_snapshot"`
	PhoneTemp     *string        `json:"phone_temp"`
}

type SessionUpdateRequest struct {
	UserID        *string        `json:"user_id"`
	VoucherID     *string        `json:"voucher_id"`
	PaymentID     *string        `json:"payment_id"`
	Status        *string        `json:"status"`
	TotalPrice    *float64       `json:"total_price"`
	FinishedAt    *int64         `json:"finished_at"`
	BoothSnapshot map[string]any `json:"booth_snapshot"`
	PhoneTemp     *string        `json:"phone_temp"`
}

type PhotoCreateRequest struct {
	SessionID   string         `json:"session_id"`
	FrameID     *string        `json:"frame_id"`
	FilterID    *string        `json:"filter_id"`
	StorageURL  string         `json:"storage_url"`
	Composition map[string]any `json:"composition"`
	RenderedURL *string        `json:"rendered_url"`
}

type PhotoUpdateRequest struct {
	FrameID     *string        `json:"frame_id"`
	FilterID    *string        `json:"filter_id"`
	StorageURL  *string        `json:"storage_url"`
	Composition map[string]any `json:"composition"`
	RenderedURL *string        `json:"rendered_url"`
}

type FrameCreateRequest struct {
	Name    string  `json:"name"`
	Theme   *string `json:"theme"`
	FileURL string  `json:"file_url"`
	Active  *bool   `json:"active"`
}

type FrameUpdateRequest struct {
	Name    string  `json:"name"`
	Theme   *string `json:"theme"`
	FileURL string  `json:"file_url"`
	Active  bool    `json:"active"`
}

type FilterCreateRequest struct {
	Name   string         `json:"name"`
	Effect map[string]any `json:"effect"`
	Active *bool          `json:"active"`
}

type FilterUpdateRequest struct {
	Name   string         `json:"name"`
	Effect map[string]any `json:"effect"`
	Active bool           `json:"active"`
}

type QRCodeCreateRequest struct {
	PhotoID  string `json:"photo_id"`
	Hash     string `json:"hash"`
	ExpireAt *int64 `json:"expire_at"`
}

type UserCreateRequest struct {
	Tel      *string `json:"tel"`
	Email    *string `json:"email"`
	Password *string `json:"password"`
	Role     *string `json:"role"`
}

type UserUpdateRequest struct {
	Tel      *string `json:"tel"`
	Email    *string `json:"email"`
	Password *string `json:"password"`
	Role     *string `json:"role"`
	Points   *int    `json:"points"`
}

type UserAdjustPointsRequest struct {
	Delta int `json:"delta"`
}

type PaymentCreateRequest struct {
	SessionID      string  `json:"session_id"`
	Method         string  `json:"method"`
	Amount         float64 `json:"amount"`
	Currency       string  `json:"currency"`
	Status         string  `json:"status"`
	TransactionRef *string `json:"transaction_ref"`
}

type PaymentUpdateRequest struct {
	Status         *string  `json:"status"`
	TransactionRef *string  `json:"transaction_ref"`
	Amount         *float64 `json:"amount"`
	Currency       *string  `json:"currency"`
	Method         *string  `json:"method"`
}

type VoucherCreateRequest struct {
	Code      string  `json:"code"`
	Type      string  `json:"type"`
	Value     float64 `json:"value"`
	Unit      string  `json:"unit"`
	MaxUsage  int     `json:"max_usage"`
	ValidFrom *int64  `json:"valid_from"`
	ValidTo   *int64  `json:"valid_to"`
	Active    *bool   `json:"active"`
}

type VoucherUpdateRequest struct {
	Value     *float64 `json:"value"`
	Unit      *string  `json:"unit"`
	MaxUsage  *int     `json:"max_usage"`
	ValidFrom *int64   `json:"valid_from"`
	ValidTo   *int64   `json:"valid_to"`
	Active    *bool    `json:"active"`
}

type VoucherRedeemRequest struct {
	Code      string   `json:"code"`
	SessionID string   `json:"session_id"`
	Tel       *string  `json:"tel"`
	Discount  *float64 `json:"discount"`
}

type VoucherRedeemResponse struct {
	Redemption VoucherRedemption `json:"redemption"`
	Voucher    Voucher           `json:"voucher"`
}

// healthzDoc godoc
// @Summary ตรวจสอบสถานะเซิร์ฟเวอร์
// @Tags Health
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /healthz [get]
func healthzDoc() {}

// boothRegisterDoc godoc
// @Summary ลงทะเบียนบูธ
// @Description ขอรับโทเคนสำหรับบูธใหม่
// @Tags Booth Token
// @Accept json
// @Produce json
// @Param payload body BoothTokenRegisterRequest true "ข้อมูลบูธ"
// @Success 200 {object} BoothTokenResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/booth/register [post]
func boothRegisterDoc() {}

// boothRegenerateDoc godoc
// @Summary สร้างโทเคนบูธใหม่
// @Description สร้างโทเคนใหม่จากโทเคนเดิม
// @Tags Booth Token
// @Produce json
// @Security BoothTokenAuth
// @Success 200 {object} BoothTokenResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/booth/regenerate-token [post]
func boothRegenerateDoc() {}

// branchListDoc godoc
// @Summary ดึงรายการสาขา
// @Tags Branches
// @Produce json
// @Success 200 {array} Branch
// @Failure 500 {object} ErrorResponse
// @Router /api/branches [get]
func branchListDoc() {}

// branchCreateDoc godoc
// @Summary สร้างสาขาใหม่
// @Tags Branches
// @Accept json
// @Produce json
// @Param payload body BranchCreateRequest true "ข้อมูลสาขา"
// @Success 201 {object} Branch
// @Failure 400 {object} ErrorResponse
// @Router /api/branches [post]
func branchCreateDoc() {}

// branchGetDoc godoc
// @Summary ดูข้อมูลสาขา
// @Tags Branches
// @Produce json
// @Param id path string true "รหัสสาขา"
// @Success 200 {object} Branch
// @Failure 404 {object} ErrorResponse
// @Router /api/branches/{id} [get]
func branchGetDoc() {}

// branchUpdateDoc godoc
// @Summary ปรับปรุงข้อมูลสาขา
// @Tags Branches
// @Accept json
// @Produce json
// @Param id path string true "รหัสสาขา"
// @Param payload body BranchUpdateRequest true "ข้อมูลที่ต้องการแก้ไข"
// @Success 200 {object} Branch
// @Failure 400 {object} ErrorResponse
// @Router /api/branches/{id} [put]
func branchUpdateDoc() {}

// branchDeleteDoc godoc
// @Summary ลบสาขา
// @Tags Branches
// @Param id path string true "รหัสสาขา"
// @Success 204 {string} string "No Content"
// @Failure 404 {object} ErrorResponse
// @Router /api/branches/{id} [delete]
func branchDeleteDoc() {}

// boothListDoc godoc
// @Summary ดึงรายการบูธ
// @Tags Booths
// @Produce json
// @Param branch_id query string false "กรองตามสาขา"
// @Success 200 {array} Booth
// @Failure 500 {object} ErrorResponse
// @Router /api/booths [get]
func boothListDoc() {}

// boothCreateDoc godoc
// @Summary สร้างบูธใหม่
// @Tags Booths
// @Accept json
// @Produce json
// @Param payload body BoothCreateRequest true "ข้อมูลบูธ"
// @Success 201 {object} Booth
// @Failure 400 {object} ErrorResponse
// @Router /api/booths [post]
func boothCreateDoc() {}

// boothGetDoc godoc
// @Summary ดูข้อมูลบูธ
// @Tags Booths
// @Produce json
// @Param id path string true "รหัสบูธ"
// @Success 200 {object} Booth
// @Failure 404 {object} ErrorResponse
// @Router /api/booths/{id} [get]
func boothGetDoc() {}

// boothUpdateDoc godoc
// @Summary ปรับปรุงข้อมูลบูธ
// @Tags Booths
// @Accept json
// @Produce json
// @Param id path string true "รหัสบูธ"
// @Param payload body BoothUpdateRequest true "ข้อมูลที่ต้องการแก้ไข"
// @Success 200 {object} Booth
// @Failure 400 {object} ErrorResponse
// @Router /api/booths/{id} [put]
func boothUpdateDoc() {}

// boothDeleteDoc godoc
// @Summary ลบบูธ
// @Tags Booths
// @Param id path string true "รหัสบูธ"
// @Success 204 {string} string "No Content"
// @Failure 404 {object} ErrorResponse
// @Router /api/booths/{id} [delete]
func boothDeleteDoc() {}

// boothLogsListDoc godoc
// @Summary ดึงบันทึกของบูธ
// @Tags Booth Logs
// @Produce json
// @Param id path string true "รหัสบูธ"
// @Param limit query int false "จำนวนรายการ (ค่าเริ่มต้น 50)"
// @Success 200 {array} BoothLog
// @Failure 404 {object} ErrorResponse
// @Router /api/booths/{id}/logs [get]
func boothLogsListDoc() {}

// boothLogsCreateDoc godoc
// @Summary บันทึกเหตุการณ์ของบูธ
// @Tags Booth Logs
// @Accept json
// @Produce json
// @Param id path string true "รหัสบูธ"
// @Param payload body BoothLogCreateRequest true "ข้อมูลบันทึก"
// @Success 201 {object} BoothLog
// @Failure 400 {object} ErrorResponse
// @Router /api/booths/{id}/logs [post]
func boothLogsCreateDoc() {}

// boothAnalyticsListDoc godoc
// @Summary ดึงข้อมูลการใช้งานบูธ
// @Tags Booth Analytics
// @Produce json
// @Param id path string true "รหัสบูธ"
// @Param limit query int false "จำนวนรายการ (ค่าเริ่มต้น 100)"
// @Success 200 {array} AnalyticsEvent
// @Failure 404 {object} ErrorResponse
// @Router /api/booths/{id}/analytics [get]
func boothAnalyticsListDoc() {}

// boothAnalyticsCreateDoc godoc
// @Summary สร้างข้อมูลการใช้งานบูธ
// @Tags Booth Analytics
// @Accept json
// @Produce json
// @Param id path string true "รหัสบูธ"
// @Param payload body BoothAnalyticsEventRequest true "ข้อมูลเหตุการณ์"
// @Success 201 {object} AnalyticsEvent
// @Failure 400 {object} ErrorResponse
// @Router /api/booths/{id}/analytics [post]
func boothAnalyticsCreateDoc() {}

// sessionListDoc godoc
// @Summary ดึงรายการเซสชัน
// @Tags Sessions
// @Produce json
// @Security BoothTokenAuth
// @Param status query string false "กรองตามสถานะ"
// @Success 200 {array} Session
// @Failure 401 {object} ErrorResponse
// @Router /api/sessions [get]
func sessionListDoc() {}

// sessionCreateDoc godoc
// @Summary สร้างเซสชัน
// @Tags Sessions
// @Accept json
// @Produce json
// @Security BoothTokenAuth
// @Param payload body SessionCreateRequest true "ข้อมูลเซสชัน"
// @Success 201 {object} Session
// @Failure 400 {object} ErrorResponse
// @Router /api/sessions [post]
func sessionCreateDoc() {}

// sessionGetDoc godoc
// @Summary ดูข้อมูลเซสชัน
// @Tags Sessions
// @Produce json
// @Security BoothTokenAuth
// @Param id path string true "รหัสเซสชัน"
// @Success 200 {object} Session
// @Failure 404 {object} ErrorResponse
// @Router /api/sessions/{id} [get]
func sessionGetDoc() {}

// sessionUpdateDoc godoc
// @Summary ปรับปรุงข้อมูลเซสชัน
// @Tags Sessions
// @Accept json
// @Produce json
// @Security BoothTokenAuth
// @Param id path string true "รหัสเซสชัน"
// @Param payload body SessionUpdateRequest true "ข้อมูลที่ต้องการแก้ไข"
// @Success 200 {object} Session
// @Failure 400 {object} ErrorResponse
// @Router /api/sessions/{id} [put]
func sessionUpdateDoc() {}

// sessionDeleteDoc godoc
// @Summary ลบเซสชัน
// @Tags Sessions
// @Security BoothTokenAuth
// @Param id path string true "รหัสเซสชัน"
// @Success 204 {string} string "No Content"
// @Failure 404 {object} ErrorResponse
// @Router /api/sessions/{id} [delete]
func sessionDeleteDoc() {}

// sessionPhotosListDoc godoc
// @Summary ดึงรายการรูปของเซสชัน
// @Tags Sessions
// @Produce json
// @Security BoothTokenAuth
// @Param id path string true "รหัสเซสชัน"
// @Success 200 {array} Photo
// @Failure 404 {object} ErrorResponse
// @Router /api/sessions/{id}/photos [get]
func sessionPhotosListDoc() {}

// sessionPaymentGetDoc godoc
// @Summary ดูข้อมูลการชำระเงินของเซสชัน
// @Tags Sessions
// @Produce json
// @Security BoothTokenAuth
// @Param id path string true "รหัสเซสชัน"
// @Success 200 {object} Payment
// @Failure 404 {object} ErrorResponse
// @Router /api/sessions/{id}/payment [get]
func sessionPaymentGetDoc() {}

// mediaPhotosListDoc godoc
// @Summary ดึงรายการรูป
// @Tags Media Photos
// @Produce json
// @Security BoothTokenAuth
// @Param session_id query string true "รหัสเซสชัน"
// @Success 200 {array} Photo
// @Failure 400 {object} ErrorResponse
// @Router /api/media/photos [get]
func mediaPhotosListDoc() {}

// mediaPhotosCreateDoc godoc
// @Summary สร้างรูปใหม่
// @Tags Media Photos
// @Accept json
// @Produce json
// @Security BoothTokenAuth
// @Param payload body PhotoCreateRequest true "ข้อมูลรูป"
// @Success 201 {object} Photo
// @Failure 400 {object} ErrorResponse
// @Router /api/media/photos [post]
func mediaPhotosCreateDoc() {}

// mediaPhotosGetDoc godoc
// @Summary ดูข้อมูลรูป
// @Tags Media Photos
// @Produce json
// @Security BoothTokenAuth
// @Param id path string true "รหัสรูป"
// @Success 200 {object} Photo
// @Failure 404 {object} ErrorResponse
// @Router /api/media/photos/{id} [get]
func mediaPhotosGetDoc() {}

// mediaPhotosUpdateDoc godoc
// @Summary ปรับปรุงข้อมูลรูป
// @Tags Media Photos
// @Accept json
// @Produce json
// @Security BoothTokenAuth
// @Param id path string true "รหัสรูป"
// @Param payload body PhotoUpdateRequest true "ข้อมูลที่ต้องการแก้ไข"
// @Success 200 {object} Photo
// @Failure 400 {object} ErrorResponse
// @Router /api/media/photos/{id} [put]
func mediaPhotosUpdateDoc() {}

// mediaPhotosDeleteDoc godoc
// @Summary ลบรูป
// @Tags Media Photos
// @Security BoothTokenAuth
// @Param id path string true "รหัสรูป"
// @Success 204 {string} string "No Content"
// @Failure 404 {object} ErrorResponse
// @Router /api/media/photos/{id} [delete]
func mediaPhotosDeleteDoc() {}

// mediaFramesListDoc godoc
// @Summary ดึงรายการกรอบรูป
// @Tags Media Frames
// @Produce json
// @Param active query bool false "ดึงเฉพาะที่เปิดใช้งาน"
// @Success 200 {array} Frame
// @Failure 500 {object} ErrorResponse
// @Router /api/media/frames [get]
func mediaFramesListDoc() {}

// mediaFramesCreateDoc godoc
// @Summary สร้างกรอบรูป
// @Tags Media Frames
// @Accept json
// @Produce json
// @Param payload body FrameCreateRequest true "ข้อมูลกรอบรูป"
// @Success 201 {object} Frame
// @Failure 400 {object} ErrorResponse
// @Router /api/media/frames [post]
func mediaFramesCreateDoc() {}

// mediaFramesGetDoc godoc
// @Summary ดูข้อมูลกรอบรูป
// @Tags Media Frames
// @Produce json
// @Param id path string true "รหัสกรอบรูป"
// @Success 200 {object} Frame
// @Failure 404 {object} ErrorResponse
// @Router /api/media/frames/{id} [get]
func mediaFramesGetDoc() {}

// mediaFramesUpdateDoc godoc
// @Summary ปรับปรุงกรอบรูป
// @Tags Media Frames
// @Accept json
// @Produce json
// @Param id path string true "รหัสกรอบรูป"
// @Param payload body FrameUpdateRequest true "ข้อมูลที่ต้องแก้ไข"
// @Success 200 {object} Frame
// @Failure 400 {object} ErrorResponse
// @Router /api/media/frames/{id} [put]
func mediaFramesUpdateDoc() {}

// mediaFramesDeleteDoc godoc
// @Summary ลบกรอบรูป
// @Tags Media Frames
// @Param id path string true "รหัสกรอบรูป"
// @Success 204 {string} string "No Content"
// @Failure 404 {object} ErrorResponse
// @Router /api/media/frames/{id} [delete]
func mediaFramesDeleteDoc() {}

// mediaFiltersListDoc godoc
// @Summary ดึงรายการฟิลเตอร์
// @Tags Media Filters
// @Produce json
// @Param active query bool false "ดึงเฉพาะที่เปิดใช้งาน"
// @Success 200 {array} Filter
// @Failure 500 {object} ErrorResponse
// @Router /api/media/filters [get]
func mediaFiltersListDoc() {}

// mediaFiltersCreateDoc godoc
// @Summary สร้างฟิลเตอร์
// @Tags Media Filters
// @Accept json
// @Produce json
// @Param payload body FilterCreateRequest true "ข้อมูลฟิลเตอร์"
// @Success 201 {object} Filter
// @Failure 400 {object} ErrorResponse
// @Router /api/media/filters [post]
func mediaFiltersCreateDoc() {}

// mediaFiltersGetDoc godoc
// @Summary ดูข้อมูลฟิลเตอร์
// @Tags Media Filters
// @Produce json
// @Param id path string true "รหัสฟิลเตอร์"
// @Success 200 {object} Filter
// @Failure 404 {object} ErrorResponse
// @Router /api/media/filters/{id} [get]
func mediaFiltersGetDoc() {}

// mediaFiltersUpdateDoc godoc
// @Summary ปรับปรุงฟิลเตอร์
// @Tags Media Filters
// @Accept json
// @Produce json
// @Param id path string true "รหัสฟิลเตอร์"
// @Param payload body FilterUpdateRequest true "ข้อมูลที่ต้องแก้ไข"
// @Success 200 {object} Filter
// @Failure 400 {object} ErrorResponse
// @Router /api/media/filters/{id} [put]
func mediaFiltersUpdateDoc() {}

// mediaFiltersDeleteDoc godoc
// @Summary ลบฟิลเตอร์
// @Tags Media Filters
// @Param id path string true "รหัสฟิลเตอร์"
// @Success 204 {string} string "No Content"
// @Failure 404 {object} ErrorResponse
// @Router /api/media/filters/{id} [delete]
func mediaFiltersDeleteDoc() {}

// mediaQRCodesCreateDoc godoc
// @Summary สร้าง QR Code
// @Tags Media QRCodes
// @Accept json
// @Produce json
// @Param payload body QRCodeCreateRequest true "ข้อมูล QR Code"
// @Success 201 {object} QRCode
// @Failure 400 {object} ErrorResponse
// @Router /api/media/qrcodes [post]
func mediaQRCodesCreateDoc() {}

// mediaQRCodesGetDoc godoc
// @Summary ดูข้อมูล QR Code
// @Tags Media QRCodes
// @Produce json
// @Param hash path string true "แฮชของ QR Code"
// @Success 200 {object} QRCode
// @Failure 404 {object} ErrorResponse
// @Router /api/media/qrcodes/{hash} [get]
func mediaQRCodesGetDoc() {}

// mediaQRCodesDeleteDoc godoc
// @Summary ลบ QR Code
// @Tags Media QRCodes
// @Param id path string true "รหัส QR Code"
// @Success 204 {string} string "No Content"
// @Failure 404 {object} ErrorResponse
// @Router /api/media/qrcodes/{id} [delete]
func mediaQRCodesDeleteDoc() {}

// userListDoc godoc
// @Summary ดึงรายการผู้ใช้
// @Tags Users
// @Produce json
// @Success 200 {array} User
// @Failure 500 {object} ErrorResponse
// @Router /api/users [get]
func userListDoc() {}

// userCreateDoc godoc
// @Summary สร้างผู้ใช้ใหม่
// @Tags Users
// @Accept json
// @Produce json
// @Param payload body UserCreateRequest true "ข้อมูลผู้ใช้"
// @Success 201 {object} User
// @Failure 400 {object} ErrorResponse
// @Router /api/users [post]
func userCreateDoc() {}

// userGetDoc godoc
// @Summary ดูข้อมูลผู้ใช้
// @Tags Users
// @Produce json
// @Param id path string true "รหัสผู้ใช้"
// @Success 200 {object} User
// @Failure 404 {object} ErrorResponse
// @Router /api/users/{id} [get]
func userGetDoc() {}

// userUpdateDoc godoc
// @Summary ปรับปรุงข้อมูลผู้ใช้
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "รหัสผู้ใช้"
// @Param payload body UserUpdateRequest true "ข้อมูลที่ต้องแก้ไข"
// @Success 200 {object} User
// @Failure 400 {object} ErrorResponse
// @Router /api/users/{id} [put]
func userUpdateDoc() {}

// userDeleteDoc godoc
// @Summary ลบผู้ใช้
// @Tags Users
// @Param id path string true "รหัสผู้ใช้"
// @Success 204 {string} string "No Content"
// @Failure 404 {object} ErrorResponse
// @Router /api/users/{id} [delete]
func userDeleteDoc() {}

// userAdjustPointsDoc godoc
// @Summary ปรับแต้มผู้ใช้
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "รหัสผู้ใช้"
// @Param payload body UserAdjustPointsRequest true "ค่าที่ต้องการปรับ"
// @Success 200 {object} User
// @Failure 400 {object} ErrorResponse
// @Router /api/users/{id}/points [post]
func userAdjustPointsDoc() {}

// paymentCreateDoc godoc
// @Summary สร้างข้อมูลการชำระเงิน
// @Tags Payments
// @Accept json
// @Produce json
// @Param payload body PaymentCreateRequest true "ข้อมูลการชำระเงิน"
// @Success 201 {object} Payment
// @Failure 400 {object} ErrorResponse
// @Router /api/payments [post]
func paymentCreateDoc() {}

// paymentGetDoc godoc
// @Summary ดูข้อมูลการชำระเงิน
// @Tags Payments
// @Produce json
// @Param id path string true "รหัสการชำระเงิน"
// @Success 200 {object} Payment
// @Failure 404 {object} ErrorResponse
// @Router /api/payments/{id} [get]
func paymentGetDoc() {}

// paymentUpdateDoc godoc
// @Summary ปรับปรุงข้อมูลการชำระเงิน
// @Tags Payments
// @Accept json
// @Produce json
// @Param id path string true "รหัสการชำระเงิน"
// @Param payload body PaymentUpdateRequest true "ข้อมูลที่ต้องแก้ไข"
// @Success 200 {object} Payment
// @Failure 400 {object} ErrorResponse
// @Router /api/payments/{id} [put]
func paymentUpdateDoc() {}

// paymentGetBySessionDoc godoc
// @Summary ดูข้อมูลการชำระเงินของเซสชัน
// @Tags Payments
// @Produce json
// @Param sessionID path string true "รหัสเซสชัน"
// @Success 200 {object} Payment
// @Failure 404 {object} ErrorResponse
// @Router /api/payments/session/{sessionID} [get]
func paymentGetBySessionDoc() {}

// voucherListDoc godoc
// @Summary ดึงรายการคูปอง
// @Tags Vouchers
// @Produce json
// @Param active query bool false "ดึงเฉพาะที่เปิดใช้งาน"
// @Success 200 {array} Voucher
// @Failure 500 {object} ErrorResponse
// @Router /api/vouchers [get]
func voucherListDoc() {}

// voucherCreateDoc godoc
// @Summary สร้างคูปองใหม่
// @Tags Vouchers
// @Accept json
// @Produce json
// @Param payload body VoucherCreateRequest true "ข้อมูลคูปอง"
// @Success 201 {object} Voucher
// @Failure 400 {object} ErrorResponse
// @Router /api/vouchers [post]
func voucherCreateDoc() {}

// voucherGetDoc godoc
// @Summary ดูข้อมูลคูปอง
// @Tags Vouchers
// @Produce json
// @Param id path string true "รหัสคูปอง"
// @Success 200 {object} Voucher
// @Failure 404 {object} ErrorResponse
// @Router /api/vouchers/{id} [get]
func voucherGetDoc() {}

// voucherUpdateDoc godoc
// @Summary ปรับปรุงคูปอง
// @Tags Vouchers
// @Accept json
// @Produce json
// @Param id path string true "รหัสคูปอง"
// @Param payload body VoucherUpdateRequest true "ข้อมูลที่ต้องแก้ไข"
// @Success 200 {object} Voucher
// @Failure 400 {object} ErrorResponse
// @Router /api/vouchers/{id} [put]
func voucherUpdateDoc() {}

// voucherDeleteDoc godoc
// @Summary ลบคูปอง
// @Tags Vouchers
// @Param id path string true "รหัสคูปอง"
// @Success 204 {string} string "No Content"
// @Failure 404 {object} ErrorResponse
// @Router /api/vouchers/{id} [delete]
func voucherDeleteDoc() {}

// voucherGetByCodeDoc godoc
// @Summary ค้นหาคูปองด้วยรหัส
// @Tags Vouchers
// @Produce json
// @Param code path string true "รหัสคูปอง"
// @Success 200 {object} Voucher
// @Failure 404 {object} ErrorResponse
// @Router /api/vouchers/code/{code} [get]
func voucherGetByCodeDoc() {}

// voucherRedeemDoc godoc
// @Summary ใช้งานคูปอง
// @Tags Vouchers
// @Accept json
// @Produce json
// @Security BoothTokenAuth
// @Param payload body VoucherRedeemRequest true "ข้อมูลการใช้งานคูปอง"
// @Success 200 {object} VoucherRedeemResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/vouchers/redeem [post]
func voucherRedeemDoc() {}
