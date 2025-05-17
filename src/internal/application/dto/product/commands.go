package product

// CreateProductCommand represents a command to create a new product
type CreateProductCommand struct {
	Name              *string                   `json:"name" validate:"required,max=200" error:"required=نام محصول الزامی است|max=نام محصول نمی‌تواند بیشتر از 200 کاراکتر باشد"`
	Description       *string                   `json:"description,omitempty" validate:"omitempty,max=2000" error:"max=توضیحات نمی‌تواند بیشتر از 2000 کاراکتر باشد"`
	Status            *StatusEnum               `json:"status" validate:"required" error:"required=وضعیت الزامی است"`
	Weight            *int                      `json:"weight" validate:"required,max=1000" error:"required=وزن الزامی است|max=وزن نمی‌تواند بیشتر از 1000 باشد"`
	FreeSend          *bool                     `json:"freeSend" validate:"required" error:"required=ارسال رایگان الزامی است"`
	LongDescription   *string                   `json:"longDescription,omitempty" validate:"omitempty,max=5000" error:"max=توضیحات کامل نمی‌تواند بیشتر از 5000 کاراکتر باشد"`
	Slug              *string                   `json:"slug" validate:"required" error:"required=نامک الزامی است"`
	SeoTags           []string                  `json:"seoTags,omitempty" validate:"omitempty" error:""`
	SiteID            *int64                    `json:"siteId" validate:"required" error:"required=سایت الزامی است"`
	Coupon            *CouponCommand            `json:"coupon,omitempty" validate:"omitempty,dive" error:""`
	ProductVariants   []ProductVariantCommand   `json:"productVariants" validate:"required,min=1,dive" error:"required=حداقل یک تنوع محصول الزامی است|min=حداقل یک تنوع محصول الزامی است"`
	ProductAttributes []ProductAttributeCommand `json:"productAttributes,omitempty" validate:"omitempty,dive" error:""`
	DiscountIDs       []int64                   `json:"discountIds,omitempty" validate:"omitempty" error:""`
	CategoryIDs       []int64                   `json:"categoryIds,omitempty" validate:"omitempty" error:""`
	MediaIDs          []int64                   `json:"mediaIds,omitempty" validate:"omitempty" error:""`
}

// UpdateProductCommand represents a command to update an existing product
type UpdateProductCommand struct {
	ID                *int64                    `json:"id" validate:"required" error:"required=محصول الزامی است"`
	SiteID            *int64                    `json:"siteId" validate:"required" error:"required=سایت الزامی است"`
	Name              *string                   `json:"name,omitempty" validate:"omitempty,max=200" error:"max=نام محصول نمی‌تواند بیشتر از 200 کاراکتر باشد"`
	Description       *string                   `json:"description,omitempty" validate:"omitempty,max=2000" error:"max=توضیحات نمی‌تواند بیشتر از 2000 کاراکتر باشد"`
	Status            *StatusEnum               `json:"status,omitempty" validate:"omitempty" error:""`
	Weight            *int                      `json:"weight,omitempty" validate:"omitempty" error:""`
	FreeSend          *bool                     `json:"freeSend,omitempty" validate:"omitempty" error:""`
	LongDescription   *string                   `json:"longDescription,omitempty" validate:"omitempty,max=5000" error:"max=توضیحات کامل نمی‌تواند بیشتر از 5000 کاراکتر باشد"`
	Slug              *string                   `json:"slug,omitempty" validate:"omitempty" error:""`
	SeoTags           []string                  `json:"seoTags,omitempty" validate:"omitempty" error:""`
	ProductVariants   []ProductVariantCommand   `json:"productVariants,omitempty" validate:"omitempty,dive" error:""`
	ProductAttributes []ProductAttributeCommand `json:"productAttributes,omitempty" validate:"omitempty,dive" error:""`
	Coupon            *CouponCommand            `json:"coupon,omitempty" validate:"omitempty,dive" error:""`
	DiscountIDs       []int64                   `json:"discountIds,omitempty" validate:"omitempty" error:""`
	CategoryIDs       []int64                   `json:"categoryIds,omitempty" validate:"omitempty" error:""`
	MediaIDs          []int64                   `json:"mediaIds,omitempty" validate:"omitempty" error:""`
}

// DeleteProductCommand represents a command to delete a product
type DeleteProductCommand struct {
	ID *int64 `json:"id" validate:"required" error:"required=محصول الزامی است"`
}
