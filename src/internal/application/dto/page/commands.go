package page

// CreatePageCommand represents a command to create a new page
type CreatePageCommand struct {
	SiteID      *int64    `json:"siteId" validate:"required,gt=0" error:"required=شناسه سایت الزامی است|gt=شناسه سایت باید بزرگتر از 0 باشد"`
	HeaderID    *int64    `json:"headerId" validate:"required,gt=0" error:"required=شناسه هدر الزامی است|gt=شناسه هدر باید بزرگتر از 0 باشد"`
	FooterID    *int64    `json:"footerId" validate:"required,gt=0" error:"required=شناسه فوتر الزامی است|gt=شناسه فوتر باید بزرگتر از 0 باشد"`
	Slug        *string   `json:"slug" validate:"required,max=200,pattern=^[a-z0-9-]+$" error:"required=نامک الزامی است|max=نامک نمی‌تواند بیشتر از 200 کاراکتر باشد|pattern=نامک فقط می‌تواند شامل حروف کوچک، اعداد و خط تیره باشد"`
	Title       *string   `json:"title" validate:"required,max=200" error:"required=عنوان الزامی است|max=عنوان نمی‌تواند بیشتر از 200 کاراکتر باشد"`
	Description *string   `json:"description,omitempty" validate:"omitempty,max=2000" error:"max=توضیحات نمی‌تواند بیشتر از 2000 کاراکتر باشد"`
	Body        *PageBody `json:"body" validate:"required" error:"required=محتوای صفحه الزامی است"`
	MediaIDs    []int64   `json:"mediaIds,omitempty" validate:"omitempty,dive,gt=0" error:"gt=شناسه‌های رسانه باید بزرگتر از 0 باشند"`
	SeoTags     []string  `json:"seoTags,omitempty" validate:"omitempty" error:""`
}

// UpdatePageCommand represents a command to update an existing page
type UpdatePageCommand struct {
	ID          *int64    `json:"id" validate:"required,gt=0" error:"required=شناسه صفحه الزامی است|gt=شناسه صفحه باید بزرگتر از 0 باشد"`
	SiteID      *int64    `json:"siteId,omitempty" validate:"omitempty,gt=0" error:"gt=شناسه سایت باید بزرگتر از 0 باشد"`
	HeaderID    *int64    `json:"headerId,omitempty" validate:"omitempty,gt=0" error:"gt=شناسه هدر باید بزرگتر از 0 باشد"`
	FooterID    *int64    `json:"footerId,omitempty" validate:"omitempty,gt=0" error:"gt=شناسه فوتر باید بزرگتر از 0 باشد"`
	Slug        *string   `json:"slug,omitempty" validate:"omitempty,max=200,pattern=^[a-z0-9-]+$" error:"max=نامک نمی‌تواند بیشتر از 200 کاراکتر باشد|pattern=نامک فقط می‌تواند شامل حروف کوچک، اعداد و خط تیره باشد"`
	Title       *string   `json:"title,omitempty" validate:"omitempty,max=200" error:"max=عنوان نمی‌تواند بیشتر از 200 کاراکتر باشد"`
	Description *string   `json:"description,omitempty" validate:"omitempty,max=2000" error:"max=توضیحات نمی‌تواند بیشتر از 2000 کاراکتر باشد"`
	Body        *PageBody `json:"body,omitempty" validate:"omitempty" error:""`
	MediaIDs    []int64   `json:"mediaIds,omitempty" validate:"omitempty,dive,gt=0" error:"gt=شناسه‌های رسانه باید بزرگتر از 0 باشند"`
	SeoTags     []string  `json:"seoTags,omitempty" validate:"omitempty" error:""`
}

// DeletePageCommand represents a command to delete a page
type DeletePageCommand struct {
	ID *int64 `json:"id" validate:"required,gt=0" error:"required=شناسه صفحه الزامی است|gt=شناسه صفحه باید بزرگتر از 0 باشد"`
}
