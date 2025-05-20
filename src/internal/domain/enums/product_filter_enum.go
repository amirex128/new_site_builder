package enums

import (
	"database/sql/driver"
	"errors"
)

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
