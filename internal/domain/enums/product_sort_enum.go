package enums

import (
	"database/sql/driver"
	"errors"
)

type ProductSortEnum string

const (
	ProductPriceLowToHighSort  ProductSortEnum = "price_low_to_high"
	ProductPriceHighToLowSort  ProductSortEnum = "price_high_to_low"
	ProductCouponHighToLowSort ProductSortEnum = "coupon_high_to_low"
	ProductNameAZSort          ProductSortEnum = "name_a_z"
	ProductNameZASort          ProductSortEnum = "name_z_a"
	ProductRecentlyAddedSort   ProductSortEnum = "recently_added"
	ProductRecentlyUpdatedSort ProductSortEnum = "recently_updated"
	ProductMostSellingSort     ProductSortEnum = "most_selling"
	ProductMostVisitedSort     ProductSortEnum = "most_visited"
	ProductMostRatedSort       ProductSortEnum = "most_rated"
	ProductMostReviewedSort    ProductSortEnum = "most_reviewed"
	ProductLeastVisitedSort    ProductSortEnum = "least_visited"
	ProductLeastRatedSort      ProductSortEnum = "least_rated"
	ProductLeastReviewedSort   ProductSortEnum = "least_reviewed"
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

func (e ProductSortEnum) IsValid() bool {
	var sortTypes = []string{
		string(ProductPriceLowToHighSort),
		string(ProductPriceHighToLowSort),
		string(ProductCouponHighToLowSort),
		string(ProductNameAZSort),
		string(ProductNameZASort),
		string(ProductRecentlyAddedSort),
		string(ProductRecentlyUpdatedSort),
		string(ProductMostSellingSort),
		string(ProductMostVisitedSort),
		string(ProductMostRatedSort),
		string(ProductMostReviewedSort),
		string(ProductLeastVisitedSort),
		string(ProductLeastRatedSort),
		string(ProductLeastReviewedSort),
	}

	for _, sortType := range sortTypes {
		if sortType == string(e) {
			return true
		}
	}
	return false
}
