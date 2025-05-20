package enums

import (
	"database/sql/driver"
	"errors"
)

type TicketStatusEnum string

const (
	NewTicketStatus        TicketStatusEnum = "new"
	InProgressTicketStatus TicketStatusEnum = "in_progress"
	ClosedTicketStatus     TicketStatusEnum = "closed"
)

func (e *TicketStatusEnum) Scan(src interface{}) error {
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
	if !TicketStatusEnum(b).IsValid() {
		return errors.New("unsupported type")
	}
	*e = TicketStatusEnum(b)
	return nil
}

func (e TicketStatusEnum) Value() (driver.Value, error) {
	if !e.IsValid() {
		return nil, errors.New("value invalid TicketStatusEnum")
	}
	return string(e), nil
}

func (e TicketStatusEnum) IsValid() bool {
	var statusTypes = []string{
		string(NewTicketStatus),
		string(InProgressTicketStatus),
		string(ClosedTicketStatus),
	}
	for _, statusType := range statusTypes {
		if statusType == string(e) {
			return true
		}
	}
	return false
}
