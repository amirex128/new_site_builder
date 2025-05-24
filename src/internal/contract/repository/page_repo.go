package repository

import (
	common "github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/domain"
)

type IPageRepository interface {
	GetAll(paginationRequestDto common.PaginationRequestDto) (common.PaginationResponseDto[domain.Page], int64, error)
	GetAllBySiteID(siteID int64, paginationRequestDto common.PaginationRequestDto) (common.PaginationResponseDto[domain.Page], int64, error)
	GetByID(id int64) (domain.Page, error)
	GetByIDAndSiteID(id, siteID int64) (domain.Page, error)
	GetBySlug(slug string, siteID int64) (domain.Page, error)
	GetByPaths(paths []string, siteID int64) ([]domain.Page, error)
	GetByIDs(ids []int64, siteID int64) ([]domain.Page, error)
	Create(page domain.Page) error
	Update(page domain.Page) error
	Delete(id int64) error
	AddMediaToPage(pageID int64, mediaIDs []int64) error
	RemoveAllMediaFromPage(pageID int64) error
}
