package page

// CreatePageCommand represents a command to create a new page
type CreatePageCommand struct {
	SiteID      *int64    `json:"siteId" nameFa:"شناسه سایت" validate:"required,gt=0"`
	HeaderID    *int64    `json:"headerId" nameFa:"شناسه سرتیتور" validate:"required,gt=0"`
	FooterID    *int64    `json:"footerId" nameFa:"شناسه فوتر" validate:"required,gt=0"`
	Slug        *string   `json:"slug" nameFa:"اسلاگ" validate:"required,slug"`
	Title       *string   `json:"title" nameFa:"عنوان" validate:"required_text=1 200"`
	Description *string   `json:"description,omitempty" nameFa:"توضیحات" validate:"optional_text=1 2000"`
	Body        *PageBody `json:"body" nameFa:"محتوا" validate:"required"`
	MediaIDs    []int64   `json:"mediaIds,omitempty" nameFa:"شناسه های مدیا" validate:"array_number_optional=0 100 1 0 false"`
	SeoTags     []string  `json:"seoTags,omitempty" nameFa:"تگ های SEO" validate:"array_string_optional=1 100"`
}

// UpdatePageCommand represents a command to update an existing page
type UpdatePageCommand struct {
	ID          *int64    `json:"id" nameFa:"شناسه" validate:"required,gt=0"`
	SiteID      *int64    `json:"siteId,omitempty" nameFa:"شناسه سایت" validate:"omitempty,gt=0"`
	HeaderID    *int64    `json:"headerId,omitempty" nameFa:"شناسه سرتیتور" validate:"omitempty,gt=0"`
	FooterID    *int64    `json:"footerId,omitempty" nameFa:"شناسه فوتر" validate:"omitempty,gt=0"`
	Slug        *string   `json:"slug,omitempty" nameFa:"اسلاگ" validate:"slug_optional"`
	Title       *string   `json:"title,omitempty" nameFa:"عنوان" validate:"optional_text=1 200"`
	Description *string   `json:"description,omitempty" nameFa:"توضیحات" validate:"optional_text=1 2000"`
	Body        *PageBody `json:"body,omitempty" nameFa:"محتوا" validate:"omitempty"`
	MediaIDs    []int64   `json:"mediaIds,omitempty" nameFa:"شناسه های مدیا" validate:"array_number_optional=0 100 1 0 false"`
	SeoTags     []string  `json:"seoTags,omitempty" nameFa:"تگ های SEO" validate:"array_string_optional=1 100"`
}

// DeletePageCommand represents a command to delete a page
type DeletePageCommand struct {
	ID *int64 `json:"id" nameFa:"شناسه" validate:"required,gt=0"`
}
