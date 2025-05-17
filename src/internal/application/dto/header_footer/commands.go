package header_footer

import (
	"github.com/amirex128/new_site_builder/src/internal/application/dto/site"
)

// CreateHeaderFooterCommand represents a command to create a new header/footer
type CreateHeaderFooterCommand struct {
	SiteID *int64                     `json:"siteId" validate:"required,gt=0" error:"required=شناسه سایت الزامی است|gt=شناسه سایت باید بزرگتر از 0 باشد"`
	Title  *string                    `json:"title" validate:"required,max=200" error:"required=عنوان الزامی است|max=عنوان نمی‌تواند بیشتر از 200 کاراکتر باشد"`
	IsMain *bool                      `json:"isMain" validate:"required" error:"required=وضعیت اصلی بودن هدر/فوتر الزامی است"`
	Body   *site.HeaderFooterBody     `json:"body" validate:"required" error:"required=محتوای هدر/فوتر الزامی است"`
	Type   *site.HeaderFooterTypeEnum `json:"type" validate:"required" error:"required=نوع هدر/فوتر الزامی است"`
}

// UpdateHeaderFooterCommand represents a command to update an existing header/footer
type UpdateHeaderFooterCommand struct {
	ID     *int64                     `json:"id" validate:"required,gt=0" error:"required=شناسه هدر/فوتر الزامی است|gt=شناسه هدر/فوتر باید بزرگتر از 0 باشد"`
	SiteID *int64                     `json:"siteId" validate:"required,gt=0" error:"required=شناسه سایت الزامی است|gt=شناسه سایت باید بزرگتر از 0 باشد"`
	Title  *string                    `json:"title,omitempty" validate:"omitempty,max=200" error:"max=عنوان نمی‌تواند بیشتر از 200 کاراکتر باشد"`
	IsMain *bool                      `json:"isMain" validate:"required" error:"required=وضعیت اصلی بودن هدر/فوتر الزامی است"`
	Body   *site.HeaderFooterBody     `json:"body,omitempty" validate:"omitempty" error:""`
	Type   *site.HeaderFooterTypeEnum `json:"type" validate:"required" error:"required=نوع هدر/فوتر الزامی است"`
}

// DeleteHeaderFooterCommand represents a command to delete a header/footer
type DeleteHeaderFooterCommand struct {
	ID *int64 `json:"id" validate:"required,gt=0" error:"required=شناسه هدر/فوتر الزامی است|gt=شناسه هدر/فوتر باید بزرگتر از 0 باشد"`
}
