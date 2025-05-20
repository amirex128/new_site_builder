package user

// LoginUserCommand represents a command to log in a user
type LoginUserCommand struct {
	Email    *string `json:"email" validate:"required,email"`
	Password *string `json:"password" validate:"required_text=8 100"`
}

// RegisterUserCommand represents a command to register a new user
type RegisterUserCommand struct {
	Email    *string `json:"email" validate:"required,email"`
	Password *string `json:"password" validate:"required_text=8 100"`
}

// RequestVerifyAndForgetUserCommand represents a command to request verification or password reset
type RequestVerifyAndForgetUserCommand struct {
	Email *string         `json:"email,omitempty" validate:"required_if=Type 0 Type 2,omitempty,email"`
	Phone *string         `json:"phone,omitempty" validate:"required_if=Type 1 Type 3,omitempty,iranian_mobile"`
	Type  *VerifyTypeEnum `json:"type" validate:"required,enum"`
}

// SmptSettings represents SMTP settings for user profile
type SmptSettings struct {
	Host     string `json:"host" validate:"required_text=1,100"`
	Port     int    `json:"port" validate:"required,min=1,max=65535"`
	Username string `json:"username" validate:"required_text=1,100"`
	Password string `json:"password" validate:"required_text=1,100"`
}

// UpdateProfileUserCommand represents a command to update a user's profile
type UpdateProfileUserCommand struct {
	FirstName          *string       `json:"firstName,omitempty" validate:"optional_text=1 100"`
	LastName           *string       `json:"lastName,omitempty" validate:"optional_text=1 100"`
	Email              *string       `json:"email,omitempty" validate:"omitempty,email"`
	Password           *string       `json:"password,omitempty" validate:"optional_text=6 100"`
	NationalCode       *string       `json:"nationalCode,omitempty" validate:"optional_text=1 100"`
	Phone              *string       `json:"phone" validate:"required,iranian_mobile"`
	AddressIDs         []int64       `json:"addressIds,omitempty" validate:"array_number_optional=0 100 1 0 false"`
	AiTypeEnum         *AiTypeEnum   `json:"aiTypeEnum,omitempty" validate:"enum_optional"`
	UseCustomEmailSmtp *StatusEnum   `json:"useCustomEmailSmtp,omitempty" validate:"enum_optional"`
	Smtp               *SmptSettings `json:"smtp,omitempty" validate:"omitempty"`
}

// UnitPriceQuery represents a nested query for unit price in charge credit request
type UnitPriceQuery struct {
	UnitPriceName  *UnitPriceNameEnum `json:"unitPriceName" validate:"required,enum"`
	UnitPriceCount *int               `json:"unitPriceCount" validate:"required,min=1,max=1000"`
	UnitPriceDay   *int               `json:"unitPriceDay,omitempty" validate:"omitempty"`
}

// ChargeCreditRequestUserCommand represents a command to request charging credit
type ChargeCreditRequestUserCommand struct {
	Gateway             *PaymentGatewaysEnum `json:"gateway" validate:"required,enum"`
	FinalFrontReturnUrl *string              `json:"finalFrontReturnUrl" validate:"required_text=1 500"`
	UnitPrices          []UnitPriceQuery     `json:"unitPrices" validate:"required,min=1,dive"`
}

// ChargeCreditVerifyUserCommand represents a command to verify charge credit
type ChargeCreditVerifyUserCommand struct {
	PaymentStatus *string           `json:"paymentStatus" validate:"required_text=1 50"`
	IsSuccess     *bool             `json:"isSuccess" validate:"required_bool"`
	OrderData     map[string]string `json:"orderData" validate:"required"`
}

// UpgradePlanRequestUserCommand represents a command to request plan upgrade
type UpgradePlanRequestUserCommand struct {
	Gateway             *PaymentGatewaysEnum `json:"gateway" validate:"required,enum"`
	FinalFrontReturnUrl *string              `json:"finalFrontReturnUrl" validate:"required_text=1 500"`
	PlanID              *int64               `json:"planId" validate:"required,gt=0"`
}

// UpgradePlanVerifyUserCommand represents a command to verify plan upgrade
type UpgradePlanVerifyUserCommand struct {
	PaymentStatus *string           `json:"paymentStatus" validate:"required_text=1 50"`
	IsSuccess     *bool             `json:"isSuccess" validate:"required_bool"`
	OrderData     map[string]string `json:"orderData" validate:"required"`
}
