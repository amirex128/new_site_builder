package product_category

// CreateCategoryCommand represents a command to create a new product product_category
type CreateCategoryCommand struct {
	Name             *string  `json:"name" validate:"required,max=200" error:"required=نام الزامی است|max=نام نمی‌تواند بیشتر از 200 کاراکتر باشد"`
	ParentCategoryID *int64   `json:"parentCategoryId,omitempty" validate:"omitempty" error:""`
	Order            *int     `json:"order" validate:"required,gt=0" error:"required=ترتیب الزامی است|gt=ترتیب باید بزرگتر از 0 باشد"`
	SiteID           *int64   `json:"siteId" validate:"required" error:"required=سایت الزامی است"`
	Description      *string  `json:"description,omitempty" validate:"omitempty,max=2000" error:"max=توضیحات نمی‌تواند بیشتر از 2000 کاراکتر باشد"`
	SeoTags          []string `json:"seoTags,omitempty" validate:"omitempty,max=10,dive,max=100" error:"max=حداکثر 10 برچسب SEO مجاز است|dive.max=هر برچسب SEO نباید بیشتر از 100 کاراکتر باشد"`
	Slug             *string  `json:"slug" validate:"required" error:"required=نامک الزامی است"`
	MediaIDs         []int64  `json:"mediaIds,omitempty" validate:"omitempty" error:""`
}

// UpdateCategoryCommand represents a command to update an existing product product_category
type UpdateCategoryCommand struct {
	ID               *int64   `json:"id" validate:"required" error:"required=دسته‌بندی الزامی است"`
	SiteID           *int64   `json:"siteId" validate:"required" error:"required=سایت الزامی است"`
	Name             *string  `json:"name" validate:"required,max=200" error:"required=نام الزامی است|max=نام نمی‌تواند بیشتر از 200 کاراکتر باشد"`
	ParentCategoryID *int64   `json:"parentCategoryId,omitempty" validate:"omitempty" error:""`
	Order            *int     `json:"order" validate:"required" error:"required=ترتیب الزامی است"`
	Description      *string  `json:"description,omitempty" validate:"omitempty,max=2000" error:"max=توضیحات نمی‌تواند بیشتر از 2000 کاراکتر باشد"`
	SeoTags          []string `json:"seoTags,omitempty" validate:"omitempty" error:""`
	Slug             *string  `json:"slug,omitempty" validate:"omitempty" error:""`
	MediaIDs         []int64  `json:"mediaIds,omitempty" validate:"omitempty" error:""`
}

// DeleteCategoryCommand represents a command to delete a product product_category
type DeleteCategoryCommand struct {
	ID *int64 `json:"id" validate:"required,gt=0" error:"required=شناسه الزامی است|gt=شناسه باید بزرگتر از 0 باشد"`
}
