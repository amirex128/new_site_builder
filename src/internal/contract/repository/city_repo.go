package repository

import (
	common "github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/domain"
)

type ICityRepository interface {
	GetAll(paginationRequestDto common.PaginationRequestDto) (common.PaginationResponseDto[domain.City], int64, error)
	GetByID(id int64) (domain.City, error)
	GetByName(name string) (domain.City, error)
	Create(city domain.City) error
	Update(city domain.City) error
	Delete(id int64) error
}
