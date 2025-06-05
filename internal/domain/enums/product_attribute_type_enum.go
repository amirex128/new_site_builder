package enums

import (
	"database/sql/driver"
	"errors"
)

type ProductAttributeTypeEnum string

const (
	PublicProductAttributeType    ProductAttributeTypeEnum = "public"
	TechnicalProductAttributeType ProductAttributeTypeEnum = "technical"
	OtherProductAttributeType     ProductAttributeTypeEnum = "other"
)

func (e *ProductAttributeTypeEnum) Scan(src interface{}) error {
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
	if !ProductAttributeTypeEnum(b).IsValid() {
		return errors.New("unsupported type")
	}
	*e = ProductAttributeTypeEnum(b)
	return nil
}

func (e ProductAttributeTypeEnum) Value() (driver.Value, error) {
	if !e.IsValid() {
		return nil, errors.New("value invalid ProductAttributeTypeEnum")
	}
	return string(e), nil
}

func (e ProductAttributeTypeEnum) IsValid() bool {
	switch e {
	case PublicProductAttributeType, TechnicalProductAttributeType, OtherProductAttributeType:
		return true
	}
	return false
}
