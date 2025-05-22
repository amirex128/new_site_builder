package paymentusecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/amirex128/new_site_builder/src/internal/domain/enums"
	"strconv"
	"time"

	"github.com/amirex128/new_site_builder/src/internal/application/usecase"
	"github.com/amirex128/new_site_builder/src/internal/contract/service"

	"github.com/gin-gonic/gin"

	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	"github.com/amirex128/new_site_builder/src/internal/application/dto/payment"
	"github.com/amirex128/new_site_builder/src/internal/contract"
	"github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/contract/repository"
	"github.com/amirex128/new_site_builder/src/internal/domain"
	"gorm.io/gorm"
)

type PaymentUsecase struct {
	*usecase.BaseUsecase
	logger      sflogger.Logger
	paymentRepo repository.IPaymentRepository
	gatewayRepo repository.IGatewayRepository
	authContext func(c *gin.Context) service.IAuthService
	container   contract.IContainer
	siteRepo    repository.ISiteRepository
}

func NewPaymentUsecase(c contract.IContainer) *PaymentUsecase {
	return &PaymentUsecase{
		BaseUsecase: &usecase.BaseUsecase{
			Logger: c.GetLogger(),
		},
		siteRepo:    c.GetSiteRepo(),
		paymentRepo: c.GetPaymentRepo(),
		gatewayRepo: c.GetGatewayRepo(),
		authContext: c.GetAuthTransientService(),
		container:   c,
	}
}

func (u *PaymentUsecase) VerifyPaymentCommand(params *payment.VerifyPaymentCommand) (any, error) {
	u.Logger.Info("VerifyPaymentCommand called", map[string]interface{}{
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
	u.Logger.Info("CreateOrUpdateGatewayCommand called", map[string]interface{}{
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
	userID, err := u.authContext(u.Ctx).GetUserID()
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
		if *params.IsActiveSaman == enums.ActiveStatus {
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
		if *params.IsActiveMellat == enums.ActiveStatus {
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
		if *params.IsActiveParsian == enums.ActiveStatus {
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
		if *params.IsActivePasargad == enums.ActiveStatus {
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
		if *params.IsActiveIranKish == enums.ActiveStatus {
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
		if *params.IsActiveMelli == enums.ActiveStatus {
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
		if *params.IsActiveAsanPardakht == enums.ActiveStatus {
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
		if *params.IsActiveSepehr == enums.ActiveStatus {
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
		if *params.IsActiveZarinPal == enums.ActiveStatus {
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
		if *params.IsActivePayIr == enums.ActiveStatus {
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
		if *params.IsActiveIdPay == enums.ActiveStatus {
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
		if *params.IsActiveYekPay == enums.ActiveStatus {
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
		if *params.IsActivePayPing == enums.ActiveStatus {
			gateway.IsActivePayPing = "Active"
		} else {
			gateway.IsActivePayPing = "Inactive"
		}
	}

	// ParbadVirtual Gateway
	if params.IsActiveParbadVirtual != nil {
		if *params.IsActiveParbadVirtual == enums.ActiveStatus {
			gateway.IsActiveParbadVirtual = "Active"
		} else {
			gateway.IsActiveParbadVirtual = "Inactive"
		}
	}
}

func (u *PaymentUsecase) GetByIdGatewayQuery(params *payment.GetByIdGatewayQuery) (any, error) {
	u.Logger.Info("GetByIdGatewayQuery called", map[string]interface{}{
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
	userID, err := u.authContext(u.Ctx).GetUserID()
	if err != nil {
		return nil, err
	}

	// Assuming site owners should be able to see their gateway
	// This logic might need adjustment based on your authorization requirements
	if gateway.UserID != userID {
		isAdmin, err := u.authContext(u.Ctx).IsAdmin()
		if err != nil || !isAdmin {
			return nil, errors.New("شما به این درگاه پرداخت دسترسی ندارید")
		}
	}

	return gateway, nil
}

func (u *PaymentUsecase) AdminGetAllGatewayQuery(params *payment.AdminGetAllGatewayQuery) (any, error) {
	u.Logger.Info("AdminGetAllGatewayQuery called", map[string]interface{}{
		"page":     params.Page,
		"pageSize": params.PageSize,
	})

	// Verify admin access
	isAdmin, err := u.authContext(u.Ctx).IsAdmin()
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
	u.Logger.Info("RequestGatewayCommand called", map[string]interface{}{
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
	gatewayAccountName := fmt.Sprintf("%s-%d", string(*params.Gateway), *params.SiteID)

	// Create payment record
	paymentData := domain.Payment{
		SiteID:             *params.SiteID,
		PaymentStatusEnum:  "Processing",
		UserType:           *params.UserType,
		TrackingNumber:     time.Now().UnixNano(),
		Gateway:            string(*params.Gateway),
		GatewayAccountName: gatewayAccountName,
		Amount:             *params.Amount,
		ServiceName:        *params.ServiceName,
		ServiceAction:      *params.ServiceAction,
		OrderID:            *params.OrderID,
		ReturnUrl:          *params.ReturnURL,
		CallVerifyUrl:      string(*params.CallVerifyURL),
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
	if *params.UserType == enums.CustomerTypeValue {
		paymentData.CustomerID = *params.UserID
		paymentData.UserID = 0 // Clear user ID for customers
	}

	// Request payment through payment gateway
	paymentURL, err := u.paymentRepo.RequestPayment(
		*params.Amount,
		*params.OrderID,
		*params.UserID,
		string(*params.Gateway),
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
	u.Logger.Info("AdminGetAllPaymentQuery called", map[string]interface{}{
		"page":     params.Page,
		"pageSize": params.PageSize,
	})

	// Verify admin access
	isAdmin, err := u.authContext(u.Ctx).IsAdmin()
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

func (u *PaymentUsecase) isGatewayActive(gateway domain.Gateway, gatewayType enums.PaymentGatewaysEnum) bool {
	switch gatewayType {
	case enums.SamanGatewayEnum:
		return gateway.IsActiveSaman == "Active"
	case enums.MellatGatewayEnum:
		return gateway.IsActiveMellat == "Active"
	case enums.ParsianGatewayEnum:
		return gateway.IsActiveParsian == "Active"
	case enums.PasargadGatewayEnum:
		return gateway.IsActivePasargad == "Active"
	case enums.IranKishGatewayEnum:
		return gateway.IsActiveIranKish == "Active"
	case enums.MelliGatewayEnum:
		return gateway.IsActiveMelli == "Active"
	case enums.AsanPardakhtGatewayEnum:
		return gateway.IsActiveAsanPardakht == "Active"
	case enums.SepehrGatewayEnum:
		return gateway.IsActiveSepehr == "Active"
	case enums.ZarinPalGatewayEnum:
		return gateway.IsActiveZarinPal == "Active"
	case enums.PayIrGatewayEnum:
		return gateway.IsActivePayIr == "Active"
	case enums.IdPayGatewayEnum:
		return gateway.IsActiveIdPay == "Active"
	case enums.YekPayGatewayEnum:
		return gateway.IsActiveYekPay == "Active"
	case enums.PayPingGatewayEnum:
		return gateway.IsActivePayPing == "Active"
	case enums.ParbadVirtualGatewayEnum:
		return gateway.IsActiveParbadVirtual == "Active"
	default:
		return false
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

	// Get user ID from order data
	userIDStr, ok := orderData["UserId"]
	if !ok {
		return false
	}

	// Parse user ID
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return false
	}

	// In a microservice architecture, this would call the User service
	// In our monolithic approach, we need to implement it directly

	// Get user repository from container
	// This is a simplified implementation - in a real app, you'd inject the repository
	userRepo := u.container.GetUserRepo()
	if userRepo == nil {
		return false
	}

	// Get the user
	user, err := userRepo.GetByID(userID)
	if err != nil {
		return false
	}

	// Get unit price repository
	unitPriceRepo := u.container.GetUnitPriceRepo()
	if unitPriceRepo == nil {
		return false
	}

	// Get all unit prices
	unitPrices, _, err := unitPriceRepo.GetAll(common.PaginationRequestDto{
		Page:     1,
		PageSize: 100,
	})
	if err != nil {
		return false
	}

	// Update user credits based on order data
	for _, unitPrice := range unitPrices {
		unitPriceNameKey := fmt.Sprintf("%s_UnitPriceName", unitPrice.Name)
		unitPriceCountKey := fmt.Sprintf("%s_UnitPriceCount", unitPrice.Name)

		if _, exists := orderData[unitPriceNameKey]; exists {
			countStr, ok := orderData[unitPriceCountKey]
			if !ok {
				continue
			}

			count, err := strconv.Atoi(countStr)
			if err != nil {
				continue
			}

			switch unitPrice.Name {
			case "SmsCredits":
				user.SmsCredits += count
			case "EmailCredits":
				user.EmailCredits += count
			case "AiCredits":
				user.AiCredits += count
			case "AiImageCredits":
				user.AiImageCredits += count
			case "StorageMbCredits":
				unitPriceDayKey := fmt.Sprintf("%s_UnitPriceDay", unitPrice.Name)
				expireDayStr, ok := orderData[unitPriceDayKey]
				if !ok {
					continue
				}

				expireDays, err := strconv.Atoi(expireDayStr)
				if err != nil {
					continue
				}

				user.StorageMbCredits += count
				expireAt := time.Now().AddDate(0, 0, expireDays)
				user.StorageMbCreditsExpireAt = &expireAt
			}
		}
	}

	// Update the user
	if err := userRepo.Update(user); err != nil {
		return false
	}

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

	// Get required data from order data
	userIDStr, ok := orderData["UserId"]
	if !ok {
		return false
	}

	planIDStr, ok := orderData["PlanId"]
	if !ok {
		return false
	}

	durationDaysStr, ok := orderData["DurationDays"]
	if !ok {
		return false
	}

	// Parse the values
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return false
	}

	planID, err := strconv.ParseInt(planIDStr, 10, 64)
	if err != nil {
		return false
	}

	durationDays, err := strconv.Atoi(durationDaysStr)
	if err != nil {
		return false
	}

	// Get repositories
	userRepo := u.container.GetUserRepo()
	if userRepo == nil {
		return false
	}

	planRepo := u.container.GetPlanRepo()
	if planRepo == nil {
		return false
	}

	// Get the user and plan
	user, err := userRepo.GetByID(userID)
	if err != nil {
		return false
	}

	plan, err := planRepo.GetByID(planID)
	if err != nil {
		return false
	}

	// Update user with plan details
	user.PlanID = &planID
	planExpiredAt := time.Now().AddDate(0, 0, durationDays)
	user.PlanExpiredAt = &planExpiredAt
	user.EmailCredits = plan.EmailCredits
	user.SmsCredits = plan.SmsCredits
	user.AiCredits = plan.AiCredits
	user.AiImageCredits = plan.AiImageCredits

	// Update storage credits if user doesn't have a plan yet
	if user.PlanID == nil {
		user.StorageMbCredits = plan.StorageMbCredits
		storageExpireAt := time.Now().AddDate(0, 0, durationDays)
		user.StorageMbCreditsExpireAt = &storageExpireAt
	}

	// Update user roles based on plan roles
	// This would require additional implementation to handle roles
	// For now, we'll skip this part

	// Update the user
	if err := userRepo.Update(user); err != nil {
		return false
	}

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

	// Parse order ID
	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil {
		return false
	}

	// Get order repository
	orderRepo := u.container.GetOrderRepo()
	if orderRepo == nil {
		return false
	}

	// Get the order with items
	order, err := orderRepo.GetByID(orderID)
	if err != nil {
		return false
	}

	// Update order status to paid
	order.OrderStatus = "Paid"
	order.UpdatedAt = time.Now()

	// Update the order
	if err := orderRepo.Update(order); err != nil {
		return false
	}

	// Get order items
	orderItems, _, err := u.container.GetOrderItemRepo().GetAllByOrderID(orderID, common.PaginationRequestDto{
		Page:     1,
		PageSize: 1000,
	})
	if err != nil {
		u.Logger.Error("Failed to get order items", map[string]interface{}{
			"error":   err.Error(),
			"orderID": orderID,
		})
		return false
	}

	// 1. Update product stock quantities
	productVariantRepo := u.container.GetProductVariantRepo()
	if productVariantRepo == nil {
		u.Logger.Error("Failed to get product variant repository", nil)
		return false
	}

	for _, item := range orderItems {
		if item.ProductVariantID > 0 {
			err := productVariantRepo.DecreaseStock(item.ProductVariantID, item.Quantity)
			if err != nil {
				u.Logger.Error("Failed to decrease product variant stock", map[string]interface{}{
					"error":            err.Error(),
					"productVariantID": item.ProductVariantID,
					"quantity":         item.Quantity,
				})
				// Continue with other items even if one fails
			}
		}
	}

	// 2. Decrease discount quantity if used
	if order.DiscountID != nil && *order.DiscountID > 0 {
		discountRepo := u.container.GetDiscountRepo()
		if discountRepo != nil {
			// Get customer ID
			customerID := order.CustomerID

			// Check if customer has already used this discount
			hasUsed, err := discountRepo.HasCustomerUsedDiscount(*order.DiscountID, customerID)
			if err != nil {
				u.Logger.Error("Failed to check if customer used discount", map[string]interface{}{
					"error":      err.Error(),
					"discountID": *order.DiscountID,
					"customerID": customerID,
				})
			} else if !hasUsed {
				// Decrease discount quantity
				if err := discountRepo.DecreaseQuantity(*order.DiscountID); err != nil {
					u.Logger.Error("Failed to decrease discount quantity", map[string]interface{}{
						"error":      err.Error(),
						"discountID": *order.DiscountID,
					})
				}

				// Record that this customer has used the discount
				if err := discountRepo.AddCustomerUsage(*order.DiscountID, customerID); err != nil {
					u.Logger.Error("Failed to record customer discount usage", map[string]interface{}{
						"error":      err.Error(),
						"discountID": *order.DiscountID,
						"customerID": customerID,
					})
				}
			}
		}
	}

	// 3. Decrease coupon quantity if applicable
	// First, get products associated with order items to check for coupons
	for _, item := range orderItems {
		if item.ProductID > 0 {
			// Get coupon for this product
			couponRepo := u.container.GetCouponRepo()
			if couponRepo != nil {
				coupon, err := couponRepo.GetByProductID(item.ProductID)
				if err == nil && coupon.ID > 0 && coupon.ExpiryDate.After(time.Now()) {
					// Coupon exists and is valid, decrease its quantity
					if err := couponRepo.DecreaseQuantity(coupon.ID); err != nil {
						u.Logger.Error("Failed to decrease coupon quantity", map[string]interface{}{
							"error":    err.Error(),
							"couponID": coupon.ID,
						})
					}
				}
			}
		}
	}

	// 4. Send confirmation email/notification
	// Get customer details
	customerRepo := u.container.GetCustomerRepo()
	if customerRepo != nil {
		customer, err := customerRepo.GetByID(order.CustomerID)
		if err == nil {
			// Get site details
			siteRepo := u.container.GetSiteRepo()
			if siteRepo != nil {
				site, err := siteRepo.GetByID(order.SiteID)
				if err == nil {
					// Prepare email content
					emailSubject := fmt.Sprintf("Order Confirmation #%d - %s", order.ID, site.Name)
					emailBody := fmt.Sprintf("Dear %s,\n\nThank you for your order #%d. Your payment has been successfully processed.\n\nOrder Details:\n",
						customer.FirstName, order.ID)

					// Add order items to email
					for _, item := range orderItems {
						// Get product details
						productRepo := u.container.GetProductRepo()
						if productRepo != nil {
							product, err := productRepo.GetByID(item.ProductID)
							if err == nil {
								emailBody += fmt.Sprintf("- %s (Quantity: %d) - %d\n",
									product.Name, item.Quantity, item.FinalPriceWithCouponDiscount)
							}
						}
					}

					emailBody += fmt.Sprintf("\nTotal: %d\n\nThank you for shopping with us!\n\n%s Team",
						order.TotalFinalPrice, site.Name)

					// In a real implementation, this would send an email
					// For now, we'll just log it
					u.Logger.Info("Order confirmation email would be sent", map[string]interface{}{
						"to":      customer.Email,
						"subject": emailSubject,
						"body":    emailBody,
					})
				}
			}
		}
	}

	// 5. Generate invoice
	// This would typically create a PDF invoice and store it
	// For now, we'll just log that an invoice would be generated
	u.Logger.Info("Invoice would be generated", map[string]interface{}{
		"orderID": order.ID,
		"amount":  order.TotalFinalPrice,
		"date":    time.Now().Format("2006-01-02"),
	})

	return true
}
