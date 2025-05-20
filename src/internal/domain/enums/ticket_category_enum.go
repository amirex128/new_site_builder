package enums

import (
	"database/sql/driver"
	"errors"
)

type TicketCategoryEnum string

const (
	TicketBugCategory            TicketCategoryEnum = "bug"
	TicketEnhancementCategory    TicketCategoryEnum = "enhancement"
	TicketFeatureRequestCategory TicketCategoryEnum = "feature_request"
	TicketQuestionCategory       TicketCategoryEnum = "question"
	TicketDocumentationCategory  TicketCategoryEnum = "documentation"
	TicketFinancialCategory      TicketCategoryEnum = "financial"
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

func (e TicketCategoryEnum) IsValid() bool {
	var categoryTypes = []string{
		string(TicketBugCategory),
		string(TicketEnhancementCategory),
		string(TicketFeatureRequestCategory),
		string(TicketQuestionCategory),
		string(TicketDocumentationCategory),
		string(TicketFinancialCategory),
	}
	for _, categoryType := range categoryTypes {
		if categoryType == string(e) {
			return true
		}
	}
	return false
}
