package repository

import (
	common_contract "go-boilerplate/src/internal/contract/common"
	"go-boilerplate/src/internal/domain"
)

type IHeaderFooterRepository interface {
	GetAll(paginationRequestDto common_contract.PaginationRequestDto) ([]domain.HeaderFooter, int64, error)
	GetAllBySiteID(siteID int64, paginationRequestDto common_contract.PaginationRequestDto) ([]domain.HeaderFooter, int64, error)
	GetByID(id int64) (domain.HeaderFooter, error)
	Create(headerFooter domain.HeaderFooter) error
	Update(headerFooter domain.HeaderFooter) error
	Delete(id int64) error
}
