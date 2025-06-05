package enums

import (
	"database/sql/driver"
	"errors"
)

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
