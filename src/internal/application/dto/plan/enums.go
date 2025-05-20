package plan

import (
	"database/sql/driver"
	"errors"
)

// DiscountTypeEnum defines discount types
type DiscountTypeEnum string

const (
	FixedDiscount      DiscountTypeEnum = "fixed"
	PercentageDiscount DiscountTypeEnum = "percentage"
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

// IsValid try to validate enum value on this type
func (e DiscountTypeEnum) IsValid() bool {
	var discountTypes = []string{
		string(FixedDiscount),
		string(PercentageDiscount),
	}

	for _, discountType := range discountTypes {
		if discountType == string(e) {
			return true
		}
	}
	return false
}
