package enums

import (
	"database/sql/driver"
	"errors"
)

type SiteTypeEnum string

const (
	ShopType     SiteTypeEnum = "shop"
	BlogType     SiteTypeEnum = "blog"
	BusinessType SiteTypeEnum = "business"
)

func (e *SiteTypeEnum) Scan(src interface{}) error {
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
	if !SiteTypeEnum(b).IsValid() {
		return errors.New("unsupported type")
	}
	*e = SiteTypeEnum(b)
	return nil
}

func (e SiteTypeEnum) Value() (driver.Value, error) {
	if !e.IsValid() {
		return nil, errors.New("value invalid SiteTypeEnum")
	}
	return string(e), nil
}

func (e SiteTypeEnum) IsValid() bool {
	var siteTypes = []string{
		string(ShopType),
		string(BlogType),
		string(BusinessType),
	}
	for _, siteType := range siteTypes {
		if siteType == string(e) {
			return true
		}
	}
	return false
}
