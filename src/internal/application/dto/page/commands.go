package page

// CreatePageCommand represents a command to create a new page
type CreatePageCommand struct {
	SiteID      *int64    `json:"siteId" validate:"required,gt=0"`
	HeaderID    *int64    `json:"headerId" validate:"required,gt=0"`
	FooterID    *int64    `json:"footerId" validate:"required,gt=0"`
	Slug        *string   `json:"slug" validate:"required,slug"`
	Title       *string   `json:"title" validate:"required_text=1 200"`
	Description *string   `json:"description,omitempty" validate:"optional_text=1 2000"`
	Body        *PageBody `json:"body" validate:"required"`
	MediaIDs    []int64   `json:"mediaIds,omitempty" validate:"array_number_optional=0 100 1 0 false"`
	SeoTags     []string  `json:"seoTags,omitempty" validate:"array_string_optional=1 100"`
}

// UpdatePageCommand represents a command to update an existing page
type UpdatePageCommand struct {
	ID          *int64    `json:"id" validate:"required,gt=0"`
	SiteID      *int64    `json:"siteId,omitempty" validate:"omitempty,gt=0"`
	HeaderID    *int64    `json:"headerId,omitempty" validate:"omitempty,gt=0"`
	FooterID    *int64    `json:"footerId,omitempty" validate:"omitempty,gt=0"`
	Slug        *string   `json:"slug,omitempty" validate:"slug_optional"`
	Title       *string   `json:"title,omitempty" validate:"optional_text=1 200"`
	Description *string   `json:"description,omitempty" validate:"optional_text=1 2000"`
	Body        *PageBody `json:"body,omitempty" validate:"omitempty"`
	MediaIDs    []int64   `json:"mediaIds,omitempty" validate:"array_number_optional=0 100 1 0 false"`
	SeoTags     []string  `json:"seoTags,omitempty" validate:"array_string_optional=1 100"`
}

// DeletePageCommand represents a command to delete a page
type DeletePageCommand struct {
	ID *int64 `json:"id" validate:"required,gt=0"`
}
