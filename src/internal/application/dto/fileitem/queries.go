package fileitem

// GetByIdsQuery for retrieving multiple file items by IDs
type GetByIdsQuery struct {
	IdsOrder      []MediaUrlRequest `json:"idsOrder" form:"idsOrder" nameFa:"شناسه ها" validate:"required"`
	IsTemporary   *bool             `json:"isTemporary" form:"isTemporary" nameFa:"آیا موقت است" validate:"required"`
	ExpireMinutes *int              `json:"expireMinutes,omitempty" form:"expireMinutes" nameFa:"مدت زمان اعتبار" validate:"required_if=IsTemporary true"`
}

// GetDeletedTreeDirectoryQuery for retrieving a tree of deleted directories
type GetDeletedTreeDirectoryQuery struct {
	// No fields needed
}

// GetDownloadFileItemByIdQuery for downloading a file item by ID
type GetDownloadFileItemByIdQuery struct {
	ID *int64 `json:"id" form:"id" nameFa:"شناسه" validate:"required"`
}

// GetTreeDirectoryQuery for retrieving a directory tree
type GetTreeDirectoryQuery struct {
	ParentFileItemID *int64 `json:"parentFileItemId,omitempty" form:"parentFileItemId" nameFa:"شناسه والد" validate:"optional"`
}
