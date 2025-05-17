package payment

// StatusEnum defines active status
type StatusEnum int

const (
	Inactive StatusEnum = iota
	Active
)

// PaymentGatewaysEnum defines payment gateway options
type PaymentGatewaysEnum int

const (
	Saman PaymentGatewaysEnum = iota
	Mellat
	Parsian
	Pasargad
	IranKish
	Melli
	AsanPardakht
	Sepehr
	ZarinPal
	PayIr
	IdPay
	YekPay
	PayPing
	ParbadVirtual
)

// VerifyPaymentEndpointEnum defines endpoints for payment verification
type VerifyPaymentEndpointEnum int

const (
	ChargeCreditVerify VerifyPaymentEndpointEnum = iota
	UpgradePlanVerify
	CreateOrderVerify
)

// UserTypeEnum defines types of users
type UserTypeEnum int

const (
	User UserTypeEnum = iota
	Customer
	Guest
)
