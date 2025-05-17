package page_usage

// PageUsage represents a usage record of a page
type PageUsage struct {
	ID       *int64         `json:"id,omitempty" validate:"omitempty"`
	EntityID *int64         `json:"entityId" validate:"required,gt=0"`
	SiteID   *int64         `json:"siteId" validate:"required,gt=0"`
	PageID   *int64         `json:"pageId" validate:"required,gt=0"`
	Type     *PageUsageEnum `json:"type" validate:"required,enum"`
}
