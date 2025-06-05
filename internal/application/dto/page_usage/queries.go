package page_usage

import (
	"github.com/amirex128/new_site_builder/internal/domain/enums"
)

// FindPageUsagesQuery represents a query to find page usages
type FindPageUsagesQuery struct {
	EntityIDs []int64             `json:"entityIds" form:"entityIds" nameFa:"شناسه های جسم" validate:"required,dive,gt=0"`
	SiteID    *int64              `json:"siteId" form:"siteId" nameFa:"شناسه سایت" validate:"required,gt=0"`
	Type      enums.PageUsageEnum `json:"type" form:"type" nameFa:"نوع" validate:"required"`
	PageID    *int64              `json:"pageId" form:"pageId" nameFa:"شناسه صفحه" validate:"required,gt=0"`
	UserID    *int64              `json:"userId" form:"userId" nameFa:"شناسه کاربر" validate:"omitempty,gt=0"`
}
