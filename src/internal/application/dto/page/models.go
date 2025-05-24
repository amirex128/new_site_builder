package page

import (
	"github.com/amirex128/new_site_builder/src/internal/application/utils/resp"
)

type PageComponent struct {
	Name    string              `json:"name"`
	Filters map[string][]string `json:"filters"`
	Sort    string              `json:"sort"`
	Title   string              `json:"title"`
	Props   interface{}         `json:"props"`
	Data    resp.PaginationJson `json:"data"`
}

type PageColumn struct {
	Id         string          `json:"id"`
	Style      interface{}     `json:"style"`
	Props      interface{}     `json:"props"`
	Components []PageComponent `json:"components"`
}

type PageRow struct {
	Id      string       `json:"id"`
	Props   interface{}  `json:"props"`
	Columns []PageColumn `json:"columns"`
}

type PageBody struct {
	Id    string      `json:"id"`
	Props interface{} `json:"props"`
	Rows  []PageRow   `json:"rows"`
}
