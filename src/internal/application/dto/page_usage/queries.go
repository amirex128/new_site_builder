package page_usage

// FindPageUsagesQuery represents a query to find page usages
type FindPageUsagesQuery struct {
	EntityIDs []int64        `json:"entityIds" validate:"required,dive,gt=0" error:"required=شناسه‌های موجودیت الزامی هستند|gt=شناسه‌های موجودیت باید بزرگتر از 0 باشند"`
	SiteID    *int64         `json:"siteId" validate:"required,gt=0" error:"required=شناسه سایت الزامی است|gt=شناسه سایت باید بزرگتر از 0 باشد"`
	Type      *PageUsageEnum `json:"type" validate:"required" error:"required=نوع استفاده الزامی است"`
}
