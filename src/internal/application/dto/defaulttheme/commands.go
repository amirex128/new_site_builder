package defaulttheme

// CreateDefaultThemeCommand represents a command to create a new default theme
type CreateDefaultThemeCommand struct {
	Name        *string `json:"name" nameFa:"نام" validate:"required_text=1 200"`
	Description *string `json:"description,omitempty" nameFa:"توضیحات" validate:"optional_text=1 1000"`
	Demo        *string `json:"demo,omitempty" nameFa:"دمو" validate:"optional_text=1 1000"`
	MediaID     *int    `json:"mediaId" nameFa:"شناسه مدیا" validate:"required,gt=0"`
	Pages       *string `json:"pages" nameFa:"صفحات" validate:"required_text=1 5000"`
}

// UpdateDefaultThemeCommand represents a command to update an existing default theme
type UpdateDefaultThemeCommand struct {
	ID          *int64  `json:"id" nameFa:"شناسه" validate:"required,gt=0"`
	Name        *string `json:"name,omitempty" nameFa:"نام" validate:"optional_text=1 200"`
	Description *string `json:"description,omitempty" nameFa:"توضیحات" validate:"optional_text=1 1000"`
	Demo        *string `json:"demo,omitempty" nameFa:"دمو" validate:"optional_text=1 1000"`
	MediaID     *int    `json:"mediaId,omitempty" nameFa:"شناسه مدیا" validate:"omitempty,gt=0"`
	Pages       *string `json:"pages,omitempty" nameFa:"صفحات" validate:"optional_text=1 5000"`
}

// DeleteDefaultThemeCommand represents a command to delete a default theme
type DeleteDefaultThemeCommand struct {
	ID *int64 `json:"id" nameFa:"شناسه" validate:"required,gt=0"`
}
