package header_footer

import "github.com/amirex128/new_site_builder/src/internal/application/dto/common_dto"

// HeaderFooterColumn corresponds to the C# HeaderFooterColumn class
type HeaderFooterColumn struct {
	Id         string                  `json:"id"`
	Style      interface{}             `json:"style"`
	Props      interface{}             `json:"props"`
	Components []HeaderFooterComponent `json:"components"`
}

// HeaderFooterComponent corresponds to the C# HeaderFooterComponent class
type HeaderFooterComponent struct {
	Name    string                    `json:"name"`
	Filters map[string][]string       `json:"filters"`
	Sort    string                    `json:"sort"`
	Title   string                    `json:"title"`
	Props   interface{}               `json:"props"`
	Data    common_dto.PaginationJson `json:"data"`
}

// HeaderFooterBody corresponds to the C# HeaderFooterBody class
type HeaderFooterBody struct {
	Id    string            `json:"id"`
	Props interface{}       `json:"props"`
	Rows  []HeaderFooterRow `json:"rows"`
}

// HeaderFooterRow corresponds to the C# HeaderFooterRow class
type HeaderFooterRow struct {
	Id      string               `json:"id"`
	Props   interface{}          `json:"props"`
	Columns []HeaderFooterColumn `json:"columns"`
}
