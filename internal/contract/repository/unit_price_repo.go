package repository

import (
	"github.com/amirex128/new_site_builder/internal/contract/common"
	"github.com/amirex128/new_site_builder/internal/domain"
)

type IUnitPriceRepository interface {
	GetAll(paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain.UnitPrice], error)
	GetByID(id int64) (*domain.UnitPrice, error)
	GetByName(name string) (*domain.UnitPrice, error)
	Create(unitPrice *domain.UnitPrice) error
	Update(unitPrice *domain.UnitPrice) error
	Delete(id int64) error
}
