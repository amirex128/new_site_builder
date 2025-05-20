package payment

import (
	"database/sql/driver"
	"errors"
)

// StatusEnum defines active status
type StatusEnum string

const (
	InactiveStatus StatusEnum = "inactive"
	ActiveStatus   StatusEnum = "active"
)

func (e *StatusEnum) Scan(src interface{}) error {
	var b []byte
	switch src.(type) {
	case []byte:
		b = src.([]byte)
	case string:
		b = []byte(src.(string))
	case nil:
		b = make([]byte, 0)
	default:
		return errors.New("unsupported type")
	}
	if !StatusEnum(b).IsValid() {
		return errors.New("unsupported type")
	}
	*e = StatusEnum(b)
	return nil
}

func (e StatusEnum) Value() (driver.Value, error) {
	if !e.IsValid() {
		return nil, errors.New("value invalid StatusEnum")
	}
	return string(e), nil
}

func (e StatusEnum) IsValid() bool {
	var types = []string{
		string(InactiveStatus),
		string(ActiveStatus),
	}
	for _, t := range types {
		if t == string(e) {
			return true
		}
	}
	return false
}

// PaymentGatewaysEnum defines payment gateway options
type PaymentGatewaysEnum string

const (
	SamanGatewayEnum         PaymentGatewaysEnum = "saman"
	MellatGatewayEnum        PaymentGatewaysEnum = "mellat"
	ParsianGatewayEnum       PaymentGatewaysEnum = "parsian"
	PasargadGatewayEnum      PaymentGatewaysEnum = "pasargad"
	IranKishGatewayEnum      PaymentGatewaysEnum = "irankish"
	MelliGatewayEnum         PaymentGatewaysEnum = "melli"
	AsanPardakhtGatewayEnum  PaymentGatewaysEnum = "asanpardakht"
	SepehrGatewayEnum        PaymentGatewaysEnum = "sepehr"
	ZarinPalGatewayEnum      PaymentGatewaysEnum = "zarinpal"
	PayIrGatewayEnum         PaymentGatewaysEnum = "payir"
	IdPayGatewayEnum         PaymentGatewaysEnum = "idpay"
	YekPayGatewayEnum        PaymentGatewaysEnum = "yekpay"
	PayPingGatewayEnum       PaymentGatewaysEnum = "payping"
	ParbadVirtualGatewayEnum PaymentGatewaysEnum = "parbadvirtual"
)

func (e *PaymentGatewaysEnum) Scan(src interface{}) error {
	var b []byte
	switch src.(type) {
	case []byte:
		b = src.([]byte)
	case string:
		b = []byte(src.(string))
	case nil:
		b = make([]byte, 0)
	default:
		return errors.New("unsupported type")
	}
	if !PaymentGatewaysEnum(b).IsValid() {
		return errors.New("unsupported type")
	}
	*e = PaymentGatewaysEnum(b)
	return nil
}

func (e PaymentGatewaysEnum) Value() (driver.Value, error) {
	if !e.IsValid() {
		return nil, errors.New("value invalid PaymentGatewaysEnum")
	}
	return string(e), nil
}

func (e PaymentGatewaysEnum) IsValid() bool {
	var types = []string{
		string(SamanGatewayEnum),
		string(MellatGatewayEnum),
		string(ParsianGatewayEnum),
		string(PasargadGatewayEnum),
		string(IranKishGatewayEnum),
		string(MelliGatewayEnum),
		string(AsanPardakhtGatewayEnum),
		string(SepehrGatewayEnum),
		string(ZarinPalGatewayEnum),
		string(PayIrGatewayEnum),
		string(IdPayGatewayEnum),
		string(YekPayGatewayEnum),
		string(PayPingGatewayEnum),
		string(ParbadVirtualGatewayEnum),
	}
	for _, t := range types {
		if t == string(e) {
			return true
		}
	}
	return false
}

// VerifyPaymentEndpointEnum defines endpoints for payment verification
type VerifyPaymentEndpointEnum string

const (
	ChargeCreditVerifyEndpoint VerifyPaymentEndpointEnum = "charge_credit_verify"
	UpgradePlanVerifyEndpoint  VerifyPaymentEndpointEnum = "upgrade_plan_verify"
	CreateOrderVerifyEndpoint  VerifyPaymentEndpointEnum = "create_order_verify"
)

func (e *VerifyPaymentEndpointEnum) Scan(src interface{}) error {
	var b []byte
	switch src.(type) {
	case []byte:
		b = src.([]byte)
	case string:
		b = []byte(src.(string))
	case nil:
		b = make([]byte, 0)
	default:
		return errors.New("unsupported type")
	}
	if !VerifyPaymentEndpointEnum(b).IsValid() {
		return errors.New("unsupported type")
	}
	*e = VerifyPaymentEndpointEnum(b)
	return nil
}

func (e VerifyPaymentEndpointEnum) Value() (driver.Value, error) {
	if !e.IsValid() {
		return nil, errors.New("value invalid VerifyPaymentEndpointEnum")
	}
	return string(e), nil
}

func (e VerifyPaymentEndpointEnum) IsValid() bool {
	var types = []string{
		string(ChargeCreditVerifyEndpoint),
		string(UpgradePlanVerifyEndpoint),
		string(CreateOrderVerifyEndpoint),
	}
	for _, t := range types {
		if t == string(e) {
			return true
		}
	}
	return false
}

// UserTypeEnum defines types of users
type UserTypeEnum string

const (
	UserTypeValue     UserTypeEnum = "user"
	CustomerTypeValue UserTypeEnum = "customer"
	GuestTypeValue    UserTypeEnum = "guest"
)

func (e *UserTypeEnum) Scan(src interface{}) error {
	var b []byte
	switch src.(type) {
	case []byte:
		b = src.([]byte)
	case string:
		b = []byte(src.(string))
	case nil:
		b = make([]byte, 0)
	default:
		return errors.New("unsupported type")
	}
	if !UserTypeEnum(b).IsValid() {
		return errors.New("unsupported type")
	}
	*e = UserTypeEnum(b)
	return nil
}

func (e UserTypeEnum) Value() (driver.Value, error) {
	if !e.IsValid() {
		return nil, errors.New("value invalid UserTypeEnum")
	}
	return string(e), nil
}

func (e UserTypeEnum) IsValid() bool {
	var types = []string{
		string(UserTypeValue),
		string(CustomerTypeValue),
		string(GuestTypeValue),
	}
	for _, t := range types {
		if t == string(e) {
			return true
		}
	}
	return false
}
