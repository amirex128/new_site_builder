package article_category

// Command DTOs for product_category CRUD operations

// CreateCategoryCommand represents a command to create a new product_category
type CreateCategoryCommand struct {
	Name             *string  `json:"name" validate:"required_text=3 200"`
	ParentCategoryID *int64   `json:"parentCategoryId" validate:"optional"`
	SiteID           *int64   `json:"siteId" validate:"required"`
	Order            *int     `json:"order" validate:"required"`
	Description      *string  `json:"description" validate:"optional_text=0 2000"`
	SeoTags          []string `json:"seoTags,omitempty" validate:"array_string_optional=0 0"`
	Slug             *string  `json:"slug" validate:"slug"`
	MediaIDs         []int64  `json:"mediaIds,omitempty" validate:"array_number_optional=0 0 0 0 false"`
}

// UpdateCategoryCommand represents a command to update an existing product_category
type UpdateCategoryCommand struct {
	ID               *int64   `json:"id" validate:"required"`
	SiteID           *int64   `json:"siteId" validate:"required"`
	Name             *string  `json:"name" validate:"required_text=3 200"`
	Order            *int     `json:"order,omitempty" validate:"optional"`
	ParentCategoryID *int64   `json:"parentCategoryId,omitempty" validate:"optional"`
	Description      *string  `json:"description,omitempty" validate:"optional_text=0 2000"`
	SeoTags          []string `json:"seoTags,omitempty" validate:"array_string_optional=0 0"`
	Slug             *string  `json:"slug,omitempty" validate:"slug_optional"`
	MediaIDs         []int64  `json:"mediaIds,omitempty" validate:"array_number_optional=0 0 0 0 false"`
}

// DeleteCategoryCommand represents a command to delete a product_category
type DeleteCategoryCommand struct {
	ID *int64 `json:"id" validate:"required"`
}
