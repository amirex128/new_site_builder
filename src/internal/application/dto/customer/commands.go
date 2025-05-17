package customer

import (
	"github.com/amirex128/new_site_builder/src/internal/application/dto/user"
)

// LoginCustomerCommand represents a command to log in a customer
type LoginCustomerCommand struct {
	Email    *string `json:"email" validate:"required,email" error:"required=ایمیل الزامی است|email=ایمیل باید معتبر باشد"`
	Password *string `json:"password" validate:"required,min=8" error:"required=رمز عبور الزامی است|min=رمز عبور باید حداقل 8 کاراکتر باشد"`
}

// RequestVerifyAndForgetCustomerCommand represents a command to request verification or password reset
type RequestVerifyAndForgetCustomerCommand struct {
	Email *string              `json:"email" validate:"required,email,max=200" error:"required=ایمیل الزامی است|email=ایمیل باید معتبر باشد|max=ایمیل نمی‌تواند بیشتر از 200 کاراکتر باشد"`
	Phone *string              `json:"phone" validate:"required,pattern=^09\\d{9}$" error:"required=شماره تلفن الزامی است|pattern=شماره تلفن نامعتبر است"`
	Type  *user.VerifyTypeEnum `json:"type" validate:"required" error:"required=نوع تأیید الزامی است"`
}

// RegisterCustomerCommand represents a command to register a new customer
type RegisterCustomerCommand struct {
	Email    *string `json:"email" validate:"required,email,max=100" error:"required=ایمیل الزامی است|email=ایمیل باید معتبر باشد|max=ایمیل نمی‌تواند بیشتر از 100 کاراکتر باشد"`
	Password *string `json:"password" validate:"required,min=8,max=100" error:"required=رمز عبور الزامی است|min=رمز عبور باید حداقل 8 کاراکتر باشد|max=رمز عبور نمی‌تواند بیشتر از 100 کاراکتر باشد"`
	SiteID   *int64  `json:"siteId" validate:"required,gt=0" error:"required=شناسه سایت الزامی است|gt=شناسه سایت باید بزرگتر از 0 باشد"`
}

// UpdateProfileCustomerCommand represents a command to update a customer's profile
type UpdateProfileCustomerCommand struct {
	FirstName    *string `json:"firstName,omitempty" validate:"omitempty,max=100" error:"max=نام نمی‌تواند بیشتر از 100 کاراکتر باشد"`
	LastName     *string `json:"lastName,omitempty" validate:"omitempty,max=100" error:"max=نام خانوادگی نمی‌تواند بیشتر از 100 کاراکتر باشد"`
	Email        *string `json:"email,omitempty" validate:"omitempty,email" error:"email=فرمت ایمیل نامعتبر است"`
	Password     *string `json:"password,omitempty" validate:"omitempty,min=8,max=100" error:"min=رمز عبور باید حداقل 8 کاراکتر باشد|max=رمز عبور نمی‌تواند بیشتر از 100 کاراکتر باشد"`
	NationalCode *string `json:"nationalCode,omitempty" validate:"omitempty,max=100" error:"max=کد ملی نمی‌تواند بیشتر از 100 کاراکتر باشد"`
	Phone        *string `json:"phone" validate:"required,pattern=^09\\d{9}$" error:"required=شماره تلفن الزامی است|pattern=شماره تلفن نامعتبر است"`
	AddressIDs   []int64 `json:"addressIds,omitempty" validate:"omitempty,dive,gt=0" error:"gt=شناسه‌های آدرس باید بزرگتر از 0 باشند"`
}
