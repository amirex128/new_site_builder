package page_usage

import (
	"database/sql/driver"
	"errors"
)

// PageUsageEnum defines usage types for pages
type PageUsageEnum string

const (
	ProductUsage      PageUsageEnum = "product"
	ArticleUsage      PageUsageEnum = "article"
	HeaderFooterUsage PageUsageEnum = "header_footer"
)

func (e *PageUsageEnum) Scan(src interface{}) error {
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
	if !PageUsageEnum(b).IsValid() {
		return errors.New("unsupported type")
	}
	*e = PageUsageEnum(b)
	return nil
}

func (e PageUsageEnum) Value() (driver.Value, error) {
	if !e.IsValid() {
		return nil, errors.New("value invalid PageUsageEnum")
	}
	return string(e), nil
}

// IsValid try to validate enum value on this type
func (e PageUsageEnum) IsValid() bool {
	var usageTypes = []string{
		string(ProductUsage),
		string(ArticleUsage),
		string(HeaderFooterUsage),
	}

	for _, usageType := range usageTypes {
		if usageType == string(e) {
			return true
		}
	}
	return false
}
