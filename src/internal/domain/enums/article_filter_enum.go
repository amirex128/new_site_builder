package enums

import (
	"database/sql/driver"
	"errors"
)

type ArticleFilterEnum string

const (
	RateRangeFilter    ArticleFilterEnum = "rate_range"
	ReviewRangeFilter  ArticleFilterEnum = "review_range"
	VisitedRangeFilter ArticleFilterEnum = "visited_range"
	AddedRangeFilter   ArticleFilterEnum = "added_range"
	UpdatedRangeFilter ArticleFilterEnum = "updated_range"
	CategoryIdsFilter  ArticleFilterEnum = "category_ids"
	ArticleIdsFilter   ArticleFilterEnum = "article_ids"
	BadgesFilter       ArticleFilterEnum = "badges"
)

func (e *ArticleFilterEnum) Scan(src interface{}) error {
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
	if !ArticleFilterEnum(b).IsValid() {
		return errors.New("unsupported type")
	}
	*e = ArticleFilterEnum(b)
	return nil
}

func (e ArticleFilterEnum) Value() (driver.Value, error) {
	if !e.IsValid() {
		return nil, errors.New("value invalid ArticleFilterEnum")
	}
	return string(e), nil
}

func (e ArticleFilterEnum) IsValid() bool {
	var filterTypes = []string{
		string(RateRangeFilter),
		string(ReviewRangeFilter),
		string(VisitedRangeFilter),
		string(AddedRangeFilter),
		string(UpdatedRangeFilter),
		string(CategoryIdsFilter),
		string(ArticleIdsFilter),
		string(BadgesFilter),
	}

	for _, filterType := range filterTypes {
		if filterType == string(e) {
			return true
		}
	}
	return false
}
