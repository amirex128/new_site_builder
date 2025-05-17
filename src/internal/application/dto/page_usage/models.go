package page_usage

// PageUsage represents a usage record of a page
type PageUsage struct {
	ID       *int64         `json:"id,omitempty" validate:"omitempty" error:""`
	EntityID *int64         `json:"entityId" validate:"required,gt=0" error:"required=شناسه موجودیت الزامی است|gt=شناسه موجودیت باید بزرگتر از 0 باشد"`
	SiteID   *int64         `json:"siteId" validate:"required,gt=0" error:"required=شناسه سایت الزامی است|gt=شناسه سایت باید بزرگتر از 0 باشد"`
	PageID   *int64         `json:"pageId" validate:"required,gt=0" error:"required=شناسه صفحه الزامی است|gt=شناسه صفحه باید بزرگتر از 0 باشد"`
	Type     *PageUsageEnum `json:"type" validate:"required" error:"required=نوع استفاده الزامی است"`
}
