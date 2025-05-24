package repository

import (
	common "github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/domain"
)

type IProvinceRepository interface {
	GetAll(paginationRequestDto common.PaginationRequestDto) (common.PaginationResponseDto[domain.Province], int64, error)
	GetByID(id int64) (domain.Province, error)
	GetByName(name string) (domain.Province, error)
	Create(province domain.Province) error
	Update(province domain.Province) error
	Delete(id int64) error
}
