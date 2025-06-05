package article

// Command DTOs for article CRUD operations

// CreateArticleCommand represents a command to create a new article
type CreateArticleCommand struct {
	Title       *string  `json:"title" nameFa:"عنوان" validate:"required_text=3 200"`
	Description *string  `json:"description" nameFa:"توضیحات" validate:"optional_text=0 2000"`
	SiteID      *int64   `json:"siteId" nameFa:"شناسه سایت" validate:"required"`
	Body        *string  `json:"body" nameFa:"محتوا" validate:"required_text=1 2147483647"`
	Slug        *string  `json:"slug" nameFa:"اسلاگ" validate:"slug"`
	SeoTags     []string `json:"seoTags,omitempty" nameFa:"تگ های سئو" validate:"array_string_optional=0 0"`
	MediaIDs    []int64  `json:"mediaIds,omitempty" nameFa:"شناسه های مدیا" validate:"array_number_optional=0 0 0 0 false"`
	CategoryIDs []int64  `json:"categoryIds,omitempty" nameFa:"شناسه های دسته بندی" validate:"array_number_optional=0 0 0 0 false"`
}

// UpdateArticleCommand represents a command to update an existing article
type UpdateArticleCommand struct {
	ID          *int64   `json:"id" nameFa:"شناسه" validate:"required"`
	SiteID      *int64   `json:"siteId" nameFa:"شناسه سایت" validate:"required"`
	Title       *string  `json:"title,omitempty" nameFa:"عنوان" validate:"optional_text=3 200"`
	Description *string  `json:"description,omitempty" nameFa:"توضیحات" validate:"optional_text=0 2000"`
	Body        *string  `json:"body,omitempty" nameFa:"محتوا" validate:"required_text=1 2147483647"`
	Slug        *string  `json:"slug,omitempty" nameFa:"اسلاگ" validate:"slug_optional"`
	SeoTags     []string `json:"seoTags,omitempty" nameFa:"تگ های سئو" validate:"array_string_optional=0 0"`
	MediaIDs    []int64  `json:"mediaIds,omitempty" nameFa:"شناسه های مدیا" validate:"array_number_optional=0 0 0 0 false"`
	CategoryIDs []int64  `json:"categoryIds,omitempty" nameFa:"شناسه های دسته بندی" validate:"array_number_optional=0 0 0 0 false"`
}

// DeleteArticleCommand represents a command to delete an article
type DeleteArticleCommand struct {
	ID *int64 `json:"id" nameFa:"شناسه" validate:"required"`
}
