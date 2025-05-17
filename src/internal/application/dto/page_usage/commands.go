package page_usage

// SyncPageUsageCommand represents a command to synchronize page usage
type SyncPageUsageCommand struct {
	EntityIDs []int64        `json:"entityIds" validate:"required,dive,gt=0" error:"required=شناسه‌های موجودیت الزامی هستند|gt=شناسه‌های موجودیت باید بزرگتر از 0 باشند"`
	SiteID    *int64         `json:"siteId" validate:"required,gt=0" error:"required=شناسه سایت الزامی است|gt=شناسه سایت باید بزرگتر از 0 باشد"`
	PageID    *int64         `json:"pageId" validate:"required,gt=0" error:"required=شناسه صفحه الزامی است|gt=شناسه صفحه باید بزرگتر از 0 باشد"`
	Type      *PageUsageEnum `json:"type" validate:"required" error:"required=نوع استفاده الزامی است"`
}
