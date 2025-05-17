package page_usage

// PageUsageEnum defines usage types for pages
type PageUsageEnum int

const (
	Product PageUsageEnum = iota
	Article
	HeaderFooter
)
