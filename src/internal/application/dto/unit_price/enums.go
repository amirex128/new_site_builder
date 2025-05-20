package unit_price

import (
	"database/sql/driver"
	"errors"
)

// UnitPriceNameEnum defines unit price types
type UnitPriceNameEnum string

const (
	StorageMbCreditsName  UnitPriceNameEnum = "storage_mb_credits"
	PageViewCreditsName   UnitPriceNameEnum = "page_view_credits"
	FormSubmitCreditsName UnitPriceNameEnum = "form_submit_credits"
	SiteCreditsName       UnitPriceNameEnum = "site_credits"
	SmsCreditsName        UnitPriceNameEnum = "sms_credits"
	EmailCreditsName      UnitPriceNameEnum = "email_credits"
	AiCreditsName         UnitPriceNameEnum = "ai_credits"
	AiImageCreditsName    UnitPriceNameEnum = "ai_image_credits"
)

func (e *UnitPriceNameEnum) Scan(src interface{}) error {
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
	if !UnitPriceNameEnum(b).IsValid() {
		return errors.New("unsupported type")
	}
	*e = UnitPriceNameEnum(b)
	return nil
}

func (e UnitPriceNameEnum) Value() (driver.Value, error) {
	if !e.IsValid() {
		return nil, errors.New("value invalid UnitPriceNameEnum")
	}
	return string(e), nil
}

func (e UnitPriceNameEnum) IsValid() bool {
	var types = []string{
		string(StorageMbCreditsName),
		string(PageViewCreditsName),
		string(FormSubmitCreditsName),
		string(SiteCreditsName),
		string(SmsCreditsName),
		string(EmailCreditsName),
		string(AiCreditsName),
		string(AiImageCreditsName),
	}
	for _, t := range types {
		if t == string(e) {
			return true
		}
	}
	return false
}

// DiscountTypeEnum defines discount types
type DiscountTypeEnum string

const (
	FixedDiscountType      DiscountTypeEnum = "fixed"
	PercentageDiscountType DiscountTypeEnum = "percentage"
)

func (e *DiscountTypeEnum) Scan(src interface{}) error {
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
	if !DiscountTypeEnum(b).IsValid() {
		return errors.New("unsupported type")
	}
	*e = DiscountTypeEnum(b)
	return nil
}

func (e DiscountTypeEnum) Value() (driver.Value, error) {
	if !e.IsValid() {
		return nil, errors.New("value invalid DiscountTypeEnum")
	}
	return string(e), nil
}

func (e DiscountTypeEnum) IsValid() bool {
	var types = []string{
		string(FixedDiscountType),
		string(PercentageDiscountType),
	}
	for _, t := range types {
		if t == string(e) {
			return true
		}
	}
	return false
}
