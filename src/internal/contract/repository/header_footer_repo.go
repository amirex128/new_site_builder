package repository

import (
	common "github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/domain"
)

type IHeaderFooterRepository interface {
	GetAll(paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain.HeaderFooter], error)
	GetAllBySiteID(siteID int64, paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain.HeaderFooter], error)
	GetByID(id int64) (domain.HeaderFooter, error)
	Create(headerFooter domain.HeaderFooter) error
	Update(headerFooter domain.HeaderFooter) error
	Delete(id int64) error
}
