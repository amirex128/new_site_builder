package service

import (
	"encoding/json"
	"errors"
	"fmt"
	repository2 "github.com/amirex128/new_site_builder/internal/contract/repository"
	domain2 "github.com/amirex128/new_site_builder/internal/domain"
	"github.com/amirex128/new_site_builder/internal/domain/enums"
	"strconv"
	"time"
)

// PaymentService implements the payment gateway operations
type PaymentService struct {
	paymentRepo repository2.IPaymentRepository
	gatewayRepo repository2.IGatewayRepository
}

// NewPaymentService creates a new payment service
func NewPaymentService(
	paymentRepo repository2.IPaymentRepository,
	gatewayRepo repository2.IGatewayRepository,
) *PaymentService {
	return &PaymentService{
		paymentRepo: paymentRepo,
		gatewayRepo: gatewayRepo,
	}
}

// RequestPayment initiates a payment request to the payment gateway
func (s *PaymentService) RequestPayment(amount int64, orderID int64, userID int64, gateway string, orderData map[string]string) (string, error) {
	// Get gateway configuration
	gatewayID, err := strconv.ParseInt(gateway, 10, 64)
	if err != nil {
		return "", errors.New("invalid gateway ID")
	}

	// Get site ID from order data
	siteIDStr, ok := orderData["SiteId"]
	if !ok {
		return "", errors.New("site ID not found in order data")
	}

	siteID, err := strconv.ParseInt(siteIDStr, 10, 64)
	if err != nil {
		return "", errors.New("invalid site ID")
	}

	// Get gateway configuration
	gatewayConfig, err := s.gatewayRepo.GetBySiteID(siteID)
	if err != nil {
		return "", err
	}

	// Check if the gateway is active
	if !s.isGatewayActive(gatewayConfig, int(gatewayID)) {
		return "", errors.New("selected gateway is not active")
	}

	// Generate tracking number
	trackingNumber := time.Now().UnixNano()

	// Serialize order data
	orderDataJSON, err := json.Marshal(orderData)
	if err != nil {
		return "", err
	}

	// Create payment record
	payment := domain2.Payment{
		SiteID:             siteID,
		PaymentStatusEnum:  "Processing",
		TrackingNumber:     trackingNumber,
		Gateway:            s.mapGatewayToString(gatewayID),
		GatewayAccountName: fmt.Sprintf("%d-%d", gatewayID, siteID),
		Amount:             amount,
		OrderID:            orderID,
		OrderData:          string(orderDataJSON),
		UserID:             userID,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	// Save payment record
	if err := s.paymentRepo.Create(&payment); err != nil {
		return "", err
	}

	// In a real implementation, this would make an HTTP request to the payment gateway
	// For now, we'll simulate it by generating a payment URL
	paymentURL := fmt.Sprintf("https://payment-gateway.example.com/pay?tracking=%d", trackingNumber)

	return paymentURL, nil
}

// VerifyPayment verifies a payment with the payment gateway
func (s *PaymentService) VerifyPayment(transactionCode string, paymentID int64) (bool, error) {
	// Get payment record
	payment, err := s.paymentRepo.GetByID(paymentID)
	if err != nil {
		return false, err
	}

	// In a real implementation, this would make an HTTP request to the payment gateway
	// For now, we'll simulate it by checking if the transaction code matches
	if payment.TransactionCode != transactionCode {
		payment.PaymentStatusEnum = "Failed"
		payment.UpdatedAt = time.Now()
		if err := s.paymentRepo.Update(payment); err != nil {
			return false, err
		}
		return false, nil
	}

	// Update payment status
	payment.PaymentStatusEnum = "Verified"
	payment.UpdatedAt = time.Now()

	if err := s.paymentRepo.Update(payment); err != nil {
		return false, err
	}

	return true, nil
}

// CancelPayment cancels a payment
func (s *PaymentService) CancelPayment(paymentID int64) error {
	// Get payment record
	payment, err := s.paymentRepo.GetByID(paymentID)
	if err != nil {
		return err
	}

	// Check if payment can be cancelled
	if payment.PaymentStatusEnum == "Verified" || payment.PaymentStatusEnum == "Succeed" {
		return errors.New("cannot cancel a verified or successful payment")
	}

	// Update payment status
	payment.PaymentStatusEnum = "Cancelled"
	payment.UpdatedAt = time.Now()

	return s.paymentRepo.Update(payment)
}

// RefundPayment refunds a payment
func (s *PaymentService) RefundPayment(paymentID int64, amount int64) error {
	// Get payment record
	payment, err := s.paymentRepo.GetByID(paymentID)
	if err != nil {
		return err
	}

	// Check if payment can be refunded
	if payment.PaymentStatusEnum != "Verified" && payment.PaymentStatusEnum != "Succeed" {
		return errors.New("only verified or successful payments can be refunded")
	}

	// Check refund amount
	if amount > payment.Amount {
		return errors.New("refund amount cannot be greater than payment amount")
	}

	// In a real implementation, this would make an HTTP request to the payment gateway
	// For now, we'll just update the payment status
	payment.PaymentStatusEnum = "Refunded"
	payment.UpdatedAt = time.Now()

	return s.paymentRepo.Update(payment)
}

// GetPaymentStatus gets the current status of a payment
func (s *PaymentService) GetPaymentStatus(paymentID int64) (enums.StatusEnum, error) {
	// Get payment record
	payment, err := s.paymentRepo.GetByID(paymentID)
	if err != nil {
		return "", err
	}

	return payment.PaymentStatusEnum, nil
}

// GetGatewayByID gets gateway configuration by ID
func (s *PaymentService) GetGatewayByID(gatewayID int64) (*domain2.Gateway, error) {
	return s.gatewayRepo.GetByID(gatewayID)
}

// GetGatewayBySiteID gets gateway configuration by site ID
func (s *PaymentService) GetGatewayBySiteID(siteID int64) (*domain2.Gateway, error) {
	return s.gatewayRepo.GetBySiteID(siteID)
}

// Helper methods

// isGatewayActive checks if a gateway is active
func (s *PaymentService) isGatewayActive(gateway *domain2.Gateway, gatewayType int) bool {
	switch gatewayType {
	case 1: // Saman
		return gateway.IsActiveSaman == "Active"
	case 2: // Mellat
		return gateway.IsActiveMellat == "Active"
	case 3: // Parsian
		return gateway.IsActiveParsian == "Active"
	case 4: // Pasargad
		return gateway.IsActivePasargad == "Active"
	case 5: // IranKish
		return gateway.IsActiveIranKish == "Active"
	case 6: // Melli
		return gateway.IsActiveMelli == "Active"
	case 7: // AsanPardakht
		return gateway.IsActiveAsanPardakht == "Active"
	case 8: // Sepehr
		return gateway.IsActiveSepehr == "Active"
	case 9: // ZarinPal
		return gateway.IsActiveZarinPal == "Active"
	case 10: // PayIr
		return gateway.IsActivePayIr == "Active"
	case 11: // IdPay
		return gateway.IsActiveIdPay == "Active"
	case 12: // YekPay
		return gateway.IsActiveYekPay == "Active"
	case 13: // PayPing
		return gateway.IsActivePayPing == "Active"
	case 14: // ParbadVirtual
		return gateway.IsActiveParbadVirtual == "Active"
	default:
		return false
	}
}

// mapGatewayToString maps gateway ID to string
func (s *PaymentService) mapGatewayToString(gateway int64) string {
	switch gateway {
	case 1:
		return "Saman"
	case 2:
		return "Mellat"
	case 3:
		return "Parsian"
	case 4:
		return "Pasargad"
	case 5:
		return "IranKish"
	case 6:
		return "Melli"
	case 7:
		return "AsanPardakht"
	case 8:
		return "Sepehr"
	case 9:
		return "ZarinPal"
	case 10:
		return "PayIr"
	case 11:
		return "IdPay"
	case 12:
		return "YekPay"
	case 13:
		return "PayPing"
	case 14:
		return "ParbadVirtual"
	default:
		return "Unknown"
	}
}
