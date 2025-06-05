package header_footer

import (
	"database/sql/driver"
	"errors"
)

// HeaderFooterTypeEnum defines header/footer types
type HeaderFooterTypeEnum string

const (
	HeaderType HeaderFooterTypeEnum = "header"
	FooterType HeaderFooterTypeEnum = "footer"
)

func (e *HeaderFooterTypeEnum) Scan(src interface{}) error {
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
	if !HeaderFooterTypeEnum(b).IsValid() {
		return errors.New("unsupported type")
	}
	*e = HeaderFooterTypeEnum(b)
	return nil
}

func (e HeaderFooterTypeEnum) Value() (driver.Value, error) {
	if !e.IsValid() {
		return nil, errors.New("value invalid HeaderFooterTypeEnum")
	}
	return string(e), nil
}

// IsValid try to validate enum value on this type
func (e HeaderFooterTypeEnum) IsValid() bool {
	var typeValues = []string{
		string(HeaderType),
		string(FooterType),
	}

	for _, typeValue := range typeValues {
		if typeValue == string(e) {
			return true
		}
	}
	return false
}
