package enums

import (
	"database/sql/driver"
	"errors"
)

type VerifyTypeEnum string

const (
	VerifyEmailType         VerifyTypeEnum = "verify_email"
	VerifyPhoneType         VerifyTypeEnum = "verify_phone"
	ForgetPasswordEmailType VerifyTypeEnum = "forget_password_email"
	ForgetPasswordPhoneType VerifyTypeEnum = "forget_password_phone"
)

func (e *VerifyTypeEnum) Scan(src interface{}) error {
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
	if !VerifyTypeEnum(b).IsValid() {
		return errors.New("unsupported type")
	}
	*e = VerifyTypeEnum(b)
	return nil
}

func (e VerifyTypeEnum) Value() (driver.Value, error) {
	if !e.IsValid() {
		return nil, errors.New("value invalid VerifyTypeEnum")
	}
	return string(e), nil
}

func (e VerifyTypeEnum) IsValid() bool {
	var types = []string{
		string(VerifyEmailType),
		string(VerifyPhoneType),
		string(ForgetPasswordEmailType),
		string(ForgetPasswordPhoneType),
	}
	for _, t := range types {
		if t == string(e) {
			return true
		}
	}
	return false
}
