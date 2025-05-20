package product_category

// CreateCategoryCommand represents a command to create a new article product_category
type CreateCategoryCommand struct {
	Name             *string  `json:"name" validate:"required_text=1 200"`
	ParentCategoryID *int64   `json:"parentCategoryId,omitempty" validate:"omitempty"`
	Order            *int     `json:"order" validate:"required,gt=0"`
	SiteID           *int64   `json:"siteId" validate:"required"`
	Description      *string  `json:"description,omitempty" validate:"optional_text=1 2000"`
	SeoTags          []string `json:"seoTags,omitempty" validate:"array_string_optional=1 100"`
	Slug             *string  `json:"slug" validate:"required,slug"`
	MediaIDs         []int64  `json:"mediaIds,omitempty" validate:"array_number_optional=0 100 1 0 false"`
}

// UpdateCategoryCommand represents a command to update an existing article product_category
type UpdateCategoryCommand struct {
	ID               *int64   `json:"id" validate:"required"`
	SiteID           *int64   `json:"siteId" validate:"required"`
	Name             *string  `json:"name" validate:"required_text=1 200"`
	ParentCategoryID *int64   `json:"parentCategoryId,omitempty" validate:"omitempty"`
	Order            *int     `json:"order" validate:"required"`
	Description      *string  `json:"description,omitempty" validate:"optional_text=1 2000"`
	SeoTags          []string `json:"seoTags,omitempty" validate:"array_string_optional=1 100"`
	Slug             *string  `json:"slug,omitempty" validate:"slug_optional"`
	MediaIDs         []int64  `json:"mediaIds,omitempty" validate:"array_number_optional=0 100 1 0 false"`
}

// DeleteCategoryCommand represents a command to delete a article product_category
type DeleteCategoryCommand struct {
	ID *int64 `json:"id" validate:"required,gt=0"`
}
