package customer

import (
	"github.com/amirex128/new_site_builder/src/internal/domain/enums"
)

// LoginCustomerCommand represents a command to log in a customer
type LoginCustomerCommand struct {
	Email    *string `json:"email" nameFa:"ایمیل" validate:"required,email"`
	Password *string `json:"password" nameFa:"رمز عبور" validate:"required_text=8,100"`
}

// RequestVerifyAndForgetCustomerCommand represents a command to request verification or password reset
type RequestVerifyAndForgetCustomerCommand struct {
	Email *string               `json:"email" nameFa:"ایمیل" validate:"required,email"`
	Phone *string               `json:"phone" nameFa:"شماره تلفن" validate:"required,iranian_mobile"`
	Type  *enums.VerifyTypeEnum `json:"type" nameFa:"نوع تایید" validate:"required,enum"`
}

// RegisterCustomerCommand represents a command to register a new customer
type RegisterCustomerCommand struct {
	Email    *string `json:"email" nameFa:"ایمیل" validate:"required,email"`
	Password *string `json:"password" nameFa:"رمز عبور" validate:"required_text=8,100"`
	SiteID   *int64  `json:"siteId" nameFa:"شناسه سایت" validate:"required,gt=0"`
}

// UpdateProfileCustomerCommand represents a command to update a customer's profile
type UpdateProfileCustomerCommand struct {
	FirstName    *string `json:"firstName,omitempty" validate:"optional_text=1 100"`
	LastName     *string `json:"lastName,omitempty" nameFa:"نام خانوادگی" validate:"optional_text=1 100"`
	Email        *string `json:"email,omitempty" nameFa:"ایمیل" validate:"omitempty,email"`
	Password     *string `json:"password,omitempty" nameFa:"رمز عبور" validate:"optional_text=8 100"`
	NationalCode *string `json:"nationalCode,omitempty" nameFa:"کد ملی" validate:"optional_text=1 100"`
	Phone        *string `json:"phone" nameFa:"شماره تلفن" validate:"required,iranian_mobile"`
	AddressIDs   []int64 `json:"addressIds,omitempty" nameFa:"شناسه های آدرس" validate:"array_number_optional=0 100 1 0 false"`
}
