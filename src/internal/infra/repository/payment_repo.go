package repository

import (
	"encoding/json"
	"fmt"
	"time"

	common "github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/domain"

	"gorm.io/gorm"
)

type PaymentRepo struct {
	database *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) *PaymentRepo {
	return &PaymentRepo{
		database: db,
	}
}

func (r *PaymentRepo) GetAll(paginationRequestDto common.PaginationRequestDto) ([]domain.Payment, int64, error) {
	var payments []domain.Payment
	var count int64

	query := r.database.Model(&domain.Payment{})
	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&payments)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return payments, count, nil
}

func (r *PaymentRepo) GetAllBySiteID(siteID int64, paginationRequestDto common.PaginationRequestDto) ([]domain.Payment, int64, error) {
	var payments []domain.Payment
	var count int64

	query := r.database.Model(&domain.Payment{}).Where("site_id = ?", siteID)
	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&payments)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return payments, count, nil
}

func (r *PaymentRepo) GetAllByOrderID(orderID int64, paginationRequestDto common.PaginationRequestDto) ([]domain.Payment, int64, error) {
	var payments []domain.Payment
	var count int64

	query := r.database.Model(&domain.Payment{}).Where("order_id = ?", orderID)
	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&payments)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return payments, count, nil
}

func (r *PaymentRepo) GetAllByCustomerID(customerID int64, paginationRequestDto common.PaginationRequestDto) ([]domain.Payment, int64, error) {
	var payments []domain.Payment
	var count int64

	query := r.database.Model(&domain.Payment{}).Where("customer_id = ?", customerID)
	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&payments)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return payments, count, nil
}

func (r *PaymentRepo) GetByID(id int64) (domain.Payment, error) {
	var payment domain.Payment
	result := r.database.First(&payment, id)
	if result.Error != nil {
		return payment, result.Error
	}
	return payment, nil
}

func (r *PaymentRepo) GetByTrackingNumber(trackingNumber string) (domain.Payment, error) {
	var payment domain.Payment
	result := r.database.Where("tracking_number = ?", trackingNumber).First(&payment)
	if result.Error != nil {
		return payment, result.Error
	}
	return payment, nil
}

func (r *PaymentRepo) Create(payment domain.Payment) error {
	result := r.database.Create(&payment)
	return result.Error
}

func (r *PaymentRepo) Update(payment domain.Payment) error {
	result := r.database.Save(&payment)
	return result.Error
}

func (r *PaymentRepo) Delete(id int64) error {
	result := r.database.Delete(&domain.Payment{}, id)
	return result.Error
}

func (r *PaymentRepo) RequestPayment(amount int64, orderID int64, userID int64, gateway string, orderData map[string]string) (string, error) {
	// Serialize orderData to JSON
	orderDataJSON, err := json.Marshal(orderData)
	if err != nil {
		return "", err
	}

	// Generate a unique tracking number
	trackingNumber := time.Now().UnixNano()

	// Create a new payment record
	payment := domain.Payment{
		PaymentStatusEnum:  "Pending",
		UserType:           "User", // Assuming UserType enum value
		TrackingNumber:     trackingNumber,
		Gateway:            gateway,
		GatewayAccountName: "Default", // This would be configured per gateway in a real implementation
		Amount:             amount,
		ServiceName:        "User",                // Assuming ServiceName enum value
		ServiceAction:      "ChargeCreditRequest", // The action being performed
		OrderID:            orderID,
		ReturnUrl:          orderData["finalFrontReturnUrl"], // From orderData
		CallVerifyUrl:      "/api/v1/payment/verify",         // Standard API endpoint for verification
		ClientIp:           "127.0.0.1",                      // In a real implementation, get from request
		OrderData:          string(orderDataJSON),
		UserID:             userID,
		CustomerID:         0, // Not applicable for a user payment
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
		IsDeleted:          false,
		SiteID:             1, // Default site ID, would be properly set in a real implementation
	}

	// Save the payment
	err = r.Create(payment)
	if err != nil {
		return "", err
	}

	// In a real implementation, this would call the actual payment gateway API
	// and return a URL or token for redirecting the user to complete payment

	// For demo purposes, return a dummy payment URL
	paymentUrl := fmt.Sprintf("https://example.com/payment/gateway?tracking=%d", trackingNumber)

	return paymentUrl, nil
}
