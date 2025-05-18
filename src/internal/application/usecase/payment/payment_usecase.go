package paymentusecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	"github.com/amirex128/new_site_builder/src/internal/application/dto/payment"
	"github.com/amirex128/new_site_builder/src/internal/contract"
	"github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/contract/repository"
	"github.com/amirex128/new_site_builder/src/internal/domain"
	"gorm.io/gorm"
)

type PaymentUsecase struct {
	logger         sflogger.Logger
	paymentRepo    repository.IPaymentRepository
	gatewayRepo    repository.IGatewayRepository
	authContextSvc common.IAuthContextService
}

func NewPaymentUsecase(c contract.IContainer) *PaymentUsecase {
	return &PaymentUsecase{
		logger:         c.GetLogger(),
		paymentRepo:    c.GetPaymentRepo(),
		gatewayRepo:    c.GetGatewayRepo(),
		authContextSvc: c.GetAuthContextService(),
	}
}

func (u *PaymentUsecase) VerifyPaymentCommand(params *payment.VerifyPaymentCommand) (any, error) {
	u.logger.Info("VerifyPaymentCommand called", map[string]interface{}{
		"transactionCode": *params.TransactionCode,
	})

	// Find payment by tracking number or transaction code
	trackingNumber, err := strconv.ParseInt(*params.TransactionCode, 10, 64)
	if err != nil {
		// If not a number, try to find by transaction code
		return nil, errors.New("invalid tracking number format")
	}

	paymentRecord, err := u.paymentRepo.GetByTrackingNumber(strconv.FormatInt(trackingNumber, 10))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("payment not found")
		}
		return nil, err
	}

	// Process payment state from result
	isSuccess := *params.Result == "success" || *params.Result == "Success"

	// Update payment status
	if isSuccess {
		paymentRecord.PaymentStatusEnum = "Succeed"
		paymentRecord.Message = "Payment was successful"
		paymentRecord.TransactionCode = *params.TransactionCode
	} else {
		paymentRecord.PaymentStatusEnum = "Failed"
		paymentRecord.Message = "Payment failed"
	}

	paymentRecord.UpdatedAt = time.Now()

	// Save updated payment
	err = u.paymentRepo.Update(paymentRecord)
	if err != nil {
		return nil, err
	}

	// Process verification callbacks based on CallVerifyUrl
	// In a microservice architecture, this would be handled by separate services
	// In our monolithic approach, we handle different cases here
	var verifyResult bool

	switch paymentRecord.CallVerifyUrl {
	case "ChargeCreditVerify":
		// Implement credit charging logic here
		verifyResult = u.handleChargeCreditVerify(paymentRecord, isSuccess)
	case "UpgradePlanVerify":
		// Implement plan upgrade logic here
		verifyResult = u.handleUpgradePlanVerify(paymentRecord, isSuccess)
	case "CreateOrderVerify":
		// Implement order finalization logic here
		verifyResult = u.handleCreateOrderVerify(paymentRecord, isSuccess)
	default:
		verifyResult = false
	}

	// Build response
	responseData := map[string]interface{}{
		"isSuccess": isSuccess,
		"message":   paymentRecord.Message,
	}

	// Add service verification result
	responseData["serviceIsSuccess"] = verifyResult

	// Redirect info
	if paymentRecord.ReturnUrl != "" {
		responseData["redirectUrl"] = paymentRecord.ReturnUrl
	}

	return responseData, nil
}

func (u *PaymentUsecase) CreateOrUpdateGatewayCommand(params *payment.CreateOrUpdateGatewayCommand) (any, error) {
	u.logger.Info("CreateOrUpdateGatewayCommand called", map[string]interface{}{
		"siteId": *params.SiteID,
	})

	// Check if a gateway exists for this site
	existingGateway, err := u.gatewayRepo.GetBySiteID(*params.SiteID)
	isCreating := false

	// If gateway not found, prepare to create a new one
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			isCreating = true
		} else {
			return nil, err
		}
	}

	// Get current user ID
	userID, err := u.authContextSvc.GetUserID()
	if err != nil {
		return nil, err
	}

	// Prepare the gateway object
	gateway := domain.Gateway{
		SiteID:    *params.SiteID,
		UserID:    userID,
		UpdatedAt: time.Now(),
	}

	if isCreating {
		gateway.CreatedAt = time.Now()
	} else {
		gateway.ID = existingGateway.ID
		gateway.CreatedAt = existingGateway.CreatedAt
	}

	// Set all gateway configuration values
	u.setGatewayValues(&gateway, params)

	// Create or update the gateway in the database
	if isCreating {
		err = u.gatewayRepo.Create(gateway)
	} else {
		err = u.gatewayRepo.Update(gateway)
	}

	if err != nil {
		return nil, err
	}

	// Return the gateway object
	return gateway, nil
}

func (u *PaymentUsecase) setGatewayValues(gateway *domain.Gateway, params *payment.CreateOrUpdateGatewayCommand) {
	// Saman Gateway
	if params.Saman != nil && params.Saman.MerchantID != nil && params.Saman.Password != nil {
		gateway.SamanMerchantId = *params.Saman.MerchantID
		gateway.SamanPassword = *params.Saman.Password
	}
	if params.IsActiveSaman != nil {
		if *params.IsActiveSaman == payment.Active {
			gateway.IsActiveSaman = "Active"
		} else {
			gateway.IsActiveSaman = "Inactive"
		}
	}

	// Mellat Gateway
	if params.Mellat != nil && params.Mellat.TerminalID != nil && params.Mellat.UserName != nil && params.Mellat.UserPassword != nil {
		gateway.MellatTerminalId = params.Mellat.TerminalID
		gateway.MellatUserName = *params.Mellat.UserName
		gateway.MellatUserPassword = *params.Mellat.UserPassword
	}
	if params.IsActiveMellat != nil {
		if *params.IsActiveMellat == payment.Active {
			gateway.IsActiveMellat = "Active"
		} else {
			gateway.IsActiveMellat = "Inactive"
		}
	}

	// Parsian Gateway
	if params.Parsian != nil && params.Parsian.LoginAccount != nil {
		gateway.ParsianLoginAccount = *params.Parsian.LoginAccount
	}
	if params.IsActiveParsian != nil {
		if *params.IsActiveParsian == payment.Active {
			gateway.IsActiveParsian = "Active"
		} else {
			gateway.IsActiveParsian = "Inactive"
		}
	}

	// Pasargad Gateway
	if params.Pasargad != nil && params.Pasargad.MerchantCode != nil && params.Pasargad.TerminalCode != nil && params.Pasargad.PrivateKey != nil {
		gateway.PasargadMerchantCode = *params.Pasargad.MerchantCode
		gateway.PasargadTerminalCode = *params.Pasargad.TerminalCode
		gateway.PasargadPrivateKey = *params.Pasargad.PrivateKey
	}
	if params.IsActivePasargad != nil {
		if *params.IsActivePasargad == payment.Active {
			gateway.IsActivePasargad = "Active"
		} else {
			gateway.IsActivePasargad = "Inactive"
		}
	}

	// IranKish Gateway
	if params.IranKish != nil && params.IranKish.TerminalID != nil && params.IranKish.AcceptorID != nil && params.IranKish.PassPhrase != nil && params.IranKish.PublicKey != nil {
		gateway.IranKishTerminalId = *params.IranKish.TerminalID
		gateway.IranKishAcceptorId = *params.IranKish.AcceptorID
		gateway.IranKishPassPhrase = *params.IranKish.PassPhrase
		gateway.IranKishPublicKey = *params.IranKish.PublicKey
	}
	if params.IsActiveIranKish != nil {
		if *params.IsActiveIranKish == payment.Active {
			gateway.IsActiveIranKish = "Active"
		} else {
			gateway.IsActiveIranKish = "Inactive"
		}
	}

	// Melli Gateway
	if params.Melli != nil && params.Melli.TerminalID != nil && params.Melli.MerchantID != nil && params.Melli.TerminalKey != nil {
		gateway.MelliTerminalId = *params.Melli.TerminalID
		gateway.MelliMerchantId = *params.Melli.MerchantID
		gateway.MelliTerminalKey = *params.Melli.TerminalKey
	}
	if params.IsActiveMelli != nil {
		if *params.IsActiveMelli == payment.Active {
			gateway.IsActiveMelli = "Active"
		} else {
			gateway.IsActiveMelli = "Inactive"
		}
	}

	// AsanPardakht Gateway
	if params.AsanPardakht != nil && params.AsanPardakht.MerchantConfigurationID != nil && params.AsanPardakht.UserName != nil && params.AsanPardakht.Password != nil && params.AsanPardakht.Key != nil && params.AsanPardakht.IV != nil {
		gateway.AsanPardakhtMerchantConfigurationId = *params.AsanPardakht.MerchantConfigurationID
		gateway.AsanPardakhtUserName = *params.AsanPardakht.UserName
		gateway.AsanPardakhtPassword = *params.AsanPardakht.Password
		gateway.AsanPardakhtKey = *params.AsanPardakht.Key
		gateway.AsanPardakhtIV = *params.AsanPardakht.IV
	}
	if params.IsActiveAsanPardakht != nil {
		if *params.IsActiveAsanPardakht == payment.Active {
			gateway.IsActiveAsanPardakht = "Active"
		} else {
			gateway.IsActiveAsanPardakht = "Inactive"
		}
	}

	// Sepehr Gateway
	if params.Sepehr != nil && params.Sepehr.TerminalID != nil {
		gateway.SepehrTerminalId = params.Sepehr.TerminalID
	}
	if params.IsActiveSepehr != nil {
		if *params.IsActiveSepehr == payment.Active {
			gateway.IsActiveSepehr = "Active"
		} else {
			gateway.IsActiveSepehr = "Inactive"
		}
	}

	// ZarinPal Gateway
	if params.ZarinPal != nil && params.ZarinPal.MerchantID != nil {
		gateway.ZarinPalMerchantId = *params.ZarinPal.MerchantID
		if params.ZarinPal.AuthorizationToken != nil {
			gateway.ZarinPalAuthorizationToken = *params.ZarinPal.AuthorizationToken
		}
		if params.ZarinPal.IsSandbox != nil {
			gateway.ZarinPalIsSandbox = params.ZarinPal.IsSandbox
		}
	}
	if params.IsActiveZarinPal != nil {
		if *params.IsActiveZarinPal == payment.Active {
			gateway.IsActiveZarinPal = "Active"
		} else {
			gateway.IsActiveZarinPal = "Inactive"
		}
	}

	// PayIr Gateway
	if params.PayIr != nil && params.PayIr.API != nil {
		gateway.PayIrApi = *params.PayIr.API
		if params.PayIr.IsTestAccount != nil {
			gateway.PayIrIsTestAccount = params.PayIr.IsTestAccount
		}
	}
	if params.IsActivePayIr != nil {
		if *params.IsActivePayIr == payment.Active {
			gateway.IsActivePayIr = "Active"
		} else {
			gateway.IsActivePayIr = "Inactive"
		}
	}

	// IdPay Gateway
	if params.IdPay != nil && params.IdPay.API != nil {
		gateway.IdPayApi = *params.IdPay.API
		if params.IdPay.IsTestAccount != nil {
			gateway.IdPayIsTestAccount = params.IdPay.IsTestAccount
		}
	}
	if params.IsActiveIdPay != nil {
		if *params.IsActiveIdPay == payment.Active {
			gateway.IsActiveIdPay = "Active"
		} else {
			gateway.IsActiveIdPay = "Inactive"
		}
	}

	// YekPay Gateway
	if params.YekPay != nil && params.YekPay.MerchantID != nil {
		gateway.YekPayMerchantId = *params.YekPay.MerchantID
	}
	if params.IsActiveYekPay != nil {
		if *params.IsActiveYekPay == payment.Active {
			gateway.IsActiveYekPay = "Active"
		} else {
			gateway.IsActiveYekPay = "Inactive"
		}
	}

	// PayPing Gateway
	if params.PayPing != nil && params.PayPing.AccessToken != nil {
		gateway.PayPingAccessToken = *params.PayPing.AccessToken
	}
	if params.IsActivePayPing != nil {
		if *params.IsActivePayPing == payment.Active {
			gateway.IsActivePayPing = "Active"
		} else {
			gateway.IsActivePayPing = "Inactive"
		}
	}

	// ParbadVirtual Gateway
	if params.IsActiveParbadVirtual != nil {
		if *params.IsActiveParbadVirtual == payment.Active {
			gateway.IsActiveParbadVirtual = "Active"
		} else {
			gateway.IsActiveParbadVirtual = "Inactive"
		}
	}
}

func (u *PaymentUsecase) GetByIdGatewayQuery(params *payment.GetByIdGatewayQuery) (any, error) {
	u.logger.Info("GetByIdGatewayQuery called", map[string]interface{}{
		"id": *params.ID,
	})

	gateway, err := u.gatewayRepo.GetByID(*params.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("درگاه پرداخت یافت نشد")
		}
		return nil, err
	}

	// Check user access
	userID, err := u.authContextSvc.GetUserID()
	if err != nil {
		return nil, err
	}

	// Assuming site owners should be able to see their gateway
	// This logic might need adjustment based on your authorization requirements
	if gateway.UserID != userID {
		isAdmin, err := u.authContextSvc.IsAdmin()
		if err != nil || !isAdmin {
			return nil, errors.New("شما به این درگاه پرداخت دسترسی ندارید")
		}
	}

	return gateway, nil
}

func (u *PaymentUsecase) AdminGetAllGatewayQuery(params *payment.AdminGetAllGatewayQuery) (any, error) {
	u.logger.Info("AdminGetAllGatewayQuery called", map[string]interface{}{
		"page":     params.Page,
		"pageSize": params.PageSize,
	})

	// Verify admin access
	isAdmin, err := u.authContextSvc.IsAdmin()
	if err != nil {
		return nil, err
	}
	if !isAdmin {
		return nil, errors.New("فقط مدیران سیستم مجاز به دسترسی به این بخش هستند")
	}

	gateways, count, err := u.gatewayRepo.GetAll(params.PaginationRequestDto)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"items": gateways,
		"total": count,
		"page":  params.Page,
		"size":  params.PageSize,
	}, nil
}

func (u *PaymentUsecase) RequestGatewayCommand(params *payment.RequestGatewayCommand) (any, error) {
	u.logger.Info("RequestGatewayCommand called", map[string]interface{}{
		"siteId":   *params.SiteID,
		"orderId":  *params.OrderID,
		"amount":   *params.Amount,
		"gateway":  *params.Gateway,
		"userType": *params.UserType,
		"userId":   *params.UserID,
	})

	// Check if we have an active gateway for this site and payment method
	gateway, err := u.gatewayRepo.GetBySiteID(*params.SiteID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("درگاه پرداخت برای این سایت یافت نشد")
		}
		return nil, err
	}

	// Check if selected gateway is active
	if !u.isGatewayActive(gateway, *params.Gateway) {
		return nil, errors.New("درگاه پرداخت انتخاب شده غیرفعال است")
	}

	// Generate gateway account name
	gatewayAccountName := fmt.Sprintf("%d-%d", *params.Gateway, *params.SiteID)

	// Create payment record
	paymentData := domain.Payment{
		SiteID:             *params.SiteID,
		PaymentStatusEnum:  "Processing",
		UserType:           u.mapUserTypeToString(*params.UserType),
		TrackingNumber:     time.Now().UnixNano(),
		Gateway:            u.mapGatewayToString(*params.Gateway),
		GatewayAccountName: gatewayAccountName,
		Amount:             *params.Amount,
		ServiceName:        *params.ServiceName,
		ServiceAction:      *params.ServiceAction,
		OrderID:            *params.OrderID,
		ReturnUrl:          *params.ReturnURL,
		CallVerifyUrl:      u.mapVerifyEndpointToString(*params.CallVerifyURL),
		ClientIp:           *params.ClientIP,
		UserID:             *params.UserID,
		CustomerID:         0, // Set appropriate customer ID if applicable
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	// Serialize order data to JSON string
	orderDataJSON, err := json.Marshal(params.OrderData)
	if err != nil {
		return nil, err
	}
	paymentData.OrderData = string(orderDataJSON)

	// Set CustomerID if user type is Customer
	if *params.UserType == payment.Customer {
		paymentData.CustomerID = *params.UserID
		paymentData.UserID = 0 // Clear user ID for customers
	}

	// Request payment through payment gateway
	paymentURL, err := u.paymentRepo.RequestPayment(
		*params.Amount,
		*params.OrderID,
		*params.UserID,
		strconv.Itoa(int(*params.Gateway)),
		params.OrderData,
	)
	if err != nil {
		return nil, err
	}

	// Create response object
	responseData := map[string]interface{}{
		"responseStatus": map[string]interface{}{
			"isSuccess": true,
			"message":   "Payment request successful",
		},
		"isSuccess":      true,
		"message":        "Payment request successful",
		"status":         "Processing",
		"trackingNumber": paymentData.TrackingNumber,
		"url":            paymentURL,
		"method":         "GET",
		"formData":       map[string]string{},
	}

	return responseData, nil
}

func (u *PaymentUsecase) AdminGetAllPaymentQuery(params *payment.AdminGetAllPaymentQuery) (any, error) {
	u.logger.Info("AdminGetAllPaymentQuery called", map[string]interface{}{
		"page":     params.Page,
		"pageSize": params.PageSize,
	})

	// Verify admin access
	isAdmin, err := u.authContextSvc.IsAdmin()
	if err != nil {
		return nil, err
	}
	if !isAdmin {
		return nil, errors.New("فقط مدیران سیستم مجاز به دسترسی به این بخش هستند")
	}

	payments, count, err := u.paymentRepo.GetAll(params.PaginationRequestDto)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"items": payments,
		"total": count,
		"page":  params.Page,
		"size":  params.PageSize,
	}, nil
}

// Helper methods

func (u *PaymentUsecase) isGatewayActive(gateway domain.Gateway, gatewayType payment.PaymentGatewaysEnum) bool {
	switch gatewayType {
	case payment.Saman:
		return gateway.IsActiveSaman == "Active"
	case payment.Mellat:
		return gateway.IsActiveMellat == "Active"
	case payment.Parsian:
		return gateway.IsActiveParsian == "Active"
	case payment.Pasargad:
		return gateway.IsActivePasargad == "Active"
	case payment.IranKish:
		return gateway.IsActiveIranKish == "Active"
	case payment.Melli:
		return gateway.IsActiveMelli == "Active"
	case payment.AsanPardakht:
		return gateway.IsActiveAsanPardakht == "Active"
	case payment.Sepehr:
		return gateway.IsActiveSepehr == "Active"
	case payment.ZarinPal:
		return gateway.IsActiveZarinPal == "Active"
	case payment.PayIr:
		return gateway.IsActivePayIr == "Active"
	case payment.IdPay:
		return gateway.IsActiveIdPay == "Active"
	case payment.YekPay:
		return gateway.IsActiveYekPay == "Active"
	case payment.PayPing:
		return gateway.IsActivePayPing == "Active"
	case payment.ParbadVirtual:
		return gateway.IsActiveParbadVirtual == "Active"
	default:
		return false
	}
}

func (u *PaymentUsecase) mapUserTypeToString(userType payment.UserTypeEnum) string {
	switch userType {
	case payment.User:
		return "User"
	case payment.Customer:
		return "Customer"
	case payment.Guest:
		return "Guest"
	default:
		return "Unknown"
	}
}

func (u *PaymentUsecase) mapGatewayToString(gateway payment.PaymentGatewaysEnum) string {
	switch gateway {
	case payment.Saman:
		return "Saman"
	case payment.Mellat:
		return "Mellat"
	case payment.Parsian:
		return "Parsian"
	case payment.Pasargad:
		return "Pasargad"
	case payment.IranKish:
		return "IranKish"
	case payment.Melli:
		return "Melli"
	case payment.AsanPardakht:
		return "AsanPardakht"
	case payment.Sepehr:
		return "Sepehr"
	case payment.ZarinPal:
		return "ZarinPal"
	case payment.PayIr:
		return "PayIr"
	case payment.IdPay:
		return "IdPay"
	case payment.YekPay:
		return "YekPay"
	case payment.PayPing:
		return "PayPing"
	case payment.ParbadVirtual:
		return "ParbadVirtual"
	default:
		return "Unknown"
	}
}

func (u *PaymentUsecase) mapVerifyEndpointToString(endpoint payment.VerifyPaymentEndpointEnum) string {
	switch endpoint {
	case payment.ChargeCreditVerify:
		return "ChargeCreditVerify"
	case payment.UpgradePlanVerify:
		return "UpgradePlanVerify"
	case payment.CreateOrderVerify:
		return "CreateOrderVerify"
	default:
		return "Unknown"
	}
}

// Handler methods for different verification endpoints

func (u *PaymentUsecase) handleChargeCreditVerify(payment domain.Payment, isSuccess bool) bool {
	if !isSuccess {
		return false
	}

	// Parse order data
	var orderData map[string]string
	err := json.Unmarshal([]byte(payment.OrderData), &orderData)
	if err != nil {
		return false
	}

	// In a microservice architecture, this would call a credit service
	// In our monolithic approach, we can implement it directly here

	// TODO: Implement credit charging logic
	// This could involve updating user credits in the database

	return true
}

func (u *PaymentUsecase) handleUpgradePlanVerify(payment domain.Payment, isSuccess bool) bool {
	if !isSuccess {
		return false
	}

	// Parse order data
	var orderData map[string]string
	err := json.Unmarshal([]byte(payment.OrderData), &orderData)
	if err != nil {
		return false
	}

	// TODO: Implement plan upgrade logic
	// This could involve updating user plan in the database

	return true
}

func (u *PaymentUsecase) handleCreateOrderVerify(payment domain.Payment, isSuccess bool) bool {
	if !isSuccess {
		return false
	}

	// Parse order data
	var orderData map[string]string
	err := json.Unmarshal([]byte(payment.OrderData), &orderData)
	if err != nil {
		return false
	}

	// Get order ID from order data
	orderIDStr, ok := orderData["OrderId"]
	if !ok {
		return false
	}

	// Parse order ID but don't use it directly since we're in a monolithic implementation
	_, err = strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil {
		return false
	}

	// TODO: Implement order finalization logic
	// This could involve updating order status in the database

	return true
}
