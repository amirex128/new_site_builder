package enums

import (
	"database/sql/driver"
	"errors"
)

type CustomerTicketPriorityEnum string

const (
	CustomerTicketLowPriority      CustomerTicketPriorityEnum = "low"
	CustomerTicketMediumPriority   CustomerTicketPriorityEnum = "medium"
	CustomerTicketHighPriority     CustomerTicketPriorityEnum = "high"
	CustomerTicketCriticalPriority CustomerTicketPriorityEnum = "critical"
)

func (e *CustomerTicketPriorityEnum) Scan(src interface{}) error {
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
	if !CustomerTicketPriorityEnum(b).IsValid() {
		return errors.New("unsupported type")
	}
	*e = CustomerTicketPriorityEnum(b)
	return nil
}

func (e CustomerTicketPriorityEnum) Value() (driver.Value, error) {
	if !e.IsValid() {
		return nil, errors.New("value invalid CustomerTicketPriorityEnum")
	}
	return string(e), nil
}

func (e CustomerTicketPriorityEnum) IsValid() bool {
	var priorityTypes = []string{
		string(CustomerTicketLowPriority),
		string(CustomerTicketMediumPriority),
		string(CustomerTicketHighPriority),
		string(CustomerTicketCriticalPriority),
	}
	for _, priorityType := range priorityTypes {
		if priorityType == string(e) {
			return true
		}
	}
	return false
}
