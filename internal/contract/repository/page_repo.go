package repository

import (
	"github.com/amirex128/new_site_builder/internal/contract/common"
	"github.com/amirex128/new_site_builder/internal/domain"
)

type IPageRepository interface {
	GetAll(paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain.Page], error)
	GetAllBySiteID(siteID int64, paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain.Page], error)
	GetByID(id int64) (*domain.Page, error)
	GetByIDAndSiteID(id, siteID int64) (*domain.Page, error)
	GetBySlug(slug string, siteID int64) (*domain.Page, error)
	GetByPaths(paths []string, siteID int64) ([]domain.Page, error)
	GetByIDs(ids []int64, siteID int64) ([]domain.Page, error)
	Create(page *domain.Page) error
	Update(page *domain.Page) error
	Delete(id int64) error
	AddMediaToPage(pageID int64, mediaIDs []int64) error
	RemoveAllMediaFromPage(pageID int64) error
}
