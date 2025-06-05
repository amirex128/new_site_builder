package article_category

// Command DTOs for product_category CRUD operations

// CreateCategoryCommand represents a command to create a new product_category
type CreateCategoryCommand struct {
	Name             *string  `json:"name" nameFa:"نام" validate:"required_text=3 200"`
	ParentCategoryID *int64   `json:"parentCategoryId" nameFa:"شناسه دسته بندی والد" validate:"optional"`
	SiteID           *int64   `json:"siteId" nameFa:"شناسه سایت" validate:"required"`
	Order            *int     `json:"order" nameFa:"ترتیب" validate:"required"`
	Description      *string  `json:"description" nameFa:"توضیحات" validate:"optional_text=0 2000"`
	SeoTags          []string `json:"seoTags,omitempty" nameFa:"تگ های سئو" validate:"array_string_optional=0 0"`
	Slug             *string  `json:"slug" nameFa:"اسلاگ" validate:"slug"`
	MediaIDs         []int64  `json:"mediaIds,omitempty" nameFa:"شناسه های مدیا" validate:"array_number_optional=0 0 0 0 false"`
}

// UpdateCategoryCommand represents a command to update an existing product_category
type UpdateCategoryCommand struct {
	ID               *int64   `json:"id" nameFa:"شناسه" validate:"required"`
	SiteID           *int64   `json:"siteId" nameFa:"شناسه سایت" validate:"required"`
	Name             *string  `json:"name" nameFa:"نام" validate:"required_text=3 200"`
	Order            *int     `json:"order,omitempty" nameFa:"ترتیب" validate:"optional"`
	ParentCategoryID *int64   `json:"parentCategoryId,omitempty" nameFa:"شناسه دسته بندی والد" validate:"optional"`
	Description      *string  `json:"description,omitempty" nameFa:"توضیحات" validate:"optional_text=0 2000"`
	SeoTags          []string `json:"seoTags,omitempty" nameFa:"تگ های سئو" validate:"array_string_optional=0 0"`
	Slug             *string  `json:"slug,omitempty" nameFa:"اسلاگ" validate:"slug_optional"`
	MediaIDs         []int64  `json:"mediaIds,omitempty" nameFa:"شناسه های مدیا" validate:"array_number_optional=0 0 0 0 false"`
}

// DeleteCategoryCommand represents a command to delete a product_category
type DeleteCategoryCommand struct {
	ID *int64 `json:"id" nameFa:"شناسه" validate:"required"`
}
