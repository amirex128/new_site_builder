package userusecase

import (
	"fmt"
	"strconv"
	"time"

	"github.com/amirex128/new_site_builder/src/internal/application/utils"
	"github.com/amirex128/new_site_builder/src/internal/application/utils/resp"
	"github.com/amirex128/new_site_builder/src/internal/domain/enums"

	"github.com/amirex128/new_site_builder/src/internal/application/usecase"
	"github.com/amirex128/new_site_builder/src/internal/contract/service"

	"github.com/gin-gonic/gin"

	"github.com/amirex128/new_site_builder/src/internal/application/dto/user"
	"github.com/amirex128/new_site_builder/src/internal/contract"
	"github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/contract/repository"
	"github.com/amirex128/new_site_builder/src/internal/domain"
)

type UserUsecase struct {
	*usecase.BaseUsecase
	userRepo      repository.IUserRepository
	planRepo      repository.IPlanRepository
	addressRepo   repository.IAddressRepository
	paymentRepo   repository.IPaymentRepository
	identitySvc   service.IIdentityService
	authContext   func(c *gin.Context) service.IAuthService
	siteRepo      repository.ISiteRepository
	messageSvc    service.IMessageService
	unitPriceRepo repository.IUnitPriceRepository
}

func NewUserUsecase(c contract.IContainer) *UserUsecase {
	return &UserUsecase{
		BaseUsecase: &usecase.BaseUsecase{
			Logger: c.GetLogger(),
		},
		siteRepo:      c.GetSiteRepo(),
		userRepo:      c.GetUserRepo(),
		planRepo:      c.GetPlanRepo(),
		addressRepo:   c.GetAddressRepo(),
		paymentRepo:   c.GetPaymentRepo(),
		identitySvc:   c.GetIdentityService(),
		authContext:   c.GetAuthTransientService(),
		messageSvc:    c.GetMessageService(),
		unitPriceRepo: c.GetUnitPriceRepo(),
	}
}

func (u *UserUsecase) UpdateProfileUserCommand(params *user.UpdateProfileUserCommand) (*resp.Response, error) {
	userId, err := u.authContext(u.Ctx).GetUserID()
	if err != nil || userId == nil {
		return nil, resp.NewError(resp.Unauthorized, err.Error())
	}
	existingUser, err := u.userRepo.GetByID(*userId)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	if params.FirstName != nil {
		existingUser.FirstName = *params.FirstName
	}
	if params.LastName != nil {
		existingUser.LastName = *params.LastName
	}
	if params.Email != nil {
		existingUser.Email = *params.Email
	}
	if params.Password != nil {
		hashedPassword, salt := u.identitySvc.HashPassword(*params.Password)
		existingUser.Password = hashedPassword
		existingUser.Salt = salt
	}
	if params.NationalCode != nil {
		existingUser.NationalCode = *params.NationalCode
	}
	if params.Phone != nil {
		existingUser.Phone = *params.Phone
	}
	if params.AiTypeEnum != nil {
		existingUser.AiTypeEnum = *params.AiTypeEnum
	}
	if params.UseCustomEmailSmtp != nil {
		existingUser.UseCustomEmailSmtp = string(*params.UseCustomEmailSmtp)
	}
	if params.Smtp != nil {
		existingUser.SmtpHost = params.Smtp.Host
		existingUser.SmtpPort = &params.Smtp.Port
		existingUser.SmtpUsername = params.Smtp.Username
		existingUser.SmtpPassword = params.Smtp.Password
	}
	existingUser.UpdatedAt = time.Now()
	err = u.userRepo.Update(existingUser)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	if len(params.AddressIDs) > 0 {
		err = u.addressRepo.RemoveAllAddressesFromUser(*userId)
		if err != nil {
			return nil, resp.NewError(resp.Internal, err.Error())
		}
		for _, addressID := range params.AddressIDs {
			err = u.addressRepo.AddAddressToUser(addressID, *userId)
			if err != nil {
				return nil, resp.NewError(resp.Internal, err.Error())
			}
		}
	}
	return resp.NewResponse(resp.Updated, "User profile updated"), nil
}

func (u *UserUsecase) GetProfileUserQuery(params *user.GetProfileUserQuery) (*resp.Response, error) {
	userId, err := u.authContext(u.Ctx).GetUserID()
	if err != nil || userId == nil {
		return nil, resp.NewError(resp.Unauthorized, err.Error())
	}
	existingUser, err := u.userRepo.GetByID(*userId)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	addresses, err := u.addressRepo.GetAllByUserID(*userId, common.PaginationRequestDto{Page: 1, PageSize: 100})
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	return resp.NewResponseData(resp.Retrieved, resp.Data{"user": existingUser, "addresses": addresses}, "user profile successful retrieved"), nil
}

func (u *UserUsecase) RegisterUserCommand(params *user.RegisterUserCommand) (*resp.Response, error) {
	u.Logger.Info("RegisterUserCommand called", map[string]interface{}{
		"email": *params.Email,
	})

	_, err := u.userRepo.GetByEmail(*params.Email)
	if err == nil {
		return nil, resp.NewError(resp.BadRequest, fmt.Sprintf("user with email %s already exists", *params.Email))
	}

	hashedPassword, salt := u.identitySvc.HashPassword(*params.Password)

	newUser := domain.User{
		Email:     *params.Email,
		Password:  hashedPassword,
		Salt:      salt,
		IsActive:  enums.InactiveStatus,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	verificationCode := utils.GenerateVerificationCode()
	newUser.VerifyCode = &verificationCode

	err = u.userRepo.Create(newUser)
	if err != nil {
		u.Logger.Error("Error creating user", map[string]interface{}{
			"error": err.Error(),
			"email": *params.Email,
		})
		return nil, resp.NewError(resp.Internal, fmt.Sprintf("Error creating user: %s", err.Error()))
	}

	token := u.identitySvc.TokenForUser(newUser).Make()
	return resp.NewResponseData(
		resp.Created,
		resp.Data{
			"token": token,
		},
		"Registration successful. Please verify your account.",
	), nil
}

func (u *UserUsecase) LoginUserCommand(params *user.LoginUserCommand) (*resp.Response, error) {
	u.Logger.Info("LoginUserCommand called", map[string]interface{}{
		"email": *params.Email,
	})

	existingUser, err := u.userRepo.GetByEmail(*params.Email)
	if err != nil {
		u.Logger.Info("Login failed: user not found", map[string]interface{}{
			"email": *params.Email,
		})
		return nil, resp.NewError(resp.Unauthorized, "invalid email or password")
	}

	if existingUser.IsActive == enums.InactiveStatus {
		u.Logger.Info("Login failed: account not active", map[string]interface{}{
			"email":  *params.Email,
			"userId": existingUser.ID,
		})
		return nil, resp.NewError(resp.Unauthorized, "account is not active")
	}

	if !u.identitySvc.VerifyPassword(*params.Password, existingUser.Password, existingUser.Salt) {
		u.Logger.Info("Login failed: invalid password", map[string]interface{}{
			"email":  *params.Email,
			"userId": existingUser.ID,
		})
		return nil, resp.NewError(resp.Unauthorized, "invalid email or password")
	}

	token := u.identitySvc.TokenForUser(existingUser).Make()

	return resp.NewResponseData(
		resp.Created,
		resp.Data{
			"token": token,
		},
		"Login successful",
	), nil
}

func (u *UserUsecase) RequestVerifyAndForgetUserCommand(params *user.RequestVerifyAndForgetUserCommand) (*resp.Response, error) {
	var existingUser domain.User
	var err error

	// Get user by email or phone based on the verification type
	if params.Type != nil && (*params.Type == enums.VerifyEmailType || *params.Type == enums.ForgetPasswordEmailType) {
		if params.Email == nil {
			return nil, resp.NewError(resp.BadRequest, "email is required for email verification")
		}
		existingUser, err = u.userRepo.GetByEmail(*params.Email)
	} else if params.Type != nil && (*params.Type == enums.VerifyPhoneType || *params.Type == enums.ForgetPasswordPhoneType) {
		if params.Phone == nil {
			return nil, resp.NewError(resp.BadRequest, "phone is required for phone verification")
		}
		existingUser, err = u.userRepo.GetByPhone(*params.Phone)
	} else {
		return nil, resp.NewError(resp.BadRequest, "invalid verification type")
	}

	if err != nil {
		return nil, resp.NewError(resp.NotFound, "user not found")
	}

	// Generate verification code
	verificationCode := utils.GenerateVerificationCode()

	// Store verification code based on type
	if params.Type != nil && (*params.Type == enums.VerifyEmailType || *params.Type == enums.ForgetPasswordEmailType) {
		existingUser.VerifyCode = &verificationCode
	} else if params.Type != nil && (*params.Type == enums.VerifyPhoneType || *params.Type == enums.ForgetPasswordPhoneType) {
		existingUser.VerifyCode = &verificationCode
	}

	err = u.userRepo.Update(existingUser)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}

	// Send verification code via email or SMS
	if params.Type != nil && (*params.Type == enums.VerifyEmailType || *params.Type == enums.ForgetPasswordEmailType) && params.Email != nil {
		msg := struct {
			To      string
			Subject string
			Body    string
		}{
			To:      *params.Email,
			Subject: "Your Verification Code",
			Body:    fmt.Sprintf("Your verification code is: %s", verificationCode),
		}
		u.messageSvc.SendEmail(msg)
	} else if params.Type != nil && (*params.Type == enums.VerifyPhoneType || *params.Type == enums.ForgetPasswordPhoneType) && params.Phone != nil {
		msg := struct {
			To   string
			Body string
		}{
			To:   *params.Phone,
			Body: fmt.Sprintf("Your verification code is: %s", verificationCode),
		}
		u.messageSvc.SendSms(msg)
	}

	return resp.NewResponseData(
		resp.Success,
		resp.Data{
			"success": true,
			"message": "Verification code sent. Please check your email or phone.",
		},
		"Verification code sent. Please check your email or phone.",
	), nil
}

func (u *UserUsecase) VerifyUserQuery(params *user.VerifyUserQuery) (*resp.Response, error) {
	// Get user by email
	existingUser, err := u.userRepo.GetByEmail(*params.Email)
	if err != nil {
		return nil, resp.NewError(resp.NotFound, "user not found")
	}

	// Convert code to string for comparison
	codeStr := strconv.Itoa(*params.Code)

	var resetToken string

	// Check verification code based on type
	if params.Type != nil && *params.Type == enums.VerifyEmailType {
		if *existingUser.VerifyCode != codeStr {
			return nil, resp.NewError(resp.BadRequest, "invalid verification code")
		}
		existingUser.IsActive = enums.ActiveStatus // Activate the user
		existingUser.VerifyEmail = ""
	} else if params.Type != nil && *params.Type == enums.VerifyPhoneType {
		if *existingUser.VerifyCode != codeStr {
			return nil, resp.NewError(resp.BadRequest, "invalid verification code")
		}
		existingUser.IsActive = enums.ActiveStatus // Activate the user
		existingUser.VerifyPhone = ""
	} else if params.Type != nil && *params.Type == enums.ForgetPasswordEmailType {
		if *existingUser.VerifyCode != codeStr {
			return nil, resp.NewError(resp.BadRequest, "invalid verification code")
		}
		// Provide a token for password reset
		resetToken = u.identitySvc.AddClaim("reset_password", "1").AddClaim("user_id", strconv.FormatInt(existingUser.ID, 10)).Make()
		existingUser.VerifyEmail = ""
	} else if params.Type != nil && *params.Type == enums.ForgetPasswordPhoneType {
		if *existingUser.VerifyCode != codeStr {
			return nil, resp.NewError(resp.BadRequest, "invalid verification code")
		}
		// Provide a token for password reset
		resetToken = u.identitySvc.AddClaim("reset_password", "1").AddClaim("user_id", strconv.FormatInt(existingUser.ID, 10)).Make()
		existingUser.VerifyPhone = ""
	} else {
		return nil, resp.NewError(resp.BadRequest, "invalid verification type")
	}

	err = u.userRepo.Update(existingUser)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}

	respData := resp.Data{
		"success": true,
		"message": "Verification successful.",
	}
	if resetToken != "" {
		respData["reset_token"] = resetToken
	}

	return resp.NewResponseData(
		resp.Success,
		respData,
		"Verification successful.",
	), nil
}

func (u *UserUsecase) ChargeCreditRequestUserCommand(params *user.ChargeCreditRequestUserCommand) (*resp.Response, error) {
	// Get the current user ID
	userID, err := u.authContext(u.Ctx).GetUserID()
	if err != nil {
		return nil, resp.NewError(resp.Unauthorized, err.Error())
	}
	if userID == nil {
		return nil, resp.NewError(resp.Unauthorized, "user not authenticated")
	}

	// Calculate total amount
	var totalAmount int64 = 0
	orderData := make(map[string]string)

	for i, unitPrice := range params.UnitPrices {
		// Fetch actual unit price from database
		unitPriceObj, err := u.unitPriceRepo.GetByName(string(*unitPrice.UnitPriceName))
		if err != nil {
			return nil, resp.NewError(resp.BadRequest, fmt.Sprintf("unit price not found: %s", *unitPrice.UnitPriceName))
		}
		var itemPrice int64 = unitPriceObj.Price * int64(*unitPrice.UnitPriceCount)

		// For storage, multiply by days if provided
		if string(*unitPrice.UnitPriceName) == "storage_mb_credits" && unitPrice.UnitPriceDay != nil {
			itemPrice = itemPrice * int64(*unitPrice.UnitPriceDay)
		}

		totalAmount += itemPrice

		// Add unit price details to order data
		orderData[fmt.Sprintf("unitPrice_%d_name", i)] = string(*unitPrice.UnitPriceName)
		orderData[fmt.Sprintf("unitPrice_%d_count", i)] = strconv.Itoa(*unitPrice.UnitPriceCount)
		if unitPrice.UnitPriceDay != nil {
			orderData[fmt.Sprintf("unitPrice_%d_days", i)] = strconv.Itoa(*unitPrice.UnitPriceDay)
		}
	}

	// Create a new order in the payment system
	orderID := time.Now().UnixNano() // Use nanoseconds for more uniqueness

	// Additional order data
	orderData["userId"] = strconv.FormatInt(*userID, 10)
	orderData["totalAmount"] = strconv.FormatInt(totalAmount, 10)
	orderData["finalFrontReturnUrl"] = *params.FinalFrontReturnUrl

	// Request payment from gateway
	paymentUrl, err := u.paymentRepo.RequestPayment(totalAmount, orderID, *userID, string(*params.Gateway), orderData)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}

	return resp.NewResponseData(
		resp.Success,
		resp.Data{
			"paymentUrl": paymentUrl,
			"orderId":    orderID,
		},
		"Payment URL generated successfully.",
	), nil
}

func (u *UserUsecase) UpgradePlanRequestUserCommand(params *user.UpgradePlanRequestUserCommand) (*resp.Response, error) {
	// Get the current user ID
	userID, err := u.authContext(u.Ctx).GetUserID()
	if err != nil {
		return nil, resp.NewError(resp.Unauthorized, fmt.Sprintf("error : %s", err.Error()))
	}
	if userID == nil {
		return nil, resp.NewError(resp.Unauthorized, "user not authenticated")
	}

	// Get the plan
	plan, err := u.planRepo.GetByID(*params.PlanID)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}

	// Calculate final price (consider applying discounts here)
	var finalPrice int64 = plan.Price
	var discountAmount int64 = 0

	// Apply discount if available
	if plan.Discount != nil && *plan.Discount > 0 {
		if plan.DiscountType == string(enums.FixedDiscountType) {
			discountAmount = *plan.Discount
			finalPrice = plan.Price - discountAmount
		} else if plan.DiscountType == string(enums.PercentageDiscountType) {
			discountAmount = (plan.Price * (*plan.Discount)) / 100
			finalPrice = plan.Price - discountAmount
		}
	}

	// Ensure price doesn't go below zero
	if finalPrice < 0 {
		finalPrice = 0
	}

	// Create order data
	orderData := make(map[string]string)
	orderData["userId"] = strconv.FormatInt(*userID, 10)
	orderData["planId"] = strconv.FormatInt(*params.PlanID, 10)
	orderData["originalPrice"] = strconv.FormatInt(plan.Price, 10)
	orderData["discountAmount"] = strconv.FormatInt(discountAmount, 10)
	orderData["finalPrice"] = strconv.FormatInt(finalPrice, 10)
	orderData["finalFrontReturnUrl"] = *params.FinalFrontReturnUrl

	// Create a new order in the payment system
	orderID := time.Now().UnixNano() // Use nanoseconds for more uniqueness

	// Request payment from gateway
	paymentUrl, err := u.paymentRepo.RequestPayment(finalPrice, orderID, *userID, string(*params.Gateway), orderData)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}

	return resp.NewResponseData(
		resp.Success,
		resp.Data{
			"paymentUrl": paymentUrl,
			"orderId":    orderID,
			"plan":       plan,
		},
		"Plan upgrade payment URL generated successfully.",
	), nil
}

func (u *UserUsecase) AdminGetAllUserQuery(params *user.AdminGetAllUserQuery) (*resp.Response, error) {
	u.Logger.Info("AdminGetAllUserQuery called", map[string]interface{}{
		"page":     params.Page,
		"pageSize": params.PageSize,
	})

	// Check admin access
	isAdmin, err := u.authContext(u.Ctx).IsAdmin()
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	if !isAdmin {
		return nil, resp.NewError(resp.Unauthorized, "Only administrators can access this resource")
	}

	result, err := u.userRepo.GetAll(params.PaginationRequestDto)
	if err != nil {
		u.Logger.Error("Error getting all users", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, resp.NewError(resp.Internal, err.Error())
	}

	count := result.TotalCount

	// Calculate total pages
	totalPages := (count + int64(params.PageSize) - 1) / int64(params.PageSize)
	if totalPages < 1 {
		totalPages = 1
	}

	return resp.NewResponseData(
		resp.Success,
		resp.Data{
			"items":     result,
			"total":     count,
			"page":      params.Page,
			"pageSize":  params.PageSize,
			"totalPage": totalPages,
		},
		"All users retrieved successfully.",
	), nil
}
