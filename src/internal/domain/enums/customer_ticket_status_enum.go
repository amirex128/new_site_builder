package enums

import (
	"database/sql/driver"
	"errors"
)

type CustomerTicketStatusEnum string

const (
	CustomerTicketNewStatus        CustomerTicketStatusEnum = "new"
	CustomerTicketInProgressStatus CustomerTicketStatusEnum = "in_progress"
	CustomerTicketClosedStatus     CustomerTicketStatusEnum = "closed"
)

func (e *CustomerTicketStatusEnum) Scan(src interface{}) error {
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
	if !CustomerTicketStatusEnum(b).IsValid() {
		return errors.New("unsupported type")
	}
	*e = CustomerTicketStatusEnum(b)
	return nil
}

func (e CustomerTicketStatusEnum) Value() (driver.Value, error) {
	if !e.IsValid() {
		return nil, errors.New("value invalid CustomerTicketStatusEnum")
	}
	return string(e), nil
}

func (e CustomerTicketStatusEnum) IsValid() bool {
	var statusTypes = []string{
		string(CustomerTicketNewStatus),
		string(CustomerTicketInProgressStatus),
		string(CustomerTicketClosedStatus),
	}
	for _, statusType := range statusTypes {
		if statusType == string(e) {
			return true
		}
	}
	return false
}
