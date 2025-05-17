package page

// PageComponent represents a component within a page column
type PageComponent struct {
	Name    *string                 `json:"name,omitempty" validate:"omitempty" error:""`
	Filters *map[string][]string    `json:"filters,omitempty" validate:"omitempty" error:""`
	Sort    *string                 `json:"sort,omitempty" validate:"omitempty" error:""`
	Title   *string                 `json:"title,omitempty" validate:"omitempty" error:""`
	Props   *map[string]interface{} `json:"props,omitempty" validate:"omitempty" error:""`
	Data    *interface{}            `json:"data,omitempty" validate:"omitempty" error:""`
}

// PageColumn represents a column within a page row
type PageColumn struct {
	ID         *string                 `json:"id,omitempty" validate:"omitempty" error:""`
	Style      *map[string]interface{} `json:"style,omitempty" validate:"omitempty" error:""`
	Props      *map[string]interface{} `json:"props,omitempty" validate:"omitempty" error:""`
	Components []*PageComponent        `json:"components,omitempty" validate:"omitempty" error:""`
}

// PageRow represents a row within a page body
type PageRow struct {
	ID      *string                 `json:"id,omitempty" validate:"omitempty" error:""`
	Props   *map[string]interface{} `json:"props,omitempty" validate:"omitempty" error:""`
	Columns []*PageColumn           `json:"columns,omitempty" validate:"omitempty" error:""`
}

// PageBody represents the body content of a page
type PageBody struct {
	ID    *string                 `json:"id,omitempty" validate:"omitempty" error:""`
	Props *map[string]interface{} `json:"props,omitempty" validate:"omitempty" error:""`
	Rows  []*PageRow              `json:"rows,omitempty" validate:"omitempty" error:""`
}
