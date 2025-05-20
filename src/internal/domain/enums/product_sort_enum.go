package enums

import (
	"database/sql/driver"
	"errors"
)

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
