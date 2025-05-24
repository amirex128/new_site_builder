package repository

import (
	common "github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/domain"
)

type ISiteRepository interface {
	GetAll(paginationRequestDto common.PaginationRequestDto) (common.PaginationResponseDto[domain.Site], int64, error)
	GetAllByUserID(userID int64, paginationRequestDto common.PaginationRequestDto) (common.PaginationResponseDto[domain.Site], int64, error)
	GetByID(id int64) (domain.Site, error)
	GetByDomain(domain string) (domain.Site, error)
	Create(site domain.Site) error
	Update(site domain.Site) error
	Delete(id int64) error
}
