package enums

import (
	"database/sql/driver"
	"errors"
)

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
