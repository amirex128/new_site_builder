package enums

import (
	"database/sql/driver"
	"errors"
)

type ArticleFilterEnum string

const (
	ArticleRateRangeFilter    ArticleFilterEnum = "rate_range"
	ArticleReviewRangeFilter  ArticleFilterEnum = "review_range"
	ArticleVisitedRangeFilter ArticleFilterEnum = "visited_range"
	ArticleAddedRangeFilter   ArticleFilterEnum = "added_range"
	ArticleUpdatedRangeFilter ArticleFilterEnum = "updated_range"
	ArticleCategoryIdsFilter  ArticleFilterEnum = "category_ids"
	ArticleArticleIdsFilter   ArticleFilterEnum = "article_ids"
	ArticleBadgesFilter       ArticleFilterEnum = "badges"
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
		string(ArticleRateRangeFilter),
		string(ArticleReviewRangeFilter),
		string(ArticleVisitedRangeFilter),
		string(ArticleAddedRangeFilter),
		string(ArticleUpdatedRangeFilter),
		string(ArticleCategoryIdsFilter),
		string(ArticleArticleIdsFilter),
		string(ArticleBadgesFilter),
	}

	for _, filterType := range filterTypes {
		if filterType == string(e) {
			return true
		}
	}
	return false
}
