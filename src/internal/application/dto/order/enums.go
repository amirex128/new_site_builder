package order

import (
	"database/sql/driver"
	"errors"
)

// CourierEnum defines types of couriers for shipping
type CourierEnum string

const (
	PostCourier  CourierEnum = "post"
	TipaxCourier CourierEnum = "tipax"
)

func (e *CourierEnum) Scan(src interface{}) error {
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
	if !CourierEnum(b).IsValid() {
		return errors.New("unsupported type")
	}
	*e = CourierEnum(b)
	return nil
}

func (e CourierEnum) Value() (driver.Value, error) {
	if !e.IsValid() {
		return nil, errors.New("value invalid CourierEnum")
	}
	return string(e), nil
}

func (e CourierEnum) IsValid() bool {
	var types = []string{
		string(PostCourier),
		string(TipaxCourier),
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
