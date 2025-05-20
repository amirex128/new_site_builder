package user

import (
	"github.com/amirex128/new_site_builder/src/internal/domain/enums"
)

// LoginUserCommand represents a command to log in a user
type LoginUserCommand struct {
	Email    *string `json:"email" validate:"required,email" nameFa:"ایمیل"`
	Password *string `json:"password" validate:"required_text=8 100" nameFa:"رمز عبور"`
}

// RegisterUserCommand represents a command to register a new user
type RegisterUserCommand struct {
	Email    *string `json:"email" validate:"required,email" nameFa:"ایمیل"`
	Password *string `json:"password" validate:"required_text=8 100" nameFa:"رمز عبور"`
}

// RequestVerifyAndForgetUserCommand represents a command to request verification or password reset
type RequestVerifyAndForgetUserCommand struct {
	Email *string               `json:"email,omitempty" validate:"required_if=Type 0 Type 2,omitempty,email" nameFa:"ایمیل"`
	Phone *string               `json:"phone,omitempty" validate:"required_if=Type 1 Type 3,omitempty,iranian_mobile" nameFa:"تلفن"`
	Type  *enums.VerifyTypeEnum `json:"type" validate:"required,enum" nameFa:"نوع"`
}

// SmptSettings represents SMTP settings for user profile
type SmptSettings struct {
	Host     string `json:"host" validate:"required_text=1,100" nameFa:"آدرس هاست"`
	Port     int    `json:"port" validate:"required,min=1,max=65535" nameFa:"پورت"`
	Username string `json:"username" validate:"required_text=1,100" nameFa:"نام کاربری"`
	Password string `json:"password" validate:"required_text=1,100" nameFa:"رمز عبور"`
}

// UpdateProfileUserCommand represents a command to update a user's profile
type UpdateProfileUserCommand struct {
	FirstName          *string           `json:"firstName,omitempty" validate:"optional_text=1 100" nameFa:"نام"`
	LastName           *string           `json:"lastName,omitempty" validate:"optional_text=1 100" nameFa:"نام خانوادگی"`
	Email              *string           `json:"email,omitempty" validate:"omitempty,email" nameFa:"ایمیل"`
	Password           *string           `json:"password,omitempty" validate:"optional_text=6 100" nameFa:"رمز عبور"`
	NationalCode       *string           `json:"nationalCode,omitempty" validate:"optional_text=1 100" nameFa:"کد ملی"`
	Phone              *string           `json:"phone" validate:"required,iranian_mobile" nameFa:"تلفن"`
	AddressIDs         []int64           `json:"addressIds,omitempty" validate:"array_number_optional=0 100 1 0 false" nameFa:"شناسه آدرس"`
	AiTypeEnum         *enums.AiTypeEnum `json:"aiTypeEnum,omitempty" validate:"enum_optional" nameFa:"نوع AI"`
	UseCustomEmailSmtp *enums.StatusEnum `json:"useCustomEmailSmtp,omitempty" validate:"enum_optional" nameFa:"استفاده از SMTP سفارشی"`
	Smtp               *SmptSettings     `json:"smtp,omitempty" validate:"omitempty" nameFa:"SMTP"`
}

// UnitPriceQuery represents a nested query for unit price in charge credit request
type UnitPriceQuery struct {
	UnitPriceName  *enums.UnitPriceNameEnum `json:"unitPriceName" validate:"required,enum" nameFa:"نام قیمت واحد"`
	UnitPriceCount *int                     `json:"unitPriceCount" validate:"required,min=1,max=1000" nameFa:"تعداد قیمت واحد"`
	UnitPriceDay   *int                     `json:"unitPriceDay,omitempty" validate:"omitempty" nameFa:"تعداد قیمت واحد"`
}

// ChargeCreditRequestUserCommand represents a command to request charging credit
type ChargeCreditRequestUserCommand struct {
	Gateway             *enums.PaymentGatewaysEnum `json:"gateway" validate:"required,enum" nameFa:"درگاه"`
	FinalFrontReturnUrl *string                    `json:"finalFrontReturnUrl" validate:"required_text=1 500" nameFa:"آدرس بازگشت نهایی"`
	UnitPrices          []UnitPriceQuery           `json:"unitPrices" validate:"required,min=1,dive" nameFa:"قیمت واحد"`
}

// ChargeCreditVerifyUserCommand represents a command to verify charge credit
type ChargeCreditVerifyUserCommand struct {
	PaymentStatus *string           `json:"paymentStatus" validate:"required_text=1 50" nameFa:"وضعیت پرداخت"`
	IsSuccess     *bool             `json:"isSuccess" validate:"required_bool" nameFa:"موفقیت"`
	OrderData     map[string]string `json:"orderData" validate:"required" nameFa:"داده سفارش"`
}

// UpgradePlanRequestUserCommand represents a command to request plan upgrade
type UpgradePlanRequestUserCommand struct {
	Gateway             *enums.PaymentGatewaysEnum `json:"gateway" validate:"required,enum" nameFa:"درگاه"`
	FinalFrontReturnUrl *string                    `json:"finalFrontReturnUrl" validate:"required_text=1 500" nameFa:"آدرس بازگشت نهایی"`
	PlanID              *int64                     `json:"planId" validate:"required,gt=0" nameFa:"شناسه طرح"`
}

// UpgradePlanVerifyUserCommand represents a command to verify plan upgrade
type UpgradePlanVerifyUserCommand struct {
	PaymentStatus *string           `json:"paymentStatus" validate:"required_text=1 50" nameFa:"وضعیت پرداخت"`
	IsSuccess     *bool             `json:"isSuccess" validate:"required_bool" nameFa:"موفقیت"`
	OrderData     map[string]string `json:"orderData" validate:"required" nameFa:"داده سفارش"`
}
