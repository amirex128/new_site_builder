package userusecase

import (
	"fmt"
	"github.com/amirex128/new_site_builder/src/internal/application/utils/nerror"
	"github.com/amirex128/new_site_builder/src/internal/domain/enums"
	"strconv"
	"time"

	"github.com/amirex128/new_site_builder/src/internal/application/usecase"
	"github.com/amirex128/new_site_builder/src/internal/contract/service"

	"github.com/gin-gonic/gin"

	"github.com/amirex128/new_site_builder/src/internal/application/dto/user"
	"github.com/amirex128/new_site_builder/src/internal/contract"
	"github.com/amirex128/new_site_builder/src/internal/contract/repository"
	"github.com/amirex128/new_site_builder/src/internal/domain"
)

type UserUsecase struct {
	*usecase.BaseUsecase
	userRepo    repository.IUserRepository
	planRepo    repository.IPlanRepository
	addressRepo repository.IAddressRepository
	paymentRepo repository.IPaymentRepository
	identitySvc service.IIdentityService
	authContext func(c *gin.Context) service.IAuthService
	siteRepo    repository.ISiteRepository
}

func NewUserUsecase(c contract.IContainer) *UserUsecase {
	return &UserUsecase{
		BaseUsecase: &usecase.BaseUsecase{
			Logger: c.GetLogger(),
		},
		siteRepo:    c.GetSiteRepo(),
		userRepo:    c.GetUserRepo(),
		planRepo:    c.GetPlanRepo(),
		addressRepo: c.GetAddressRepo(),
		paymentRepo: c.GetPaymentRepo(),
		identitySvc: c.GetIdentityService(),
		authContext: c.GetAuthTransientService(),
	}
}

func (u *UserUsecase) UpdateProfileUserCommand(params *user.UpdateProfileUserCommand) (any, error) {
	// Implementation for updating a user's profile
	fmt.Println(params)

	// In a real implementation, get the user ID from the auth context
	userID := int64(1)

	existingUser, err := u.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
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
		// In a real implementation, hash the password before storing
		existingUser.Password = *params.Password
	}

	if params.NationalCode != nil {
		existingUser.NationalCode = *params.NationalCode
	}

	if params.Phone != nil {
		existingUser.Phone = *params.Phone
	}

	if params.AiTypeEnum != nil {
		existingUser.AiTypeEnum = string(*params.AiTypeEnum)
	}

	if params.UseCustomEmailSmtp != nil {
		existingUser.UseCustomEmailSmtp = string(*params.UseCustomEmailSmtp)
	}

	if params.Smtp != nil {
		// In a real implementation, encrypt sensitive information like SMTP password
		existingUser.SmtpHost = params.Smtp.Host
		existingUser.SmtpPort = &params.Smtp.Port
		existingUser.SmtpUsername = params.Smtp.Username
		existingUser.SmtpPassword = params.Smtp.Password
	}

	existingUser.UpdatedAt = time.Now()

	err = u.userRepo.Update(existingUser)
	if err != nil {
		return nil, err
	}

	// Handle address IDs if provided
	if len(params.AddressIDs) > 0 {
		// First, remove all existing addresses
		err = u.addressRepo.RemoveAllAddressesFromUser(userID)
		if err != nil {
			return nil, err
		}

		// Then add the new addresses
		for _, addressID := range params.AddressIDs {
			err = u.addressRepo.AddAddressToUser(addressID, userID)
			if err != nil {
				// Log error but continue
				u.Logger.Errorf("Failed to assign address %d to user %d: %v", addressID, userID, err)
			}
		}
	}

	return existingUser, nil
}

func (u *UserUsecase) GetProfileUserQuery(params *user.GetProfileUserQuery) (any, error) {
	// Implementation to get a user's profile
	fmt.Println(params)

	// In a real implementation, get the user ID from the auth context
	userID := int64(1)

	result, err := u.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	// Get user addresses
	addresses, err := u.addressRepo.GetAllByUserID(userID)
	if err != nil {
		u.Logger.Errorf("Failed to get addresses for user %d: %v", userID, err)
	}

	return map[string]interface{}{
		"user":      result,
		"addresses": addresses,
	}, nil
}

func (u *UserUsecase) RegisterUserCommand(params *user.RegisterUserCommand) (any, error) {
	_, err := u.userRepo.GetByEmail(*params.Email)
	if err == nil {
		return nil, nerror.NewNError(nerror.BadRequest, "user with email %s already exists", *params.Email)
	}

	// Generate salt and hash password
	hashedPassword, salt := u.identitySvc.HashPassword(*params.Password)

	// Create new user
	newUser := domain.User{
		Email:     *params.Email,
		Password:  hashedPassword,
		Salt:      salt,
		IsActive:  "0", // Not active until verification
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		IsDeleted: false,
	}

	// Generate verification code
	verificationCode := "123456" // In a real implementation, generate a random code
	newUser.VerifyEmail = verificationCode

	err = u.userRepo.Create(newUser)
	if err != nil {
		return nil, err
	}

	// TODO: Send verification email

	return map[string]interface{}{
		"success": true,
		"message": "Registration successful. Please verify your account.",
	}, nil
}

func (u *UserUsecase) LoginUserCommand(params *user.LoginUserCommand) (any, error) {
	// Get user by email
	existingUser, err := u.userRepo.GetByEmail(*params.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	// Check if user is active
	if existingUser.IsActive != "1" {
		return nil, fmt.Errorf("account is not active")
	}

	// Verify password
	if !u.identitySvc.VerifyPassword(*params.Password, existingUser.Password, existingUser.Salt) {
		return nil, fmt.Errorf("invalid email or password")
	}

	// Generate JWT token
	token := u.identitySvc.TokenForUser(existingUser).Make()

	return map[string]interface{}{
		"token": token,
		"user":  existingUser,
	}, nil
}

func (u *UserUsecase) RequestVerifyAndForgetUserCommand(params *user.RequestVerifyAndForgetUserCommand) (any, error) {
	var existingUser domain.User
	var err error

	// Get user by email or phone based on the verification type
	if params.Type != nil && (*params.Type == enums.VerifyEmailType || *params.Type == enums.ForgetPasswordEmailType) {
		if params.Email == nil {
			return nil, fmt.Errorf("email is required for email verification")
		}
		existingUser, err = u.userRepo.GetByEmail(*params.Email)
	} else if params.Type != nil && (*params.Type == enums.VerifyPhoneType || *params.Type == enums.ForgetPasswordPhoneType) {
		if params.Phone == nil {
			return nil, fmt.Errorf("phone is required for phone verification")
		}
		existingUser, err = u.userRepo.GetByPhone(*params.Phone)
	} else {
		return nil, fmt.Errorf("invalid verification type")
	}

	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// Generate verification code
	verificationCode := "123456" // Example code, in a real implementation generate a random code

	// Store verification code based on type
	if params.Type != nil && (*params.Type == enums.VerifyEmailType || *params.Type == enums.ForgetPasswordEmailType) {
		existingUser.VerifyEmail = verificationCode
	} else if params.Type != nil && (*params.Type == enums.VerifyPhoneType || *params.Type == enums.ForgetPasswordPhoneType) {
		existingUser.VerifyPhone = verificationCode
	}

	err = u.userRepo.Update(existingUser)
	if err != nil {
		return nil, err
	}

	// TODO: Send verification code via email or SMS

	return map[string]interface{}{
		"success": true,
		"message": "Verification code sent. Please check your email or phone.",
	}, nil
}

func (u *UserUsecase) VerifyUserQuery(params *user.VerifyUserQuery) (any, error) {
	// Get user by email
	existingUser, err := u.userRepo.GetByEmail(*params.Email)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// Convert code to string for comparison
	codeStr := strconv.Itoa(*params.Code)

	// Check verification code based on type
	if params.Type != nil && *params.Type == enums.VerifyEmailType {
		if existingUser.VerifyEmail != codeStr {
			return nil, fmt.Errorf("invalid verification code")
		}
		existingUser.IsActive = "1" // Activate the user
		existingUser.VerifyEmail = ""
	} else if params.Type != nil && *params.Type == enums.VerifyPhoneType {
		if existingUser.VerifyPhone != codeStr {
			return nil, fmt.Errorf("invalid verification code")
		}
		existingUser.IsActive = "1" // Activate the user
		existingUser.VerifyPhone = ""
	} else if params.Type != nil && *params.Type == enums.ForgetPasswordEmailType {
		if existingUser.VerifyEmail != codeStr {
			return nil, fmt.Errorf("invalid verification code")
		}
		// In a real implementation, provide a token for password reset instead
		existingUser.VerifyEmail = ""
	} else if params.Type != nil && *params.Type == enums.ForgetPasswordPhoneType {
		if existingUser.VerifyPhone != codeStr {
			return nil, fmt.Errorf("invalid verification code")
		}
		// In a real implementation, provide a token for password reset instead
		existingUser.VerifyPhone = ""
	} else {
		return nil, fmt.Errorf("invalid verification type")
	}

	err = u.userRepo.Update(existingUser)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"success": true,
		"message": "Verification successful.",
	}, nil
}

func (u *UserUsecase) ChargeCreditRequestUserCommand(params *user.ChargeCreditRequestUserCommand) (any, error) {
	// Get the current user ID
	userID, err := u.authContext(u.Ctx).GetUserID()
	if err != nil {
		return nil, err
	}
	if userID == 0 {
		return nil, fmt.Errorf("user not authenticated")
	}

	// Calculate total amount
	var totalAmount int64 = 0
	orderData := make(map[string]string)

	for i, unitPrice := range params.UnitPrices {
		// In a real implementation, fetch actual unit prices from database
		// For now, using dummy values
		var itemPrice int64 = 1000 * int64(*unitPrice.UnitPriceCount)

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
	orderID := time.Now().Unix() // Dummy order ID

	// Additional order data
	orderData["userId"] = strconv.FormatInt(userID, 10)
	orderData["totalAmount"] = strconv.FormatInt(totalAmount, 10)
	orderData["finalFrontReturnUrl"] = *params.FinalFrontReturnUrl

	// Request payment from gateway
	paymentUrl, err := u.paymentRepo.RequestPayment(totalAmount, orderID, userID, string(*params.Gateway), orderData)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"paymentUrl": paymentUrl,
		"orderId":    orderID,
	}, nil
}

func (u *UserUsecase) UpgradePlanRequestUserCommand(params *user.UpgradePlanRequestUserCommand) (any, error) {
	// Get the current user ID
	userID, err := u.authContext(u.Ctx).GetUserID()
	if err != nil {
		return nil, err
	}
	if userID == 0 {
		return nil, fmt.Errorf("user not authenticated")
	}

	// Get the plan
	plan, err := u.planRepo.GetByID(*params.PlanID)
	if err != nil {
		return nil, err
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
	orderData["userId"] = strconv.FormatInt(userID, 10)
	orderData["planId"] = strconv.FormatInt(*params.PlanID, 10)
	orderData["originalPrice"] = strconv.FormatInt(plan.Price, 10)
	orderData["discountAmount"] = strconv.FormatInt(discountAmount, 10)
	orderData["finalPrice"] = strconv.FormatInt(finalPrice, 10)
	orderData["finalFrontReturnUrl"] = *params.FinalFrontReturnUrl

	// Create a new order in the payment system
	orderID := time.Now().Unix() // Dummy order ID

	// Request payment from gateway
	paymentUrl, err := u.paymentRepo.RequestPayment(finalPrice, orderID, userID, string(*params.Gateway), orderData)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"paymentUrl": paymentUrl,
		"orderId":    orderID,
		"plan":       plan,
	}, nil
}

func (u *UserUsecase) AdminGetAllUserQuery(params *user.AdminGetAllUserQuery) (any, error) {
	// Implementation for admin to get all users
	fmt.Println(params)

	result, count, err := u.userRepo.GetAll(params.PaginationRequestDto)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"items": result,
		"total": count,
	}, nil
}
