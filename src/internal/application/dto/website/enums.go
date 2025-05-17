package website

// ArticleFilterEnum defines article filter options
type ArticleFilterEnum int

const (
	RateRange ArticleFilterEnum = iota
	ReviewRange
	VisitedRange
	AddedRange
	UpdatedRange
	CategoryIds
	ArticleIds
	Badges
)

// ArticleSortEnum defines article sorting options
type ArticleSortEnum int

const (
	TitleAZ ArticleSortEnum = iota
	TitleZA
	RecentlyAdded
	RecentlyUpdated
	MostVisited
	LeastVisited
	MostRated
	LeastRated
	MostReviewed
	LeastReviewed
)
