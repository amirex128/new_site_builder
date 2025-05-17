package page_usage

// SyncPageUsageCommand represents a command to synchronize page usage
type SyncPageUsageCommand struct {
	EntityIDs []int64        `json:"entityIds" validate:"array_number=1,100,1,0,false"`
	SiteID    *int64         `json:"siteId" validate:"required,gt=0"`
	PageID    *int64         `json:"pageId" validate:"required,gt=0"`
	Type      *PageUsageEnum `json:"type" validate:"required,enum"`
}
