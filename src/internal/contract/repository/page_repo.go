package repository

import (
	common_contract "go-boilerplate/src/internal/contract/common"
	"go-boilerplate/src/internal/domain"
)

type IPageRepository interface {
	GetAll(paginationRequestDto common_contract.PaginationRequestDto) ([]domain.Page, int64, error)
	GetAllBySiteID(siteID int64, paginationRequestDto common_contract.PaginationRequestDto) ([]domain.Page, int64, error)
	GetByID(id int64) (domain.Page, error)
	GetBySlug(slug string, siteID int64) (domain.Page, error)
	Create(page domain.Page) error
	Update(page domain.Page) error
	Delete(id int64) error
}
