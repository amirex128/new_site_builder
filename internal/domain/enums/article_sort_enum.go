package enums

import (
	"database/sql/driver"
	"errors"
)

type ArticleSortEnum string

const (
	ArticleTitleAZSort         ArticleSortEnum = "title_a_z"
	ArticleTitleZASort         ArticleSortEnum = "title_z_a"
	ArticleRecentlyAddedSort   ArticleSortEnum = "recently_added"
	ArticleRecentlyUpdatedSort ArticleSortEnum = "recently_updated"
	ArticleMostVisitedSort     ArticleSortEnum = "most_visited"
	ArticleLeastVisitedSort    ArticleSortEnum = "least_visited"
	ArticleMostRatedSort       ArticleSortEnum = "most_rated"
	ArticleLeastRatedSort      ArticleSortEnum = "least_rated"
	ArticleMostReviewedSort    ArticleSortEnum = "most_reviewed"
	ArticleLeastReviewedSort   ArticleSortEnum = "least_reviewed"
)

func (e *ArticleSortEnum) Scan(src interface{}) error {
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
	if !ArticleSortEnum(b).IsValid() {
		return errors.New("unsupported type")
	}
	*e = ArticleSortEnum(b)
	return nil
}

func (e ArticleSortEnum) Value() (driver.Value, error) {
	if !e.IsValid() {
		return nil, errors.New("value invalid ArticleSortEnum")
	}
	return string(e), nil
}

func (e ArticleSortEnum) IsValid() bool {
	var sortTypes = []string{
		string(ArticleTitleAZSort),
		string(ArticleTitleZASort),
		string(ArticleRecentlyAddedSort),
		string(ArticleRecentlyUpdatedSort),
		string(ArticleMostVisitedSort),
		string(ArticleLeastVisitedSort),
		string(ArticleMostRatedSort),
		string(ArticleLeastRatedSort),
		string(ArticleMostReviewedSort),
		string(ArticleLeastReviewedSort),
	}

	for _, sortType := range sortTypes {
		if sortType == string(e) {
			return true
		}
	}
	return false
}
