package enums

import (
	"database/sql/driver"
	"errors"
)

type TicketPriorityEnum string

const (
	TicketLowPriority      TicketPriorityEnum = "low"
	TicketMediumPriority   TicketPriorityEnum = "medium"
	TicketHighPriority     TicketPriorityEnum = "high"
	TicketCriticalPriority TicketPriorityEnum = "critical"
)

func (e *TicketPriorityEnum) Scan(src interface{}) error {
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
	if !TicketPriorityEnum(b).IsValid() {
		return errors.New("unsupported type")
	}
	*e = TicketPriorityEnum(b)
	return nil
}

func (e TicketPriorityEnum) Value() (driver.Value, error) {
	if !e.IsValid() {
		return nil, errors.New("value invalid TicketPriorityEnum")
	}
	return string(e), nil
}

func (e TicketPriorityEnum) IsValid() bool {
	var priorityTypes = []string{
		string(TicketLowPriority),
		string(TicketMediumPriority),
		string(TicketHighPriority),
		string(TicketCriticalPriority),
	}
	for _, priorityType := range priorityTypes {
		if priorityType == string(e) {
			return true
		}
	}
	return false
}
