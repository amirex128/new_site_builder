package page_usage

// FindPageUsagesQuery represents a query to find page usages
type FindPageUsagesQuery struct {
	EntityIDs []int64        `json:"entityIds" form:"entityIds" validate:"required,dive,gt=0"`
	SiteID    *int64         `json:"siteId" form:"siteId" validate:"required,gt=0"`
	Type      *PageUsageEnum `json:"type" form:"type" validate:"required"`
	PageID    *int64         `json:"pageId" form:"pageId" validate:"required,gt=0"`
	UserID    *int64         `json:"userId" form:"userId" validate:"omitempty,gt=0"`
}
