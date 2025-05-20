package enums

import (
	"database/sql/driver"
	"errors"
)

type CustomerTicketCategoryEnum string

const (
	CustomerTicketBugCategory            CustomerTicketCategoryEnum = "bug"
	CustomerTicketEnhancementCategory    CustomerTicketCategoryEnum = "enhancement"
	CustomerTicketFeatureRequestCategory CustomerTicketCategoryEnum = "feature_request"
	CustomerTicketQuestionCategory       CustomerTicketCategoryEnum = "question"
	CustomerTicketDocumentationCategory  CustomerTicketCategoryEnum = "documentation"
	CustomerTicketFinancialCategory      CustomerTicketCategoryEnum = "financial"
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

func (e CustomerTicketCategoryEnum) IsValid() bool {
	var categoryTypes = []string{
		string(CustomerTicketBugCategory),
		string(CustomerTicketEnhancementCategory),
		string(CustomerTicketFeatureRequestCategory),
		string(CustomerTicketQuestionCategory),
		string(CustomerTicketDocumentationCategory),
		string(CustomerTicketFinancialCategory),
	}
	for _, categoryType := range categoryTypes {
		if categoryType == string(e) {
			return true
		}
	}
	return false
}
