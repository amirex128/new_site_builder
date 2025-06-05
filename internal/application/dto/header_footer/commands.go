package header_footer

// CreateHeaderFooterCommand represents a command to create a new header/footer
type CreateHeaderFooterCommand struct {
	SiteID *int64                `json:"siteId" nameFa:"شناسه سایت" validate:"required,gt=0"`
	Title  *string               `json:"title" nameFa:"عنوان" validate:"required_text=1 200"`
	IsMain *bool                 `json:"isMain" nameFa:"آیا اصلی است" validate:"required_bool"`
	Body   *HeaderFooterBody     `json:"body" nameFa:"محتوا" validate:"required"`
	Type   *HeaderFooterTypeEnum `json:"type" nameFa:"نوع" validate:"required,enum"`
}

// UpdateHeaderFooterCommand represents a command to update an existing header/footer
type UpdateHeaderFooterCommand struct {
	ID     *int64                `json:"id" nameFa:"شناسه" validate:"required,gt=0"`
	SiteID *int64                `json:"siteId" nameFa:"شناسه سایت" validate:"required,gt=0"`
	Title  *string               `json:"title,omitempty" nameFa:"عنوان" validate:"optional_text=1 200"`
	IsMain *bool                 `json:"isMain" nameFa:"آیا اصلی است" validate:"required_bool"`
	Body   *HeaderFooterBody     `json:"body,omitempty" nameFa:"محتوا" validate:"omitempty"`
	Type   *HeaderFooterTypeEnum `json:"type" nameFa:"نوع" validate:"required,enum"`
}

// DeleteHeaderFooterCommand represents a command to delete a header/footer
type DeleteHeaderFooterCommand struct {
	ID *int64 `json:"id" nameFa:"شناسه" validate:"required,gt=0"`
}
