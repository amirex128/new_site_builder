package enums

import (
	"database/sql/driver"
	"errors"
)

type DiscountTypeEnum string

const (
	FixedDiscountType      DiscountTypeEnum = "fixed"
	PercentageDiscountType DiscountTypeEnum = "percentage"
)

func (e *DiscountTypeEnum) Scan(src interface{}) error {
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
	if !DiscountTypeEnum(b).IsValid() {
		return errors.New("unsupported type")
	}
	*e = DiscountTypeEnum(b)
	return nil
}

func (e DiscountTypeEnum) Value() (driver.Value, error) {
	if !e.IsValid() {
		return nil, errors.New("value invalid DiscountTypeEnum")
	}
	return string(e), nil
}

func (e DiscountTypeEnum) IsValid() bool {
	var discountTypes = []string{
		string(FixedDiscountType),
		string(PercentageDiscountType),
	}
	for _, discountType := range discountTypes {
		if discountType == string(e) {
			return true
		}
	}
	return false
}
