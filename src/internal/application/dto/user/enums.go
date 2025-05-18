package user

// VerifyTypeEnum defines verification types
type VerifyTypeEnum int

const (
	VerifyEmail VerifyTypeEnum = iota
	VerifyPhone
	ForgetPasswordEmail
	ForgetPasswordPhone
)

// UnitPriceNameEnum defines unit price types
type UnitPriceNameEnum int

const (
	StorageMbCredits UnitPriceNameEnum = iota
	PageViewCredits
	FormSubmitCredits
	SiteCredits
	SmsCredits
	EmailCredits
	AiCredits
	AiImageCredits
)

// PaymentGatewaysEnum defines payment gateway types
type PaymentGatewaysEnum int

const (
	ZarinPal PaymentGatewaysEnum = iota
	IDPay
	NextPay
)

// ServiceNameEnum defines service names
type ServiceNameEnum int

const (
	User ServiceNameEnum = iota
	Order
	Payment
)

// VerifyPaymentEndpointEnum defines verify payment endpoint types
type VerifyPaymentEndpointEnum int

const (
	ChargeCreditVerify VerifyPaymentEndpointEnum = iota
	UpgradePlanVerify
)

// UserTypeEnum defines user types
type UserTypeEnum int

const (
	UserType UserTypeEnum = iota
	CustomerType
)

// StatusEnum defines status types
type StatusEnum int

const (
	Disabled StatusEnum = iota
	Enabled
)

// AiTypeEnum defines AI types
type AiTypeEnum int

const (
	GPT35 AiTypeEnum = iota
	GPT4
	Claude
)

// DiscountType defines discount types
type DiscountType int

const (
	Fixed DiscountType = iota
	Percentage
)
