package enums

import (
	"database/sql/driver"
	"errors"
)

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
