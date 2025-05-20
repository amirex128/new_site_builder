package enums

import (
	"database/sql/driver"
	"errors"
)

type ArticleSortEnum string

const (
	TitleAZSort         ArticleSortEnum = "title_a_z"
	TitleZASort         ArticleSortEnum = "title_z_a"
	RecentlyAddedSort   ArticleSortEnum = "recently_added"
	RecentlyUpdatedSort ArticleSortEnum = "recently_updated"
	MostVisitedSort     ArticleSortEnum = "most_visited"
	LeastVisitedSort    ArticleSortEnum = "least_visited"
	MostRatedSort       ArticleSortEnum = "most_rated"
	LeastRatedSort      ArticleSortEnum = "least_rated"
	MostReviewedSort    ArticleSortEnum = "most_reviewed"
	LeastReviewedSort   ArticleSortEnum = "least_reviewed"
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
		string(TitleAZSort),
		string(TitleZASort),
		string(RecentlyAddedSort),
		string(RecentlyUpdatedSort),
		string(MostVisitedSort),
		string(LeastVisitedSort),
		string(MostRatedSort),
		string(LeastRatedSort),
		string(MostReviewedSort),
		string(LeastReviewedSort),
	}

	for _, sortType := range sortTypes {
		if sortType == string(e) {
			return true
		}
	}
	return false
}
