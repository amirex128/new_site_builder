package ticket

import (
	"database/sql/driver"
	"errors"
)

// TicketStatusEnum defines the status of a ticket
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

// IsValid try to validate enum value on this type
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

// TicketCategoryEnum defines the category of a ticket
type TicketCategoryEnum string

const (
	BugCategory            TicketCategoryEnum = "bug"
	EnhancementCategory    TicketCategoryEnum = "enhancement"
	FeatureRequestCategory TicketCategoryEnum = "feature_request"
	QuestionCategory       TicketCategoryEnum = "question"
	DocumentationCategory  TicketCategoryEnum = "documentation"
	FinancialCategory      TicketCategoryEnum = "financial"
)

func (e *TicketCategoryEnum) Scan(src interface{}) error {
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
	if !TicketCategoryEnum(b).IsValid() {
		return errors.New("unsupported type")
	}
	*e = TicketCategoryEnum(b)
	return nil
}

func (e TicketCategoryEnum) Value() (driver.Value, error) {
	if !e.IsValid() {
		return nil, errors.New("value invalid TicketCategoryEnum")
	}
	return string(e), nil
}

// IsValid try to validate enum value on this type
func (e TicketCategoryEnum) IsValid() bool {
	var categoryTypes = []string{
		string(BugCategory),
		string(EnhancementCategory),
		string(FeatureRequestCategory),
		string(QuestionCategory),
		string(DocumentationCategory),
		string(FinancialCategory),
	}

	for _, categoryType := range categoryTypes {
		if categoryType == string(e) {
			return true
		}
	}
	return false
}

// TicketPriorityEnum defines the priority of a ticket
type TicketPriorityEnum string

const (
	LowPriority      TicketPriorityEnum = "low"
	MediumPriority   TicketPriorityEnum = "medium"
	HighPriority     TicketPriorityEnum = "high"
	CriticalPriority TicketPriorityEnum = "critical"
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

// IsValid try to validate enum value on this type
func (e TicketPriorityEnum) IsValid() bool {
	var priorityTypes = []string{
		string(LowPriority),
		string(MediumPriority),
		string(HighPriority),
		string(CriticalPriority),
	}

	for _, priorityType := range priorityTypes {
		if priorityType == string(e) {
			return true
		}
	}
	return false
}
