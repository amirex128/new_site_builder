package product

// DiscountTypeEnum defines types of discounts
type DiscountTypeEnum int

const (
	Percentage DiscountTypeEnum = iota
	Value
)

// ProductAttributeTypeEnum defines types of product attributes
type ProductAttributeTypeEnum int

const (
	Public ProductAttributeTypeEnum = iota
	Technical
	Other
)

// StatusEnum defines product status options
type StatusEnum int

const (
	Inactive StatusEnum = iota
	Active
)

// ProductFilterEnum defines product filter options
type ProductFilterEnum int

const (
	PriceRange ProductFilterEnum = iota
	RatingRange
	CouponRange
	SellingRange
	VisitedRange
	ReviewRange
	UpdatedRange
	AddedRange
	WeightRange
	CategoryIds
	ProductIds
	ProductAttributes
	ProductVariant
	Badges
	FreeSend
)

// ProductSortEnum defines product sorting options
type ProductSortEnum int

const (
	PriceLowToHigh ProductSortEnum = iota
	PriceHighToLow
	CouponHighToLow
	NameAZ
	NameZA
	RecentlyAdded
	RecentlyUpdated
	MostSelling
	MostVisited
	MostRated
	MostReviewed
	LeastVisited
	LeastRated
	LeastReviewed
)
