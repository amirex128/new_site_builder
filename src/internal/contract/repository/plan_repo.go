package repository

import (
	common "github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/domain"
)

type IPlanRepository interface {
	GetAll(paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain.Plan], error)
	GetByID(id int64) (*domain.Plan, error)
	Create(plan *domain.Plan) error
	Update(plan *domain.Plan) error
	Delete(id int64) error
}
