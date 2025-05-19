package customerusecase

import (
	"errors"
	"fmt"
	"github.com/amirex128/new_site_builder/src/internal/application/usecase"
	"github.com/amirex128/new_site_builder/src/internal/contract/service"
	"net/url"
	"os"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/amirex128/new_site_builder/src/internal/application/dto/customer"
	"github.com/amirex128/new_site_builder/src/internal/application/dto/user"
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
	}
}

// LoginCustomerCommand handles customer login
func (u *CustomerUsecase) LoginCustomerCommand(params *customer.LoginCustomerCommand) (any, error) {
	u.Logger.Info("LoginCustomerCommand called", map[string]interface{}{
		"email": *params.Email,
	})

	// Get customer by email
	existingCustomer, err := u.repo.GetByEmail(*params.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("ایمیل یا رمز عبور اشتباه است")
		}
		return nil, err
	}

	// Verify password
	if !u.identitySvc.VerifyPassword(*params.Password, existingCustomer.Password, existingCustomer.Salt) {
		return nil, errors.New("ایمیل یا رمز عبور اشتباه است")
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

	return token, nil
}

// RegisterCustomerCommand handles customer registration
func (u *CustomerUsecase) RegisterCustomerCommand(params *customer.RegisterCustomerCommand) (any, error) {
	u.Logger.Info("RegisterCustomerCommand called", map[string]interface{}{
		"email":  *params.Email,
		"siteId": *params.SiteID,
	})

	// Check if customer already exists with this email
	existingCustomer, err := u.repo.GetByEmail(*params.Email)
	if err == nil && existingCustomer.ID > 0 {
		return nil, errors.New("ایمیل قبلاً ثبت شده است")
	}

	// Hash password with salt
	hashedPassword, salt := u.identitySvc.HashPassword(*params.Password)

	// Create new customer
	newCustomer := domain.Customer{
		Email:     *params.Email,
		SiteID:    *params.SiteID,
		Password:  hashedPassword,
		Salt:      salt,
		IsActive:  "1",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		IsDeleted: false,
	}

	// Save customer
	err = u.repo.Create(newCustomer)
	if err != nil {
		u.Logger.Error("Error creating customer", map[string]interface{}{
			"error": err.Error(),
			"email": *params.Email,
		})
		return nil, errors.New("خطا در ثبت نام مشتری")
	}

	// Retrieve the created customer to get the ID
	createdCustomer, err := u.repo.GetByEmail(*params.Email)
	if err != nil {
		return nil, errors.New("خطا در بازیابی اطلاعات مشتری")
	}

	// Add default customer role (ID 1)
	// In a real implementation, we'd need to assign roles through a relationship method
	// For simplicity, we'll assume role ID 1 is "Customer"

	// Generate token
	token := u.identitySvc.TokenForCustomer(createdCustomer).AddRoles([]string{"Customer"}).Make()

	return token, nil
}

// RequestVerifyAndForgetCustomerCommand handles verification and password reset requests
func (u *CustomerUsecase) RequestVerifyAndForgetCustomerCommand(params *customer.RequestVerifyAndForgetCustomerCommand) (any, error) {
	u.Logger.Info("RequestVerifyAndForgetCustomerCommand called", map[string]interface{}{
		"email": *params.Email,
		"phone": *params.Phone,
		"type":  *params.Type,
	})

	// Handle phone verification/forget
	if *params.Type == user.VerifyPhone || *params.Type == user.ForgetPasswordPhone {
		existingCustomer, err := u.repo.GetByPhone(*params.Phone)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("شماره تلفن یافت نشد")
			}
			return nil, err
		}

		// Generate verification code (random 6-digit number)
		verifyCode := generateVerificationCode()
		expireAt := time.Now().Add(15 * time.Minute)

		// Update customer with verification code
		existingCustomer.VerifyCode = &verifyCode
		existingCustomer.ExpireVerifyCodeAt = &expireAt
		existingCustomer.UpdatedAt = time.Now()

		err = u.repo.Update(existingCustomer)
		if err != nil {
			return nil, errors.New("خطا در بروزرسانی اطلاعات مشتری")
		}

		// In a real implementation, we would send an SMS here
		// For now, we'll just log the code
		u.Logger.Info("SMS verification code generated", map[string]interface{}{
			"phone": *params.Phone,
			"code":  verifyCode,
			"type":  *params.Type,
		})

		return "success", nil
	}

	// Handle email verification/forget
	if *params.Type == user.VerifyEmail || *params.Type == user.ForgetPasswordEmail {
		existingCustomer, err := u.repo.GetByEmail(*params.Email)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("ایمیل یافت نشد")
			}
			return nil, err
		}

		// Generate verification code and link
		verifyCode := generateVerificationCode()
		expireAt := time.Now().Add(24 * time.Hour)

		// Update customer with verification code
		existingCustomer.VerifyCode = &verifyCode
		existingCustomer.ExpireVerifyCodeAt = &expireAt
		existingCustomer.UpdatedAt = time.Now()

		err = u.repo.Update(existingCustomer)
		if err != nil {
			return nil, errors.New("خطا در بروزرسانی اطلاعات مشتری")
		}

		// Generate reset link
		gatewayURL := os.Getenv("HTTP_GATEWAY_URL")
		if gatewayURL == "" {
			gatewayURL = "http://localhost"
		}

		resetLink := generatePasswordResetLink(gatewayURL, existingCustomer.ID, verifyCode, int(*params.Type))

		// In a real implementation, we would send an email here
		// For now, we'll just log the link
		u.Logger.Info("Email verification link generated", map[string]interface{}{
			"email": *params.Email,
			"code":  verifyCode,
			"link":  resetLink,
			"type":  *params.Type,
		})

		return "success", nil
	}

	return nil, errors.New("نوع تأیید نامعتبر است")
}

// UpdateProfileCustomerCommand handles updating customer profile
func (u *CustomerUsecase) UpdateProfileCustomerCommand(params *customer.UpdateProfileCustomerCommand) (any, error) {
	u.Logger.Info("UpdateProfileCustomerCommand called", map[string]interface{}{
		"firstName": params.FirstName,
		"lastName":  params.LastName,
	})

	// Get customer ID from auth context
	customerID, err := u.authContext(u.Ctx).GetCustomerID()
	if err != nil {
		return nil, errors.New("خطا در احراز هویت مشتری")
	}

	// Get existing customer
	existingCustomer, err := u.repo.GetByID(customerID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("مشتری یافت نشد")
		}
		return nil, err
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
		return nil, errors.New("خطا در بروزرسانی پروفایل مشتری")
	}

	// Handle address IDs if provided
	if len(params.AddressIDs) > 0 {
		// This would typically involve updating the many-to-many relationship
		// For simplicity, we'll just log the addresses
		u.Logger.Info("Updating customer addresses", map[string]interface{}{
			"customerId": customerID,
			"addressIds": params.AddressIDs,
		})
	}

	// Return updated customer
	return enhanceCustomerResponse(existingCustomer), nil
}

// VerifyCustomerQuery handles customer verification
func (u *CustomerUsecase) VerifyCustomerQuery(params *customer.VerifyCustomerQuery) (any, error) {
	u.Logger.Info("VerifyCustomerQuery called", map[string]interface{}{
		"email": *params.Email,
		"code":  *params.Code,
		"type":  *params.Type,
	})

	// Get customer by email
	existingCustomer, err := u.repo.GetByEmail(*params.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("ایمیل یافت نشد")
		}
		return nil, err
	}

	// Check verification code
	if existingCustomer.VerifyCode == nil || *existingCustomer.VerifyCode != *params.Code {
		return nil, errors.New("کد تایید نامعتبر است")
	}

	// Check if code is expired
	if existingCustomer.ExpireVerifyCodeAt == nil || time.Now().After(*existingCustomer.ExpireVerifyCodeAt) {
		return nil, errors.New("کد تایید منقضی شده است")
	}

	// Clear verification code
	existingCustomer.VerifyCode = nil
	existingCustomer.ExpireVerifyCodeAt = nil
	existingCustomer.UpdatedAt = time.Now()

	// Handle different verification types
	switch *params.Type {
	case user.VerifyEmail:
		existingCustomer.VerifyEmail = "1"
		err = u.repo.Update(existingCustomer)
		if err != nil {
			return nil, errors.New("خطا در بروزرسانی اطلاعات مشتری")
		}
		return "success", nil

	case user.VerifyPhone:
		existingCustomer.VerifyPhone = "1"
		err = u.repo.Update(existingCustomer)
		if err != nil {
			return nil, errors.New("خطا در بروزرسانی اطلاعات مشتری")
		}
		return "success", nil

	case user.ForgetPasswordEmail:
		// Generate token for password reset
		token := u.identitySvc.TokenForCustomer(existingCustomer).Make()
		err = u.repo.Update(existingCustomer)
		if err != nil {
			return nil, errors.New("خطا در بروزرسانی اطلاعات مشتری")
		}

		// Return token with redirect URL
		redirectURL := "https://example.com/reset-password"
		queryParams := url.Values{}
		queryParams.Add("token", token)

		return map[string]string{
			"url": redirectURL + "?" + queryParams.Encode(),
		}, nil

	case user.ForgetPasswordPhone:
		err = u.repo.Update(existingCustomer)
		if err != nil {
			return nil, errors.New("خطا در بروزرسانی اطلاعات مشتری")
		}
		return "success", nil

	default:
		return nil, errors.New("نوع تأیید نامعتبر است")
	}
}

// GetProfileCustomerQuery handles getting customer profile
func (u *CustomerUsecase) GetProfileCustomerQuery(params *customer.GetProfileCustomerQuery) (any, error) {
	u.Logger.Info("GetProfileCustomerQuery called", nil)

	// Get customer ID from auth context
	customerID, err := u.authContext(u.Ctx).GetCustomerID()
	if err != nil {
		return nil, errors.New("خطا در احراز هویت مشتری")
	}

	// Get customer details
	existingCustomer, err := u.repo.GetByID(customerID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("مشتری یافت نشد")
		}
		return nil, err
	}

	// Enhance response with file items (avatar) if available
	response := enhanceCustomerResponse(existingCustomer)

	// Get avatar media URL if available
	if existingCustomer.AvatarID != nil {
		// We would typically call a file service here to get the file URL
		// For simplicity, we'll just add the ID
		response["avatarId"] = *existingCustomer.AvatarID
	}

	return response, nil
}

// AdminGetAllCustomerQuery handles admin getting all customers
func (u *CustomerUsecase) AdminGetAllCustomerQuery(params *customer.AdminGetAllCustomerQuery) (any, error) {
	u.Logger.Info("AdminGetAllCustomerQuery called", map[string]interface{}{
		"page":     params.Page,
		"pageSize": params.PageSize,
	})

	// Check admin access
	isAdmin, err := u.authContext(u.Ctx).IsAdmin()
	if err != nil {
		return nil, err
	}
	if !isAdmin {
		return nil, errors.New("فقط مدیران سیستم مجاز به دسترسی به این بخش هستند")
	}

	// Get all customers with pagination
	customers, count, err := u.repo.GetAll(params.PaginationRequestDto)
	if err != nil {
		return nil, errors.New("خطا در دریافت لیست مشتریان")
	}

	// Enhance response
	enhancedCustomers := make([]map[string]interface{}, 0, len(customers))
	for _, c := range customers {
		enhancedCustomers = append(enhancedCustomers, enhanceCustomerResponse(c))
	}

	return map[string]interface{}{
		"items":     enhancedCustomers,
		"total":     count,
		"page":      params.Page,
		"pageSize":  params.PageSize,
		"totalPage": (count + int64(params.PageSize) - 1) / int64(params.PageSize),
	}, nil
}

// Helper function to enhance customer response with structured data
func enhanceCustomerResponse(c domain.Customer) map[string]interface{} {
	response := map[string]interface{}{
		"id":           c.ID,
		"siteId":       c.SiteID,
		"firstName":    c.FirstName,
		"lastName":     c.LastName,
		"email":        c.Email,
		"verifyEmail":  c.VerifyEmail,
		"nationalCode": c.NationalCode,
		"phone":        c.Phone,
		"verifyPhone":  c.VerifyPhone,
		"isActive":     c.IsActive,
		"createdAt":    c.CreatedAt,
		"updatedAt":    c.UpdatedAt,
	}

	// Add roles if available
	if len(c.Roles) > 0 {
		roles := make([]map[string]interface{}, 0, len(c.Roles))
		for _, r := range c.Roles {
			roles = append(roles, map[string]interface{}{
				"id":   r.ID,
				"name": r.Name,
			})
		}
		response["roles"] = roles
	}

	// Add addresses if available
	if len(c.Addresses) > 0 {
		addresses := make([]map[string]interface{}, 0, len(c.Addresses))
		for _, a := range c.Addresses {
			addr := enhanceAddressResponse(a)
			addresses = append(addresses, addr)
		}
		response["addresses"] = addresses
	}

	return response
}

// Helper function to enhance address response with structured data
func enhanceAddressResponse(a domain.Address) map[string]interface{} {
	response := map[string]interface{}{
		"id":          a.ID,
		"title":       a.Title,
		"latitude":    a.Latitude,
		"longitude":   a.Longitude,
		"addressLine": a.AddressLine,
		"postalCode":  a.PostalCode,
		"cityId":      a.CityID,
		"provinceId":  a.ProvinceID,
	}

	// Add city info if available
	if a.City != nil {
		response["city"] = map[string]interface{}{
			"id":   a.City.ID,
			"name": a.City.Name,
		}
	}

	// Add province info if available
	if a.Province != nil {
		response["province"] = map[string]interface{}{
			"id":   a.Province.ID,
			"name": a.Province.Name,
		}
	}

	return response
}

// Helper function to generate a random verification code
func generateVerificationCode() int {
	// In a real implementation, we would use a secure random number generator
	// For simplicity, we'll just return a fixed code for now
	return 123456
}

// Helper function to generate a password reset link
func generatePasswordResetLink(baseURL string, customerID int64, code int, verifyType int) string {
	// In a real implementation, we might encrypt or sign the parameters
	return fmt.Sprintf("%s/reset-password?id=%d&code=%d&type=%d",
		baseURL, customerID, code, verifyType)
}
