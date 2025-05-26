package paymentusecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/amirex128/new_site_builder/src/internal/domain/enums"

	"github.com/amirex128/new_site_builder/src/internal/application/usecase"
	"github.com/amirex128/new_site_builder/src/internal/contract/service"

	"github.com/gin-gonic/gin"

	"github.com/amirex128/new_site_builder/src/internal/application/dto/payment"
	"github.com/amirex128/new_site_builder/src/internal/application/utils/resp"
	"github.com/amirex128/new_site_builder/src/internal/contract"
	"github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/contract/repository"
	"github.com/amirex128/new_site_builder/src/internal/domain"
	"gorm.io/gorm"
)

type PaymentUsecase struct {
	*usecase.BaseUsecase
	paymentRepo repository.IPaymentRepository
	gatewayRepo repository.IGatewayRepository
	container   contract.IContainer
	siteRepo    repository.ISiteRepository
}

func NewPaymentUsecase(c contract.IContainer) *PaymentUsecase {
	return &PaymentUsecase{
		BaseUsecase: &usecase.BaseUsecase{
			Logger:      c.GetLogger(),
			AuthContext: c.GetAuthTransientService(),
		},
		siteRepo:    c.GetSiteRepo(),
		paymentRepo: c.GetPaymentRepo(),
		gatewayRepo: c.GetGatewayRepo(),
		container:   c,
	}
}

func (u *PaymentUsecase) VerifyPaymentCommand(params *payment.VerifyPaymentCommand) (*resp.Response, error) {
	trackingNumber, err := strconv.ParseInt(*params.TransactionCode, 10, 64)
	if err != nil {
		return nil, resp.NewError(resp.BadRequest, "invalid tracking number format")
	}
	paymentRecord, err := u.paymentRepo.GetByTrackingNumber(strconv.FormatInt(trackingNumber, 10))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, resp.NewError(resp.NotFound, "payment not found")
		}
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	isSuccess := *params.Result == "success" || *params.Result == "Success"
	if isSuccess {
		paymentRecord.PaymentStatusEnum = "Succeed"
		paymentRecord.Message = "Payment was successful"
		paymentRecord.TransactionCode = *params.TransactionCode
	} else {
		paymentRecord.PaymentStatusEnum = "Failed"
		paymentRecord.Message = "Payment failed"
	}
	paymentRecord.UpdatedAt = time.Now()
	err = u.paymentRepo.Update(paymentRecord)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	var verifyResult bool
	switch paymentRecord.CallVerifyUrl {
	case "ChargeCreditVerify":
		verifyResult = u.handleChargeCreditVerify(paymentRecord, isSuccess)
	case "UpgradePlanVerify":
		verifyResult = u.handleUpgradePlanVerify(paymentRecord, isSuccess)
	case "CreateOrderVerify":
		verifyResult = u.handleCreateOrderVerify(paymentRecord, isSuccess)
	default:
		verifyResult = false
	}
	responseData := map[string]interface{}{
		"isSuccess":        isSuccess,
		"message":          paymentRecord.Message,
		"serviceIsSuccess": verifyResult,
	}
	if paymentRecord.ReturnUrl != "" {
		responseData["redirectUrl"] = paymentRecord.ReturnUrl
	}
	return resp.NewResponseData(resp.Retrieved, responseData, "وضعیت پرداخت"), nil
}

func (u *PaymentUsecase) CreateOrUpdateGatewayCommand(params *payment.CreateOrUpdateGatewayCommand) (*resp.Response, error) {
	existingGateway, err := u.gatewayRepo.GetBySiteID(*params.SiteID)
	isCreating := false
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			isCreating = true
		} else {
			return nil, resp.NewError(resp.Internal, err.Error())
		}
	}
	userID, err := u.AuthContext(u.Ctx).GetUserID()
	if err != nil {
		return nil, resp.NewError(resp.Unauthorized, err.Error())
	}
	gateway := domain.Gateway{
		SiteID:    *params.SiteID,
		UserID:    *userID,
		UpdatedAt: time.Now(),
	}
	if isCreating {
		gateway.CreatedAt = time.Now()
	} else {
		gateway.ID = existingGateway.ID
		gateway.CreatedAt = existingGateway.CreatedAt
	}
	u.setGatewayValues(&gateway, params)
	if isCreating {
		err = u.gatewayRepo.Create(&gateway)
	} else {
		err = u.gatewayRepo.Update(&gateway)
	}
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	return resp.NewResponseData(resp.Success, map[string]interface{}{"gateway": gateway}, "درگاه با موفقیت ایجاد یا بروزرسانی شد"), nil
}

func (u *PaymentUsecase) setGatewayValues(gateway *domain.Gateway, params *payment.CreateOrUpdateGatewayCommand) {
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
	if params.IsActiveParbadVirtual != nil {
		if *params.IsActiveParbadVirtual == enums.ActiveStatus {
			gateway.IsActiveParbadVirtual = "Active"
		} else {
			gateway.IsActiveParbadVirtual = "Inactive"
		}
	}
}

func (u *PaymentUsecase) GetByIdGatewayQuery(params *payment.GetByIdGatewayQuery) (*resp.Response, error) {
	gateway, err := u.gatewayRepo.GetByID(*params.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, resp.NewError(resp.NotFound, "درگاه پرداخت یافت نشد")
		}
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	userID, err := u.AuthContext(u.Ctx).GetUserID()
	if err != nil {
		return nil, resp.NewError(resp.Unauthorized, err.Error())
	}
	if userID != nil && gateway.UserID != *userID {
		isAdmin, err := u.AuthContext(u.Ctx).IsAdmin()
		if err != nil || !isAdmin {
			return nil, resp.NewError(resp.Unauthorized, "شما به این درگاه پرداخت دسترسی ندارید")
		}
	}
	return resp.NewResponseData(resp.Retrieved, map[string]interface{}{"gateway": gateway}, "درگاه با موفقیت دریافت شد"), nil
}

func (u *PaymentUsecase) AdminGetAllGatewayQuery(params *payment.AdminGetAllGatewayQuery) (*resp.Response, error) {
	isAdmin, err := u.AuthContext(u.Ctx).IsAdmin()
	if err != nil || !isAdmin {
		return nil, resp.NewError(resp.Unauthorized, err.Error())
	}
	gatewaysResult, err := u.gatewayRepo.GetAll(params.PaginationRequestDto)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	return resp.NewResponseData(resp.Retrieved, map[string]interface{}{
		"items":     gatewaysResult.Items,
		"total":     gatewaysResult.TotalCount,
		"page":      gatewaysResult.PageNumber,
		"pageSize":  params.PageSize,
		"totalPage": gatewaysResult.TotalPages,
	}, "لیست درگاه ها با موفقیت دریافت شد"), nil
}

func (u *PaymentUsecase) RequestGatewayCommand(params *payment.RequestGatewayCommand) (*resp.Response, error) {
	gateway, err := u.gatewayRepo.GetBySiteID(*params.SiteID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, resp.NewError(resp.NotFound, "درگاه پرداخت برای این سایت یافت نشد")
		}
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	if !u.isGatewayActive(gateway, *params.Gateway) {
		return nil, resp.NewError(resp.BadRequest, "درگاه پرداخت انتخاب شده غیرفعال است")
	}
	gatewayAccountName := fmt.Sprintf("%s-%d", string(*params.Gateway), *params.SiteID)
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
		CustomerID:         0,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}
	orderDataJSON, err := json.Marshal(params.OrderData)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	paymentData.OrderData = string(orderDataJSON)
	if *params.UserType == enums.CustomerTypeValue {
		paymentData.CustomerID = *params.UserID
		paymentData.UserID = 0
	}
	paymentURL, err := u.paymentRepo.RequestPayment(
		*params.Amount,
		*params.OrderID,
		*params.UserID,
		string(*params.Gateway),
		params.OrderData,
	)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}
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
	return resp.NewResponseData(resp.Retrieved, responseData, "درخواست پرداخت با موفقیت انجام شد"), nil
}

func (u *PaymentUsecase) AdminGetAllPaymentQuery(params *payment.AdminGetAllPaymentQuery) (*resp.Response, error) {
	isAdmin, err := u.AuthContext(u.Ctx).IsAdmin()
	if err != nil || !isAdmin {
		return nil, resp.NewError(resp.Unauthorized, err.Error())
	}
	paymentsResult, err := u.paymentRepo.GetAll(params.PaginationRequestDto)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	return resp.NewResponseData(resp.Retrieved, map[string]interface{}{
		"items":     paymentsResult.Items,
		"total":     paymentsResult.TotalCount,
		"page":      paymentsResult.PageNumber,
		"pageSize":  params.PageSize,
		"totalPage": paymentsResult.TotalPages,
	}, "لیست پرداخت ها با موفقیت دریافت شد"), nil
}

func (u *PaymentUsecase) isGatewayActive(gateway *domain.Gateway, gatewayType enums.PaymentGatewaysEnum) bool {
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

func (u *PaymentUsecase) handleChargeCreditVerify(payment *domain.Payment, isSuccess bool) bool {
	if !isSuccess {
		return false
	}
	var orderData map[string]string
	err := json.Unmarshal([]byte(payment.OrderData), &orderData)
	if err != nil {
		return false
	}
	userIDStr, ok := orderData["UserId"]
	if !ok {
		return false
	}
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return false
	}
	userRepo := u.container.GetUserRepo()
	if userRepo == nil {
		return false
	}
	user, err := userRepo.GetByID(userID)
	if err != nil {
		return false
	}
	unitPriceRepo := u.container.GetUnitPriceRepo()
	if unitPriceRepo == nil {
		return false
	}
	unitPricesResult, err := unitPriceRepo.GetAll(common.PaginationRequestDto{
		Page:     1,
		PageSize: 100,
	})
	if err != nil {
		return false
	}
	for _, unitPrice := range unitPricesResult.Items {
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
	if err := userRepo.Update(user); err != nil {
		return false
	}
	return true
}

func (u *PaymentUsecase) handleUpgradePlanVerify(payment *domain.Payment, isSuccess bool) bool {
	if !isSuccess {
		return false
	}
	var orderData map[string]string
	err := json.Unmarshal([]byte(payment.OrderData), &orderData)
	if err != nil {
		return false
	}
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
	userRepo := u.container.GetUserRepo()
	if userRepo == nil {
		return false
	}
	planRepo := u.container.GetPlanRepo()
	if planRepo == nil {
		return false
	}
	user, err := userRepo.GetByID(userID)
	if err != nil {
		return false
	}
	plan, err := planRepo.GetByID(planID)
	if err != nil {
		return false
	}
	user.PlanID = &planID
	planExpiredAt := time.Now().AddDate(0, 0, durationDays)
	user.PlanExpiredAt = &planExpiredAt
	user.EmailCredits = plan.EmailCredits
	user.SmsCredits = plan.SmsCredits
	user.AiCredits = plan.AiCredits
	user.AiImageCredits = plan.AiImageCredits
	if user.PlanID == nil {
		user.StorageMbCredits = plan.StorageMbCredits
		storageExpireAt := time.Now().AddDate(0, 0, durationDays)
		user.StorageMbCreditsExpireAt = &storageExpireAt
	}
	if err := userRepo.Update(user); err != nil {
		return false
	}
	return true
}

func (u *PaymentUsecase) handleCreateOrderVerify(payment *domain.Payment, isSuccess bool) bool {
	if !isSuccess {
		return false
	}
	var orderData map[string]string
	err := json.Unmarshal([]byte(payment.OrderData), &orderData)
	if err != nil {
		return false
	}
	orderIDStr, ok := orderData["OrderId"]
	if !ok {
		return false
	}
	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil {
		return false
	}
	orderRepo := u.container.GetOrderRepo()
	if orderRepo == nil {
		return false
	}
	order, err := orderRepo.GetByID(orderID)
	if err != nil {
		return false
	}
	order.OrderStatus = "Paid"
	order.UpdatedAt = time.Now()
	if err := orderRepo.Update(order); err != nil {
		return false
	}
	orderItemsResult, err := u.container.GetOrderItemRepo().GetAllByOrderID(orderID, common.PaginationRequestDto{
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
	productVariantRepo := u.container.GetProductVariantRepo()
	if productVariantRepo == nil {
		u.Logger.Error("Failed to get product variant repository", nil)
		return false
	}
	for _, item := range orderItemsResult.Items {
		if item.ProductVariantID > 0 {
			err := productVariantRepo.DecreaseStock(item.ProductVariantID, item.Quantity)
			if err != nil {
				u.Logger.Error("Failed to decrease product variant stock", map[string]interface{}{
					"error":            err.Error(),
					"productVariantID": item.ProductVariantID,
					"quantity":         item.Quantity,
				})
			}
		}
	}
	if order.DiscountID != nil && *order.DiscountID > 0 {
		discountRepo := u.container.GetDiscountRepo()
		if discountRepo != nil {
			customerID := order.CustomerID
			hasUsed, err := discountRepo.HasCustomerUsedDiscount(*order.DiscountID, customerID)
			if err != nil {
				u.Logger.Error("Failed to check if customer used discount", map[string]interface{}{
					"error":      err.Error(),
					"discountID": *order.DiscountID,
					"customerID": customerID,
				})
			} else if !hasUsed {
				if err := discountRepo.DecreaseQuantity(*order.DiscountID); err != nil {
					u.Logger.Error("Failed to decrease discount quantity", map[string]interface{}{
						"error":      err.Error(),
						"discountID": *order.DiscountID,
					})
				}
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
	for _, item := range orderItemsResult.Items {
		if item.ProductID > 0 {
			couponRepo := u.container.GetCouponRepo()
			if couponRepo != nil {
				coupon, err := couponRepo.GetByProductID(item.ProductID)
				if err == nil && coupon.ID > 0 && coupon.ExpiryDate.After(time.Now()) {
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
	customerRepo := u.container.GetCustomerRepo()
	if customerRepo != nil {
		customer, err := customerRepo.GetByID(order.CustomerID)
		if err == nil {
			siteRepo := u.container.GetSiteRepo()
			if siteRepo != nil {
				site, err := siteRepo.GetByID(order.SiteID)
				if err == nil {
					emailSubject := fmt.Sprintf("Order Confirmation #%d - %s", order.ID, site.Name)
					emailBody := fmt.Sprintf("Dear %s,\n\nThank you for your order #%d. Your payment has been successfully processed.\n\nOrder Details:\n",
						customer.FirstName, order.ID)
					for _, item := range orderItemsResult.Items {
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
					u.Logger.Info("Order confirmation email would be sent", map[string]interface{}{
						"to":      customer.Email,
						"subject": emailSubject,
						"body":    emailBody,
					})
				}
			}
		}
	}
	u.Logger.Info("Invoice would be generated", map[string]interface{}{
		"orderID": order.ID,
		"amount":  order.TotalFinalPrice,
		"date":    time.Now().Format("2006-01-02"),
	})
	return true
}
