package enums

import (
	"database/sql/driver"
	"errors"
)

type StatusEnum string

const (
	InactiveStatus StatusEnum = "inactive"
	ActiveStatus   StatusEnum = "active"
	PendingStatus  StatusEnum = "pending"
	DeletedStatus  StatusEnum = "deleted"
	DisabledStatus StatusEnum = "disabled"
	EnabledStatus  StatusEnum = "enabled"
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

func (e StatusEnum) IsValid() bool {
	var statusTypes = []string{
		string(InactiveStatus),
		string(ActiveStatus),
		string(PendingStatus),
		string(DeletedStatus),
		string(DisabledStatus),
		string(EnabledStatus),
	}
	for _, statusType := range statusTypes {
		if statusType == string(e) {
			return true
		}
	}
	return false
}
