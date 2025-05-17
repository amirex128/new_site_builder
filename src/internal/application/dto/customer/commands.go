package customer

import (
	"github.com/amirex128/new_site_builder/src/internal/application/dto/user"
)

// LoginCustomerCommand represents a command to log in a customer
type LoginCustomerCommand struct {
	Email    *string `json:"email" validate:"required,email"`
	Password *string `json:"password" validate:"required_text=8,100"`
}

// RequestVerifyAndForgetCustomerCommand represents a command to request verification or password reset
type RequestVerifyAndForgetCustomerCommand struct {
	Email *string              `json:"email" validate:"required,email"`
	Phone *string              `json:"phone" validate:"required,iranian_mobile"`
	Type  *user.VerifyTypeEnum `json:"type" validate:"required,enum"`
}

// RegisterCustomerCommand represents a command to register a new customer
type RegisterCustomerCommand struct {
	Email    *string `json:"email" validate:"required,email"`
	Password *string `json:"password" validate:"required_text=8,100"`
	SiteID   *int64  `json:"siteId" validate:"required,gt=0"`
}

// UpdateProfileCustomerCommand represents a command to update a customer's profile
type UpdateProfileCustomerCommand struct {
	FirstName    *string `json:"firstName,omitempty" validate:"optional_text=1,100"`
	LastName     *string `json:"lastName,omitempty" validate:"optional_text=1,100"`
	Email        *string `json:"email,omitempty" validate:"omitempty,email"`
	Password     *string `json:"password,omitempty" validate:"optional_text=8,100"`
	NationalCode *string `json:"nationalCode,omitempty" validate:"optional_text=1,100"`
	Phone        *string `json:"phone" validate:"required,iranian_mobile"`
	AddressIDs   []int64 `json:"addressIds,omitempty" validate:"array_number_optional=0,100,1,0,false"`
}
