package repository

import (
	common_contract "go-boilerplate/src/internal/contract/common"
	"go-boilerplate/src/internal/domain"
)

type ISiteRepository interface {
	GetAll(paginationRequestDto common_contract.PaginationRequestDto) ([]domain.Site, int64, error)
	GetAllByUserID(userID int64, paginationRequestDto common_contract.PaginationRequestDto) ([]domain.Site, int64, error)
	GetByID(id int64) (domain.Site, error)
	GetByDomain(domain string) (domain.Site, error)
	Create(site domain.Site) error
	Update(site domain.Site) error
	Delete(id int64) error
}
