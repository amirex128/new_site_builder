package header_footer

// CreateHeaderFooterCommand represents a command to create a new header/footer
type CreateHeaderFooterCommand struct {
	SiteID *int64                `json:"siteId" validate:"required,gt=0"`
	Title  *string               `json:"title" validate:"required_text=1 200"`
	IsMain *bool                 `json:"isMain" validate:"required_bool"`
	Body   *HeaderFooterBody     `json:"body" validate:"required"`
	Type   *HeaderFooterTypeEnum `json:"type" validate:"required,enum"`
}

// UpdateHeaderFooterCommand represents a command to update an existing header/footer
type UpdateHeaderFooterCommand struct {
	ID     *int64                `json:"id" validate:"required,gt=0"`
	SiteID *int64                `json:"siteId" validate:"required,gt=0"`
	Title  *string               `json:"title,omitempty" validate:"optional_text=1 200"`
	IsMain *bool                 `json:"isMain" validate:"required_bool"`
	Body   *HeaderFooterBody     `json:"body,omitempty" validate:"omitempty"`
	Type   *HeaderFooterTypeEnum `json:"type" validate:"required,enum"`
}

// DeleteHeaderFooterCommand represents a command to delete a header/footer
type DeleteHeaderFooterCommand struct {
	ID *int64 `json:"id" validate:"required,gt=0"`
}
