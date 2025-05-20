package site

import (
	"database/sql/driver"
	"errors"
)

// DomainTypeEnum defines domain types
type DomainTypeEnum string

const (
	DomainType    DomainTypeEnum = "domain"
	SubdomainType DomainTypeEnum = "subdomain"
)

func (e *DomainTypeEnum) Scan(src interface{}) error {
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
	if !DomainTypeEnum(b).IsValid() {
		return errors.New("unsupported type")
	}
	*e = DomainTypeEnum(b)
	return nil
}

func (e DomainTypeEnum) Value() (driver.Value, error) {
	if !e.IsValid() {
		return nil, errors.New("value invalid DomainTypeEnum")
	}
	return string(e), nil
}

// IsValid try to validate enum value on this type
func (e DomainTypeEnum) IsValid() bool {
	var domainTypes = []string{
		string(DomainType),
		string(SubdomainType),
	}

	for _, domainType := range domainTypes {
		if domainType == string(e) {
			return true
		}
	}
	return false
}

// SiteTypeEnum defines site types
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

// IsValid try to validate enum value on this type
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

// StatusEnum defines status types
type StatusEnum string

const (
	ActiveStatus   StatusEnum = "active"
	InactiveStatus StatusEnum = "inactive"
	PendingStatus  StatusEnum = "pending"
	DeletedStatus  StatusEnum = "deleted"
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

// IsValid try to validate enum value on this type
func (e StatusEnum) IsValid() bool {
	var statusTypes = []string{
		string(ActiveStatus),
		string(InactiveStatus),
		string(PendingStatus),
		string(DeletedStatus),
	}

	for _, statusType := range statusTypes {
		if statusType == string(e) {
			return true
		}
	}
	return false
}
