package defaulttheme

// CreateDefaultThemeCommand represents a command to create a new default theme
type CreateDefaultThemeCommand struct {
	Name        *string `json:"name" validate:"required,max=200" error:"required=نام قالب الزامی است|max=نام قالب نباید بیشتر از 200 کاراکتر باشد"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=1000" error:"max=توضیحات نباید بیشتر از 1000 کاراکتر باشد"`
	Demo        *string `json:"demo,omitempty" validate:"omitempty,max=1000" error:"max=آدرس دمو نباید بیشتر از 1000 کاراکتر باشد"`
	MediaID     *int    `json:"mediaId" validate:"required,gt=0" error:"required=شناسه رسانه الزامی است|gt=شناسه رسانه باید بزرگتر از 0 باشد"`
	Pages       *string `json:"pages" validate:"required,max=5000" error:"required=صفحات الزامی هستند|max=محتوای صفحات نباید بیشتر از 5000 کاراکتر باشد"`
}

// UpdateDefaultThemeCommand represents a command to update an existing default theme
type UpdateDefaultThemeCommand struct {
	ID          *int64  `json:"id" validate:"required,gt=0" error:"required=شناسه الزامی است|gt=شناسه قالب باید بزرگتر از 0 باشد"`
	Name        *string `json:"name,omitempty" validate:"omitempty,max=200" error:"max=نام قالب نباید بیشتر از 200 کاراکتر باشد"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=1000" error:"max=توضیحات نباید بیشتر از 1000 کاراکتر باشد"`
	Demo        *string `json:"demo,omitempty" validate:"omitempty,max=1000" error:"max=آدرس دمو نباید بیشتر از 1000 کاراکتر باشد"`
	MediaID     *int    `json:"mediaId,omitempty" validate:"omitempty,gt=0" error:"gt=شناسه رسانه باید بزرگتر از 0 باشد"`
	Pages       *string `json:"pages,omitempty" validate:"omitempty,max=5000" error:"max=محتوای صفحات نباید بیشتر از 5000 کاراکتر باشد"`
}

// DeleteDefaultThemeCommand represents a command to delete a default theme
type DeleteDefaultThemeCommand struct {
	ID *int64 `json:"id" validate:"required,gt=0" error:"required=شناسه الزامی است|gt=شناسه باید بزرگتر از 0 باشد"`
}
