package enums

import (
	"database/sql/driver"
	"errors"
)

type UserTypeEnum string

const (
	UserTypeValue     UserTypeEnum = "user"
	CustomerTypeValue UserTypeEnum = "customer"
	GuestTypeValue    UserTypeEnum = "guest"
)

func (e *UserTypeEnum) Scan(src interface{}) error {
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
	if !UserTypeEnum(b).IsValid() {
		return errors.New("unsupported type")
	}
	*e = UserTypeEnum(b)
	return nil
}

func (e UserTypeEnum) Value() (driver.Value, error) {
	if !e.IsValid() {
		return nil, errors.New("value invalid UserTypeEnum")
	}
	return string(e), nil
}

func (e UserTypeEnum) IsValid() bool {
	var types = []string{
		string(UserTypeValue),
		string(CustomerTypeValue),
		string(GuestTypeValue),
	}
	for _, t := range types {
		if t == string(e) {
			return true
		}
	}
	return false
}
