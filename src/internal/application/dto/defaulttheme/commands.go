package defaulttheme

// CreateDefaultThemeCommand represents a command to create a new default theme
type CreateDefaultThemeCommand struct {
	Name        *string `json:"name" validate:"required_text=1,200"`
	Description *string `json:"description,omitempty" validate:"optional_text=1,1000"`
	Demo        *string `json:"demo,omitempty" validate:"optional_text=1,1000"`
	MediaID     *int    `json:"mediaId" validate:"required,gt=0"`
	Pages       *string `json:"pages" validate:"required_text=1,5000"`
}

// UpdateDefaultThemeCommand represents a command to update an existing default theme
type UpdateDefaultThemeCommand struct {
	ID          *int64  `json:"id" validate:"required,gt=0"`
	Name        *string `json:"name,omitempty" validate:"optional_text=1,200"`
	Description *string `json:"description,omitempty" validate:"optional_text=1,1000"`
	Demo        *string `json:"demo,omitempty" validate:"optional_text=1,1000"`
	MediaID     *int    `json:"mediaId,omitempty" validate:"omitempty,gt=0"`
	Pages       *string `json:"pages,omitempty" validate:"optional_text=1,5000"`
}

// DeleteDefaultThemeCommand represents a command to delete a default theme
type DeleteDefaultThemeCommand struct {
	ID *int64 `json:"id" validate:"required,gt=0"`
}
