package page_usage

// FindPageUsagesQuery represents a query to find page usages
type FindPageUsagesQuery struct {
	EntityIDs []int64        `json:"entityIds" validate:"required,dive,gt=0"`
	SiteID    *int64         `json:"siteId" validate:"required,gt=0"`
	Type      *PageUsageEnum `json:"type" validate:"required"`
}
