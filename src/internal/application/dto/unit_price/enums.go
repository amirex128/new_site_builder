package unit_price

// UnitPriceNameEnum defines unit price types
type UnitPriceNameEnum int

const (
	StorageMbCredits UnitPriceNameEnum = iota
	PageViewCredits
	FormSubmitCredits
	SiteCredits
	SmsCredits
	EmailCredits
	AiCredits
	AiImageCredits
)

// DiscountTypeEnum defines discount types
type DiscountTypeEnum int

const (
	Fixed DiscountTypeEnum = iota
	Percentage
)
