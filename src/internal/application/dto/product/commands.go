package product

import (
	"github.com/amirex128/new_site_builder/src/internal/domain/enums"
)

// CreateProductCommand represents a command to create a new article
type CreateProductCommand struct {
	Name              *string                   `json:"name" validate:"required_text=1 200"`
	Description       *string                   `json:"description,omitempty" validate:"optional_text=1 2000"`
	Status            *enums.StatusEnum         `json:"status" validate:"required,enum"`
	Weight            *int                      `json:"weight" validate:"required,max=1000"`
	FreeSend          *bool                     `json:"freeSend" validate:"required_bool"`
	LongDescription   *string                   `json:"longDescription,omitempty" validate:"optional_text=1 5000"`
	Slug              *string                   `json:"slug" validate:"required,slug"`
	SeoTags           []string                  `json:"seoTags,omitempty" validate:"array_string_optional=1 100"`
	SiteID            *int64                    `json:"siteId" validate:"required"`
	Coupon            *CouponCommand            `json:"coupon,omitempty" validate:"omitempty,dive"`
	ProductVariants   []ProductVariantCommand   `json:"productVariants" validate:"required,min=1,dive"`
	ProductAttributes []ProductAttributeCommand `json:"productAttributes,omitempty" validate:"omitempty,dive"`
	DiscountIDs       []int64                   `json:"discountIds,omitempty" validate:"array_number_optional=0 100 1 0 false"`
	CategoryIDs       []int64                   `json:"categoryIds,omitempty" validate:"array_number_optional=0 100 1 0 false"`
	MediaIDs          []int64                   `json:"mediaIds,omitempty" validate:"array_number_optional=0 100 1 0 false"`
}

// UpdateProductCommand represents a command to update an existing article
type UpdateProductCommand struct {
	ID                *int64                    `json:"id" validate:"required"`
	SiteID            *int64                    `json:"siteId" validate:"required"`
	Name              *string                   `json:"name,omitempty" validate:"optional_text=1 200"`
	Description       *string                   `json:"description,omitempty" validate:"optional_text=1 2000"`
	Status            *enums.StatusEnum         `json:"status,omitempty" validate:"enum_optional"`
	Weight            *int                      `json:"weight,omitempty" validate:"omitempty"`
	FreeSend          *bool                     `json:"freeSend,omitempty" validate:"optional_bool"`
	LongDescription   *string                   `json:"longDescription,omitempty" validate:"optional_text=1 5000"`
	Slug              *string                   `json:"slug,omitempty" validate:"slug_optional"`
	SeoTags           []string                  `json:"seoTags,omitempty" validate:"array_string_optional=1 100"`
	ProductVariants   []ProductVariantCommand   `json:"productVariants,omitempty" validate:"omitempty,dive"`
	ProductAttributes []ProductAttributeCommand `json:"productAttributes,omitempty" validate:"omitempty,dive"`
	Coupon            *CouponCommand            `json:"coupon,omitempty" validate:"omitempty,dive"`
	DiscountIDs       []int64                   `json:"discountIds,omitempty" validate:"array_number_optional=0 100 1 0 false"`
	CategoryIDs       []int64                   `json:"categoryIds,omitempty" validate:"array_number_optional=0 100 1 0 false"`
	MediaIDs          []int64                   `json:"mediaIds,omitempty" validate:"array_number_optional=0 100 1 0 false"`
}

// DeleteProductCommand represents a command to delete a article
type DeleteProductCommand struct {
	ID *int64 `json:"id" validate:"required"`
}
