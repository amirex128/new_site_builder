package customer_ticket

import (
	"database/sql/driver"
	"errors"
)

// CustomerTicketStatusEnum defines the status of a customer ticket
type CustomerTicketStatusEnum string

const (
	NewStatus        CustomerTicketStatusEnum = "new"
	InProgressStatus CustomerTicketStatusEnum = "in_progress"
	ClosedStatus     CustomerTicketStatusEnum = "closed"
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

// IsValid try to validate enum value on this type
func (e CustomerTicketStatusEnum) IsValid() bool {
	var statusTypes = []string{
		string(NewStatus),
		string(InProgressStatus),
		string(ClosedStatus),
	}

	for _, statusType := range statusTypes {
		if statusType == string(e) {
			return true
		}
	}
	return false
}

// CustomerTicketCategoryEnum defines the category of a customer ticket
type CustomerTicketCategoryEnum string

const (
	BugCategory            CustomerTicketCategoryEnum = "bug"
	EnhancementCategory    CustomerTicketCategoryEnum = "enhancement"
	FeatureRequestCategory CustomerTicketCategoryEnum = "feature_request"
	QuestionCategory       CustomerTicketCategoryEnum = "question"
	DocumentationCategory  CustomerTicketCategoryEnum = "documentation"
	FinancialCategory      CustomerTicketCategoryEnum = "financial"
)

func (e *CustomerTicketCategoryEnum) Scan(src interface{}) error {
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
	if !CustomerTicketCategoryEnum(b).IsValid() {
		return errors.New("unsupported type")
	}
	*e = CustomerTicketCategoryEnum(b)
	return nil
}

func (e CustomerTicketCategoryEnum) Value() (driver.Value, error) {
	if !e.IsValid() {
		return nil, errors.New("value invalid CustomerTicketCategoryEnum")
	}
	return string(e), nil
}

// IsValid try to validate enum value on this type
func (e CustomerTicketCategoryEnum) IsValid() bool {
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

// CustomerTicketPriorityEnum defines the priority of a customer ticket
type CustomerTicketPriorityEnum string

const (
	LowPriority      CustomerTicketPriorityEnum = "low"
	MediumPriority   CustomerTicketPriorityEnum = "medium"
	HighPriority     CustomerTicketPriorityEnum = "high"
	CriticalPriority CustomerTicketPriorityEnum = "critical"
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

// IsValid try to validate enum value on this type
func (e CustomerTicketPriorityEnum) IsValid() bool {
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
