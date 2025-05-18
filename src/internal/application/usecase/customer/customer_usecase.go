package customerusecase

import (
	"fmt"
	"strconv"
	"time"

	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	"github.com/amirex128/new_site_builder/src/internal/application/dto/customer"
	"github.com/amirex128/new_site_builder/src/internal/application/dto/user"
	"github.com/amirex128/new_site_builder/src/internal/contract"
	"github.com/amirex128/new_site_builder/src/internal/contract/repository"
	"github.com/amirex128/new_site_builder/src/internal/domain"
)

type CustomerUsecase struct {
	logger       sflogger.Logger
	customerRepo repository.ICustomerRepository
	addressRepo  repository.IAddressRepository
	siteRepo     repository.ISiteRepository
}

func NewCustomerUsecase(c contract.IContainer) *CustomerUsecase {
	return &CustomerUsecase{
		logger:       c.GetLogger(),
		customerRepo: c.GetCustomerRepo(),
		addressRepo:  c.GetAddressRepo(),
		siteRepo:     c.GetSiteRepo(),
	}
}

func (u *CustomerUsecase) UpdateProfileCustomerCommand(params *customer.UpdateProfileCustomerCommand) (any, error) {
	// Implementation for updating a customer's profile
	fmt.Println(params)

	// In a real implementation, get the customer ID from the auth context
	customerID := int64(1)

	existingCustomer, err := u.customerRepo.GetByID(customerID)
	if err != nil {
		return nil, err
	}

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
		// In a real implementation, hash the password before storing
		existingCustomer.Password = *params.Password
	}

	if params.NationalCode != nil {
		existingCustomer.NationalCode = *params.NationalCode
	}

	if params.Phone != nil {
		existingCustomer.Phone = *params.Phone
	}

	existingCustomer.UpdatedAt = time.Now()

	err = u.customerRepo.Update(existingCustomer)
	if err != nil {
		return nil, err
	}

	// Handle address IDs if provided
	if len(params.AddressIDs) > 0 {
		// First, remove all existing addresses
		// Implementation would need to handle the many-to-many relationship
		// For now, log that this functionality is not yet implemented
		u.logger.Warnf("Address management for customers not fully implemented")

		// In a complete implementation, would need methods like:
		// err = u.addressRepo.RemoveAllAddressesFromCustomer(customerID)
		// For each address: u.addressRepo.AddAddressToCustomer(addressID, customerID)
	}

	return existingCustomer, nil
}

func (u *CustomerUsecase) GetProfileCustomerQuery(params *customer.GetProfileCustomerQuery) (any, error) {
	// Implementation to get a customer's profile
	fmt.Println(params)

	// In a real implementation, get the customer ID from the auth context
	customerID := int64(1)

	result, err := u.customerRepo.GetByID(customerID)
	if err != nil {
		return nil, err
	}

	// Get customer addresses
	addresses, err := u.addressRepo.GetAllByCustomerID(customerID)
	if err != nil {
		u.logger.Errorf("Failed to get addresses for customer %d: %v", customerID, err)
	}

	return map[string]interface{}{
		"customer":  result,
		"addresses": addresses,
	}, nil
}

func (u *CustomerUsecase) RegisterCustomerCommand(params *customer.RegisterCustomerCommand) (any, error) {
	// Implementation for customer registration
	fmt.Println(params)

	// Check if customer already exists
	_, err := u.customerRepo.GetByEmail(*params.Email)
	if err == nil {
		return nil, fmt.Errorf("customer with email %s already exists", *params.Email)
	}

	// Verify site exists
	_, err = u.siteRepo.GetByID(*params.SiteID)
	if err != nil {
		return nil, fmt.Errorf("site with ID %d not found", *params.SiteID)
	}

	// Generate salt and hash password
	salt := "random_salt" // In a real implementation, generate a secure random salt
	// In a real implementation, hash the password with the salt
	hashedPassword := *params.Password + salt

	// Generate a verification code
	verificationCode := "123456" // Example code, in a real implementation generate a random code

	// Create new customer
	newCustomer := domain.Customer{
		Email:       *params.Email,
		Password:    hashedPassword,
		Salt:        salt,
		SiteID:      *params.SiteID,
		IsActive:    "0", // Activate after verification
		VerifyEmail: verificationCode,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		IsDeleted:   false,
	}

	err = u.customerRepo.Create(newCustomer)
	if err != nil {
		return nil, err
	}

	// TODO: Send verification email or SMS

	return map[string]interface{}{
		"success": true,
		"message": "Registration successful. Please verify your account.",
	}, nil
}

func (u *CustomerUsecase) LoginCustomerCommand(params *customer.LoginCustomerCommand) (any, error) {
	// Implementation for customer login
	fmt.Println(params)

	// Get customer by email
	existingCustomer, err := u.customerRepo.GetByEmail(*params.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	// Check if customer is active
	if existingCustomer.IsActive != "1" {
		return nil, fmt.Errorf("account is not active")
	}

	// In a real implementation, hash the provided password with the stored salt
	// and compare with the stored hashed password
	hashedPassword := *params.Password + existingCustomer.Salt
	if hashedPassword != existingCustomer.Password {
		return nil, fmt.Errorf("invalid email or password")
	}

	// Generate JWT token
	token := "dummy_jwt_token" // In a real implementation, generate a proper JWT token

	return map[string]interface{}{
		"token":    token,
		"customer": existingCustomer,
	}, nil
}

func (u *CustomerUsecase) RequestVerifyAndForgetCustomerCommand(params *customer.RequestVerifyAndForgetCustomerCommand) (any, error) {
	// Implementation for requesting verification or password reset
	fmt.Println(params)

	// Get customer by email or phone
	var existingCustomer domain.Customer
	var err error

	if params.Email != nil {
		existingCustomer, err = u.customerRepo.GetByEmail(*params.Email)
	} else if params.Phone != nil {
		existingCustomer, err = u.customerRepo.GetByPhone(*params.Phone)
	} else {
		return nil, fmt.Errorf("email or phone is required")
	}

	if err != nil {
		return nil, fmt.Errorf("customer not found")
	}

	// Generate verification code
	verificationCode := "123456" // Example code, in a real implementation generate a random code

	// Store verification code based on type
	if params.Type != nil && (*params.Type == user.VerifyEmail || *params.Type == user.ForgetPasswordEmail) {
		existingCustomer.VerifyEmail = verificationCode
	} else if params.Type != nil && (*params.Type == user.VerifyPhone || *params.Type == user.ForgetPasswordPhone) {
		existingCustomer.VerifyPhone = verificationCode
	}

	err = u.customerRepo.Update(existingCustomer)
	if err != nil {
		return nil, err
	}

	// TODO: Send verification code via email or SMS

	return map[string]interface{}{
		"success": true,
		"message": "Verification code sent. Please check your email or phone.",
	}, nil
}

func (u *CustomerUsecase) VerifyCustomerQuery(params *customer.VerifyCustomerQuery) (any, error) {
	// Implementation for customer verification
	fmt.Println(params)

	// Get customer by email
	existingCustomer, err := u.customerRepo.GetByEmail(*params.Email)
	if err != nil {
		return nil, fmt.Errorf("customer not found")
	}

	// Convert code to string for comparison
	codeStr := strconv.Itoa(*params.Code)

	// Check verification code based on type
	if params.Type == nil {
		return nil, fmt.Errorf("verification type is required")
	}

	switch *params.Type {
	case user.VerifyEmail:
		if existingCustomer.VerifyEmail != codeStr {
			return nil, fmt.Errorf("invalid verification code")
		}
		existingCustomer.IsActive = "1"
		existingCustomer.VerifyEmail = ""
	case user.VerifyPhone:
		if existingCustomer.VerifyPhone != codeStr {
			return nil, fmt.Errorf("invalid verification code")
		}
		existingCustomer.IsActive = "1"
		existingCustomer.VerifyPhone = ""
	case user.ForgetPasswordEmail:
		if existingCustomer.VerifyEmail != codeStr {
			return nil, fmt.Errorf("invalid verification code")
		}
		// In a real implementation, provide a token for password reset instead
		existingCustomer.VerifyEmail = ""
	case user.ForgetPasswordPhone:
		if existingCustomer.VerifyPhone != codeStr {
			return nil, fmt.Errorf("invalid verification code")
		}
		// In a real implementation, provide a token for password reset instead
		existingCustomer.VerifyPhone = ""
	default:
		return nil, fmt.Errorf("invalid verification type")
	}

	err = u.customerRepo.Update(existingCustomer)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"success": true,
		"message": "Verification successful.",
	}, nil
}

func (u *CustomerUsecase) AdminGetAllCustomerQuery(params *customer.AdminGetAllCustomerQuery) (any, error) {
	// Implementation for admin to get all customers
	fmt.Println(params)

	result, count, err := u.customerRepo.GetAll(params.PaginationRequestDto)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"items": result,
		"total": count,
	}, nil
}
