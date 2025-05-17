package order

// CourierEnum defines types of couriers for shipping
type CourierEnum int

const (
	Post CourierEnum = iota
	Tipax
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
