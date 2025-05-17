package repository

import (
	common "github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/domain"
)

type IPageRepository interface {
	GetAll(paginationRequestDto common.PaginationRequestDto) ([]domain.Page, int64, error)
	GetAllBySiteID(siteID int64, paginationRequestDto common.PaginationRequestDto) ([]domain.Page, int64, error)
	GetByID(id int64) (domain.Page, error)
	GetBySlug(slug string, siteID int64) (domain.Page, error)
	Create(page domain.Page) error
	Update(page domain.Page) error
	Delete(id int64) error
}
