package product

import (
	"database/sql/driver"
	"errors"
)

// DiscountTypeEnum defines types of discounts
type DiscountTypeEnum string

const (
	PercentageDiscount DiscountTypeEnum = "percentage"
	ValueDiscount      DiscountTypeEnum = "value"
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

// IsValid try to validate enum value on this type
func (e DiscountTypeEnum) IsValid() bool {
	var discountTypes = []string{
		string(PercentageDiscount),
		string(ValueDiscount),
	}

	for _, discountType := range discountTypes {
		if discountType == string(e) {
			return true
		}
	}
	return false
}

// ProductAttributeTypeEnum defines types of product attributes
type ProductAttributeTypeEnum string

const (
	PublicAttribute    ProductAttributeTypeEnum = "public"
	TechnicalAttribute ProductAttributeTypeEnum = "technical"
	OtherAttribute     ProductAttributeTypeEnum = "other"
)

func (e *ProductAttributeTypeEnum) Scan(src interface{}) error {
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
	if !ProductAttributeTypeEnum(b).IsValid() {
		return errors.New("unsupported type")
	}
	*e = ProductAttributeTypeEnum(b)
	return nil
}

func (e ProductAttributeTypeEnum) Value() (driver.Value, error) {
	if !e.IsValid() {
		return nil, errors.New("value invalid ProductAttributeTypeEnum")
	}
	return string(e), nil
}

// IsValid try to validate enum value on this type
func (e ProductAttributeTypeEnum) IsValid() bool {
	var attributeTypes = []string{
		string(PublicAttribute),
		string(TechnicalAttribute),
		string(OtherAttribute),
	}

	for _, attributeType := range attributeTypes {
		if attributeType == string(e) {
			return true
		}
	}
	return false
}

// StatusEnum defines article status options
type StatusEnum string

const (
	InactiveStatus StatusEnum = "inactive"
	ActiveStatus   StatusEnum = "active"
)

func (e *StatusEnum) Scan(src interface{}) error {
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
	if !StatusEnum(b).IsValid() {
		return errors.New("unsupported type")
	}
	*e = StatusEnum(b)
	return nil
}

func (e StatusEnum) Value() (driver.Value, error) {
	if !e.IsValid() {
		return nil, errors.New("value invalid StatusEnum")
	}
	return string(e), nil
}

// IsValid try to validate enum value on this type
func (e StatusEnum) IsValid() bool {
	var statusTypes = []string{
		string(InactiveStatus),
		string(ActiveStatus),
	}

	for _, statusType := range statusTypes {
		if statusType == string(e) {
			return true
		}
	}
	return false
}

// ProductFilterEnum defines article filter options
type ProductFilterEnum string

const (
	PriceRangeFilter        ProductFilterEnum = "price_range"
	RatingRangeFilter       ProductFilterEnum = "rating_range"
	CouponRangeFilter       ProductFilterEnum = "coupon_range"
	SellingRangeFilter      ProductFilterEnum = "selling_range"
	VisitedRangeFilter      ProductFilterEnum = "visited_range"
	ReviewRangeFilter       ProductFilterEnum = "review_range"
	UpdatedRangeFilter      ProductFilterEnum = "updated_range"
	AddedRangeFilter        ProductFilterEnum = "added_range"
	WeightRangeFilter       ProductFilterEnum = "weight_range"
	CategoryIdsFilter       ProductFilterEnum = "category_ids"
	ProductIdsFilter        ProductFilterEnum = "product_ids"
	ProductAttributesFilter ProductFilterEnum = "product_attributes"
	ProductVariantFilter    ProductFilterEnum = "product_variant"
	BadgesFilter            ProductFilterEnum = "badges"
	FreeSendFilter          ProductFilterEnum = "free_send"
)

func (e *ProductFilterEnum) Scan(src interface{}) error {
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
	if !ProductFilterEnum(b).IsValid() {
		return errors.New("unsupported type")
	}
	*e = ProductFilterEnum(b)
	return nil
}

func (e ProductFilterEnum) Value() (driver.Value, error) {
	if !e.IsValid() {
		return nil, errors.New("value invalid ProductFilterEnum")
	}
	return string(e), nil
}

// IsValid try to validate enum value on this type
func (e ProductFilterEnum) IsValid() bool {
	var filterTypes = []string{
		string(PriceRangeFilter),
		string(RatingRangeFilter),
		string(CouponRangeFilter),
		string(SellingRangeFilter),
		string(VisitedRangeFilter),
		string(ReviewRangeFilter),
		string(UpdatedRangeFilter),
		string(AddedRangeFilter),
		string(WeightRangeFilter),
		string(CategoryIdsFilter),
		string(ProductIdsFilter),
		string(ProductAttributesFilter),
		string(ProductVariantFilter),
		string(BadgesFilter),
		string(FreeSendFilter),
	}

	for _, filterType := range filterTypes {
		if filterType == string(e) {
			return true
		}
	}
	return false
}

// ProductSortEnum defines article sorting options
type ProductSortEnum string

const (
	PriceLowToHighSort  ProductSortEnum = "price_low_to_high"
	PriceHighToLowSort  ProductSortEnum = "price_high_to_low"
	CouponHighToLowSort ProductSortEnum = "coupon_high_to_low"
	NameAZSort          ProductSortEnum = "name_a_z"
	NameZASort          ProductSortEnum = "name_z_a"
	RecentlyAddedSort   ProductSortEnum = "recently_added"
	RecentlyUpdatedSort ProductSortEnum = "recently_updated"
	MostSellingSort     ProductSortEnum = "most_selling"
	MostVisitedSort     ProductSortEnum = "most_visited"
	MostRatedSort       ProductSortEnum = "most_rated"
	MostReviewedSort    ProductSortEnum = "most_reviewed"
	LeastVisitedSort    ProductSortEnum = "least_visited"
	LeastRatedSort      ProductSortEnum = "least_rated"
	LeastReviewedSort   ProductSortEnum = "least_reviewed"
)

func (e *ProductSortEnum) Scan(src interface{}) error {
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
	if !ProductSortEnum(b).IsValid() {
		return errors.New("unsupported type")
	}
	*e = ProductSortEnum(b)
	return nil
}

func (e ProductSortEnum) Value() (driver.Value, error) {
	if !e.IsValid() {
		return nil, errors.New("value invalid ProductSortEnum")
	}
	return string(e), nil
}

// IsValid try to validate enum value on this type
func (e ProductSortEnum) IsValid() bool {
	var sortTypes = []string{
		string(PriceLowToHighSort),
		string(PriceHighToLowSort),
		string(CouponHighToLowSort),
		string(NameAZSort),
		string(NameZASort),
		string(RecentlyAddedSort),
		string(RecentlyUpdatedSort),
		string(MostSellingSort),
		string(MostVisitedSort),
		string(MostRatedSort),
		string(MostReviewedSort),
		string(LeastVisitedSort),
		string(LeastRatedSort),
		string(LeastReviewedSort),
	}

	for _, sortType := range sortTypes {
		if sortType == string(e) {
			return true
		}
	}
	return false
}
