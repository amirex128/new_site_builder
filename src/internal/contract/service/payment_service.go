package service

import (
	"github.com/amirex128/new_site_builder/src/internal/domain"
	"github.com/amirex128/new_site_builder/src/internal/domain/enums"
)

// IPaymentService defines the interface for payment gateway operations
type IPaymentService interface {
	// RequestPayment initiates a payment request to the payment gateway
	RequestPayment(amount int64, orderID int64, userID int64, gateway string, orderData map[string]string) (string, error)

	// VerifyPayment verifies a payment with the payment gateway
	VerifyPayment(transactionCode string, paymentID int64) (bool, error)

	// CancelPayment cancels a payment
	CancelPayment(paymentID int64) error

	// RefundPayment refunds a payment
	RefundPayment(paymentID int64, amount int64) error

	// GetPaymentStatus gets the current status of a payment
	GetPaymentStatus(paymentID int64) (enums.StatusEnum, error)

	// GetGatewayByID gets gateway configuration by ID
	GetGatewayByID(gatewayID int64) (*domain.Gateway, error)

	// GetGatewayBySiteID gets gateway configuration by site ID
	GetGatewayBySiteID(siteID int64) (*domain.Gateway, error)
}
