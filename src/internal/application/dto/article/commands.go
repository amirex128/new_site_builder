package article

// Command DTOs for article CRUD operations

// CreateArticleCommand represents a command to create a new article
type CreateArticleCommand struct {
	Title       *string  `json:"title" validate:"required_text=3 200"`
	Description *string  `json:"description" validate:"optional_text=0 2000"`
	SiteID      *int64   `json:"siteId" validate:"required"`
	Body        *string  `json:"body" validate:"required_text=1 2147483647"`
	Slug        *string  `json:"slug" validate:"slug"`
	SeoTags     []string `json:"seoTags,omitempty" validate:"array_string_optional=0 0"`
	MediaIDs    []int64  `json:"mediaIds,omitempty" validate:"array_number_optional=0 0 0 0 false"`
	CategoryIDs []int64  `json:"categoryIds,omitempty" validate:"array_number_optional=0 0 0 0 false"`
}

// UpdateArticleCommand represents a command to update an existing article
type UpdateArticleCommand struct {
	ID          *int64   `json:"id" validate:"required"`
	SiteID      *int64   `json:"siteId" validate:"required"`
	Title       *string  `json:"title,omitempty" validate:"optional_text=3 200"`
	Description *string  `json:"description,omitempty" validate:"optional_text=0 2000"`
	Body        *string  `json:"body,omitempty" validate:"required_text=1 2147483647"`
	Slug        *string  `json:"slug,omitempty" validate:"slug_optional"`
	SeoTags     []string `json:"seoTags,omitempty" validate:"array_string_optional=0 0"`
	MediaIDs    []int64  `json:"mediaIds,omitempty" validate:"array_number_optional=0 0 0 0 false"`
	CategoryIDs []int64  `json:"categoryIds,omitempty" validate:"array_number_optional=0 0 0 0 false"`
}

// DeleteArticleCommand represents a command to delete an article
type DeleteArticleCommand struct {
	ID *int64 `json:"id" validate:"required"`
}
