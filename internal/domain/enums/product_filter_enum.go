package enums

import (
	"database/sql/driver"
	"errors"
)

type ProductFilterEnum string

const (
	ProductPriceRangeFilter   ProductFilterEnum = "price_range"
	ProductRatingRangeFilter  ProductFilterEnum = "rating_range"
	ProductCouponRangeFilter  ProductFilterEnum = "coupon_range"
	ProductSellingRangeFilter ProductFilterEnum = "selling_range"
	ProductVisitedRangeFilter ProductFilterEnum = "visited_range"
	ProductReviewRangeFilter  ProductFilterEnum = "review_range"
	ProductUpdatedRangeFilter ProductFilterEnum = "updated_range"
	ProductAddedRangeFilter   ProductFilterEnum = "added_range"
	ProductWeightRangeFilter  ProductFilterEnum = "weight_range"
	ProductCategoryIdsFilter  ProductFilterEnum = "category_ids"
	ProductProductIdsFilter   ProductFilterEnum = "product_ids"
	ProductAttributesFilter   ProductFilterEnum = "product_attributes"
	ProductVariantFilter      ProductFilterEnum = "product_variant"
	ProductBadgesFilter       ProductFilterEnum = "badges"
	ProductFreeSendFilter     ProductFilterEnum = "free_send"
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
		string(ProductPriceRangeFilter),
		string(ProductRatingRangeFilter),
		string(ProductCouponRangeFilter),
		string(ProductSellingRangeFilter),
		string(ProductVisitedRangeFilter),
		string(ProductReviewRangeFilter),
		string(ProductUpdatedRangeFilter),
		string(ProductAddedRangeFilter),
		string(ProductWeightRangeFilter),
		string(ProductCategoryIdsFilter),
		string(ProductProductIdsFilter),
		string(ProductAttributesFilter),
		string(ProductVariantFilter),
		string(ProductBadgesFilter),
		string(ProductFreeSendFilter),
	}

	for _, filterType := range filterTypes {
		if filterType == string(e) {
			return true
		}
	}
	return false
}
