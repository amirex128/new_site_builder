package site

// DomainTypeEnum defines domain types
type DomainTypeEnum int

const (
	Domain DomainTypeEnum = iota
	Subdomain
)

// SiteTypeEnum defines site types
type SiteTypeEnum int

const (
	Shop SiteTypeEnum = iota
	Blog
	Business
)

// StatusEnum defines status types
type StatusEnum int

const (
	Active StatusEnum = iota
	Inactive
	Pending
	Deleted
)
