package page_usage

import (
	"github.com/amirex128/new_site_builder/internal/domain/enums"
)

// SyncPageUsageCommand represents a command to synchronize page usage
type SyncPageUsageCommand struct {
	EntityIDs []int64             `json:"entityIds" nameFa:"شناسه های جسم" validate:"array_number=1 100 1 0 false"`
	SiteID    *int64              `json:"siteId" nameFa:"شناسه سایت" validate:"required,gt=0"`
	PageID    *int64              `json:"pageId" nameFa:"شناسه صفحه" validate:"required,gt=0"`
	Type      enums.PageUsageEnum `json:"type" nameFa:"نوع" validate:"required,enum"`
}
