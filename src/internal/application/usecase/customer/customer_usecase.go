package customerusecase

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/amirex128/new_site_builder/src/internal/domain/enums"

	"github.com/amirex128/new_site_builder/src/internal/application/usecase"
	"github.com/amirex128/new_site_builder/src/internal/application/utils"
	"github.com/amirex128/new_site_builder/src/internal/application/utils/resp"
	"github.com/amirex128/new_site_builder/src/internal/contract/service"

	"github.com/gin-gonic/gin"

	"github.com/amirex128/new_site_builder/src/internal/application/dto/customer"
	"github.com/amirex128/new_site_builder/src/internal/contract"
	"github.com/amirex128/new_site_builder/src/internal/contract/repository"
	"github.com/amirex128/new_site_builder/src/internal/domain"
	"gorm.io/gorm"
)

type CustomerUsecase struct {
	*usecase.BaseUsecase
	repo         repository.ICustomerRepository
	fileItemRepo repository.IFileItemRepository
	addressRepo  repository.IAddressRepository
	roleRepo     repository.IRoleRepository
	authContext  func(c *gin.Context) service.IAuthService
	identitySvc  service.IIdentityService
	messageSvc   service.IMessageService
}

func NewCustomerUsecase(c contract.IContainer) *CustomerUsecase {
	return &CustomerUsecase{
		BaseUsecase: &usecase.BaseUsecase{
			Logger: c.GetLogger(),
		},
		repo:         c.GetCustomerRepo(),
		fileItemRepo: c.GetFileItemRepo(),
		addressRepo:  c.GetAddressRepo(),
		roleRepo:     c.GetRoleRepo(),
		authContext:  c.GetAuthTransientService(),
		identitySvc:  c.GetIdentityService(),
		messageSvc:   c.GetMessageService(),
	}
}

// LoginCustomerCommand handles customer login
func (u *CustomerUsecase) LoginCustomerCommand(params *customer.LoginCustomerCommand) (*resp.Response, error) {
	u.Logger.Info("LoginCustomerCommand called", map[string]interface{}{
		"email": *params.Email,
	})

	// Get customer by email
	existingCustomer, err := u.repo.GetByEmail(*params.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, resp.NewError(resp.Unauthorized, "ایمیل یا رمز عبور اشتباه است")
		}
		return nil, resp.NewError(resp.Internal, err.Error())
	}

	// Verify password
	if !u.identitySvc.VerifyPassword(*params.Password, existingCustomer.Password, existingCustomer.Salt) {
		return nil, resp.NewError(resp.Unauthorized, "ایمیل یا رمز عبور اشتباه است")
	}

	// Get roles for customer
	roleNames := make([]string, 0, len(existingCustomer.Roles))
	for _, role := range existingCustomer.Roles {
		roleNames = append(roleNames, role.Name)
	}

	// If no roles, add default "Customer" role
	if len(roleNames) == 0 {
		roleNames = append(roleNames, "Customer")
	}

	// Generate token
	token := u.identitySvc.TokenForCustomer(existingCustomer).AddRoles(roleNames).Make()

	return resp.NewResponseData(
		resp.Created,
		resp.Data{
			"token": token,
		},
		"Login successful",
	), nil
}

// RegisterCustomerCommand handles customer registration
func (u *CustomerUsecase) RegisterCustomerCommand(params *customer.RegisterCustomerCommand) (*resp.Response, error) {
	u.Logger.Info("RegisterCustomerCommand called", map[string]interface{}{
		"email":  *params.Email,
		"siteId": *params.SiteID,
	})

	// Check if customer already exists with this email
	existingCustomer, err := u.repo.GetByEmail(*params.Email)
	if err == nil && existingCustomer.ID > 0 {
		return nil, resp.NewError(resp.BadRequest, "ایمیل قبلاً ثبت شده است")
	}

	// Hash password with salt
	hashedPassword, salt := u.identitySvc.HashPassword(*params.Password)

	// Create new customer
	newCustomer := domain.Customer{
		Email:     *params.Email,
		SiteID:    *params.SiteID,
		Password:  hashedPassword,
		Salt:      salt,
		IsActive:  enums.ActiveStatus,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		IsDeleted: false,
	}

	// Generate verification code
	verificationCode := utils.GenerateVerificationCode()
	newCustomer.VerifyCode = &verificationCode

	// Save customer
	err = u.repo.Create(newCustomer)
	if err != nil {
		u.Logger.Error("Error creating customer", map[string]interface{}{
			"error": err.Error(),
			"email": *params.Email,
		})
		return nil, resp.NewError(resp.Internal, "خطا در ثبت نام مشتری")
	}

	// Retrieve the created customer to get the ID
	createdCustomer, err := u.repo.GetByEmail(*params.Email)
	if err != nil {
		return nil, resp.NewError(resp.Internal, "خطا در بازیابی اطلاعات مشتری")
	}

	// Add default customer role
	defaultRole, err := u.roleRepo.GetByName("Customer")
	if err == nil {
		err = u.roleRepo.AddRoleToCustomer(defaultRole.ID, createdCustomer.ID)
		if err != nil {
			u.Logger.Error("Error assigning default role to customer", map[string]interface{}{
				"error":      err.Error(),
				"customerId": createdCustomer.ID,
			})
			// Continue despite role assignment error
		}
	}

	// Generate token
	token := u.identitySvc.TokenForCustomer(createdCustomer).AddRoles([]string{"Customer"}).Make()

	// Send verification email
	msg := struct {
		To      string
		Subject string
		Body    string
	}{
		To:      *params.Email,
		Subject: "Verify Your Account",
		Body:    fmt.Sprintf("Your verification code is: %s", verificationCode),
	}
	u.messageSvc.SendEmail(msg)

	return resp.NewResponseData(
		resp.Created,
		resp.Data{
			"token": token,
		},
		"Registration successful. Please verify your account.",
	), nil
}

// RequestVerifyAndForgetCustomerCommand handles verification and password reset requests
func (u *CustomerUsecase) RequestVerifyAndForgetCustomerCommand(params *customer.RequestVerifyAndForgetCustomerCommand) (*resp.Response, error) {
	u.Logger.Info("RequestVerifyAndForgetCustomerCommand called", map[string]interface{}{
		"email": *params.Email,
		"phone": *params.Phone,
		"type":  *params.Type,
	})

	// Handle phone verification/forget
	if *params.Type == enums.VerifyPhoneType || *params.Type == enums.ForgetPasswordPhoneType {
		if params.Phone == nil {
			return nil, resp.NewError(resp.BadRequest, "phone is required for phone verification")
		}

		existingCustomer, err := u.repo.GetByPhone(*params.Phone)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, resp.NewError(resp.NotFound, "شماره تلفن یافت نشد")
			}
			return nil, resp.NewError(resp.Internal, err.Error())
		}

		// Generate verification code
		verifyCode := utils.GenerateVerificationCode()
		verifyCodeStr := verifyCode
		expireAt := time.Now().Add(15 * time.Minute)

		// Update customer with verification code
		existingCustomer.VerifyCode = &verifyCodeStr
		existingCustomer.ExpireVerifyCodeAt = &expireAt
		existingCustomer.UpdatedAt = time.Now()

		err = u.repo.Update(existingCustomer)
		if err != nil {
			return nil, resp.NewError(resp.Internal, "خطا در بروزرسانی اطلاعات مشتری")
		}

		// Send verification code via SMS
		msg := struct {
			To   string
			Body string
		}{
			To:   *params.Phone,
			Body: fmt.Sprintf("Your verification code is: %s", verifyCode),
		}
		u.messageSvc.SendSms(msg)

		return resp.NewResponseData(
			resp.Success,
			resp.Data{
				"success": true,
				"message": "Verification code sent. Please check your phone.",
			},
			"Verification code sent. Please check your phone.",
		), nil
	}

	// Handle email verification/forget
	if *params.Type == enums.VerifyEmailType || *params.Type == enums.ForgetPasswordEmailType {
		if params.Email == nil {
			return nil, resp.NewError(resp.BadRequest, "email is required for email verification")
		}

		existingCustomer, err := u.repo.GetByEmail(*params.Email)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, resp.NewError(resp.NotFound, "ایمیل یافت نشد")
			}
			return nil, resp.NewError(resp.Internal, err.Error())
		}

		// Generate verification code
		verifyCode := utils.GenerateVerificationCode()
		verifyCodeStr := verifyCode
		expireAt := time.Now().Add(24 * time.Hour)

		// Update customer with verification code
		existingCustomer.VerifyCode = &verifyCodeStr
		existingCustomer.ExpireVerifyCodeAt = &expireAt
		existingCustomer.UpdatedAt = time.Now()

		err = u.repo.Update(existingCustomer)
		if err != nil {
			return nil, resp.NewError(resp.Internal, "خطا در بروزرسانی اطلاعات مشتری")
		}

		// Generate reset link
		gatewayURL := os.Getenv("HTTP_GATEWAY_URL")
		if gatewayURL == "" {
			gatewayURL = "http://localhost"
		}

		resetLink := generatePasswordResetLink(gatewayURL, existingCustomer.ID, verifyCodeStr, string(*params.Type))

		// Send verification code via email
		msg := struct {
			To      string
			Subject string
			Body    string
		}{
			To:      *params.Email,
			Subject: "Your Verification Code",
			Body:    fmt.Sprintf("Your verification code is: %s\nOr use this link: %s", verifyCode, resetLink),
		}
		u.messageSvc.SendEmail(msg)

		return resp.NewResponseData(
			resp.Success,
			resp.Data{
				"success": true,
				"message": "Verification code sent. Please check your email.",
			},
			"Verification code sent. Please check your email.",
		), nil
	}

	return nil, resp.NewError(resp.BadRequest, "نوع تأیید نامعتبر است")
}

// UpdateProfileCustomerCommand handles updating customer profile
func (u *CustomerUsecase) UpdateProfileCustomerCommand(params *customer.UpdateProfileCustomerCommand) (*resp.Response, error) {
	u.Logger.Info("UpdateProfileCustomerCommand called", map[string]interface{}{
		"firstName": params.FirstName,
		"lastName":  params.LastName,
	})

	// Get customer ID from auth context
	customerID, err := u.authContext(u.Ctx).GetCustomerID()
	if err != nil {
		return nil, resp.NewError(resp.Unauthorized, "خطا در احراز هویت مشتری")
	}

	// Get existing customer
	existingCustomer, err := u.repo.GetByID(customerID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, resp.NewError(resp.NotFound, "مشتری یافت نشد")
		}
		return nil, resp.NewError(resp.Internal, err.Error())
	}

	// Update fields if provided
	if params.FirstName != nil {
		existingCustomer.FirstName = *params.FirstName
	}

	if params.LastName != nil {
		existingCustomer.LastName = *params.LastName
	}

	if params.Email != nil {
		existingCustomer.Email = *params.Email
	}

	if params.Password != nil {
		hashedPassword, salt := u.identitySvc.HashPassword(*params.Password)
		existingCustomer.Password = hashedPassword
		existingCustomer.Salt = salt
	}

	if params.NationalCode != nil {
		existingCustomer.NationalCode = *params.NationalCode
	}

	if params.Phone != nil {
		existingCustomer.Phone = *params.Phone
	}

	// Update timestamp
	existingCustomer.UpdatedAt = time.Now()

	// Save customer
	err = u.repo.Update(existingCustomer)
	if err != nil {
		return nil, resp.NewError(resp.Internal, "خطا در بروزرسانی پروفایل مشتری")
	}

	// Handle address IDs if provided
	if len(params.AddressIDs) > 0 {
		// First, remove all existing addresses
		err = u.addressRepo.RemoveAllAddressesFromCustomer(customerID)
		if err != nil {
			return nil, resp.NewError(resp.Internal, "خطا در بروزرسانی آدرس‌های مشتری")
		}

		// Then add the new addresses
		for _, addressID := range params.AddressIDs {
			err = u.addressRepo.AddAddressToCustomer(addressID, customerID)
			if err != nil {
				return nil, resp.NewError(resp.Internal, "خطا در افزودن آدرس به مشتری")
			}
		}
	}

	// Return updated customer with enhanced response
	return resp.NewResponseData(
		resp.Updated,
		resp.Data{
			"customer": existingCustomer,
		},
		"Customer profile updated successfully.",
	), nil
}

// VerifyCustomerQuery handles customer verification
func (u *CustomerUsecase) VerifyCustomerQuery(params *customer.VerifyCustomerQuery) (*resp.Response, error) {
	u.Logger.Info("VerifyCustomerQuery called", map[string]interface{}{
		"email": *params.Email,
		"code":  *params.Code,
		"type":  *params.Type,
	})

	// Get customer by email
	existingCustomer, err := u.repo.GetByEmail(*params.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, resp.NewError(resp.NotFound, "ایمیل یافت نشد")
		}
		return nil, resp.NewError(resp.Internal, err.Error())
	}

	// Check verification code
	if existingCustomer.VerifyCode == nil || *existingCustomer.VerifyCode != *params.Code {
		return nil, resp.NewError(resp.BadRequest, "کد تایید نامعتبر است")
	}

	// Check if code is expired
	if existingCustomer.ExpireVerifyCodeAt == nil || time.Now().After(*existingCustomer.ExpireVerifyCodeAt) {
		return nil, resp.NewError(resp.BadRequest, "کد تایید منقضی شده است")
	}

	// Clear verification code
	existingCustomer.VerifyCode = nil
	existingCustomer.ExpireVerifyCodeAt = nil
	existingCustomer.UpdatedAt = time.Now()

	var resetToken string

	// Handle different verification types
	switch *params.Type {
	case enums.VerifyEmailType:
		existingCustomer.VerifyEmail = enums.ActiveStatus
		err = u.repo.Update(existingCustomer)
		if err != nil {
			return nil, resp.NewError(resp.Internal, "خطا در بروزرسانی اطلاعات مشتری")
		}
		return resp.NewResponseData(
			resp.Success,
			resp.Data{
				"success": true,
				"message": "Email verification successful.",
			},
			"Email verification successful.",
		), nil

	case enums.VerifyPhoneType:
		existingCustomer.VerifyPhone = enums.ActiveStatus
		err = u.repo.Update(existingCustomer)
		if err != nil {
			return nil, resp.NewError(resp.Internal, "خطا در بروزرسانی اطلاعات مشتری")
		}
		return resp.NewResponseData(
			resp.Success,
			resp.Data{
				"success": true,
				"message": "Phone verification successful.",
			},
			"Phone verification successful.",
		), nil

	case enums.ForgetPasswordEmailType:
		// Generate token for password reset with claims
		resetToken = u.identitySvc.AddClaim("reset_password", "1").AddClaim("customer_id", strconv.FormatInt(existingCustomer.ID, 10)).Make()
		err = u.repo.Update(existingCustomer)
		if err != nil {
			return nil, resp.NewError(resp.Internal, "خطا در بروزرسانی اطلاعات مشتری")
		}

		// Get reset URL from environment or use default
		resetURL := os.Getenv("RESET_PASSWORD_URL")
		if resetURL == "" {
			resetURL = "https://example.com/reset-password"
		}

		// Add token to URL
		queryParams := url.Values{}
		queryParams.Add("token", resetToken)
		resetURL = resetURL + "?" + queryParams.Encode()

		return resp.NewResponseData(
			resp.Success,
			resp.Data{
				"success":     true,
				"message":     "Password reset verification successful.",
				"reset_token": resetToken,
				"reset_url":   resetURL,
			},
			"Password reset verification successful.",
		), nil

	case enums.ForgetPasswordPhoneType:
		// Generate token for password reset with claims
		resetToken = u.identitySvc.AddClaim("reset_password", "1").AddClaim("customer_id", strconv.FormatInt(existingCustomer.ID, 10)).Make()
		err = u.repo.Update(existingCustomer)
		if err != nil {
			return nil, resp.NewError(resp.Internal, "خطا در بروزرسانی اطلاعات مشتری")
		}
		return resp.NewResponseData(
			resp.Success,
			resp.Data{
				"success":     true,
				"message":     "Password reset verification successful.",
				"reset_token": resetToken,
			},
			"Password reset verification successful.",
		), nil

	default:
		return nil, resp.NewError(resp.BadRequest, "نوع تأیید نامعتبر است")
	}
}

// GetProfileCustomerQuery handles getting customer profile
func (u *CustomerUsecase) GetProfileCustomerQuery(params *customer.GetProfileCustomerQuery) (*resp.Response, error) {
	u.Logger.Info("GetProfileCustomerQuery called", nil)

	// Get customer ID from auth context
	customerID, err := u.authContext(u.Ctx).GetCustomerID()
	if err != nil {
		return nil, resp.NewError(resp.Unauthorized, "خطا در احراز هویت مشتری")
	}

	// Get customer details
	existingCustomer, err := u.repo.GetByID(customerID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, resp.NewError(resp.NotFound, "مشتری یافت نشد")
		}
		return nil, resp.NewError(resp.Internal, err.Error())
	}

	// Get customer addresses
	addresses, err := u.addressRepo.GetAllByCustomerID(customerID)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}

	return resp.NewResponseData(
		resp.Retrieved,
		resp.Data{
			"customer": existingCustomer,
			"address":  addresses,
		},
		"Customer profile retrieved successfully.",
	), nil
}

// AdminGetAllCustomerQuery handles admin getting all customers
func (u *CustomerUsecase) AdminGetAllCustomerQuery(params *customer.AdminGetAllCustomerQuery) (*resp.Response, error) {
	u.Logger.Info("AdminGetAllCustomerQuery called", map[string]interface{}{
		"page":     params.Page,
		"pageSize": params.PageSize,
	})

	// Check admin access
	isAdmin, err := u.authContext(u.Ctx).IsAdmin()
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	if !isAdmin {
		return nil, resp.NewError(resp.Unauthorized, "فقط مدیران سیستم مجاز به دسترسی به این بخش هستند")
	}

	// Get all customers with pagination
	customers, count, err := u.repo.GetAll(params.PaginationRequestDto)
	if err != nil {
		return nil, resp.NewError(resp.Internal, "خطا در دریافت لیست مشتریان")
	}

	// Calculate total pages
	totalPages := (count + int64(params.PageSize) - 1) / int64(params.PageSize)
	if totalPages < 1 {
		totalPages = 1
	}

	return resp.NewResponseData(
		resp.Success,
		resp.Data{
			"items":     customers,
			"total":     count,
			"page":      params.Page,
			"pageSize":  params.PageSize,
			"totalPage": totalPages,
		},
		"All customers retrieved successfully.",
	), nil
}

// Helper function to generate a password reset link
func generatePasswordResetLink(baseURL string, customerID int64, code string, verifyType string) string {
	// In a real implementation, we might encrypt or sign the parameters
	return fmt.Sprintf("%s/reset-password?id=%d&code=%s&type=%s",
		baseURL, customerID, code, verifyType)
}
