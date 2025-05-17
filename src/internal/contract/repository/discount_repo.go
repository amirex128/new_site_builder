package repository

import (
	common "github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/domain"
)

type IDiscountRepository interface {
	GetAll(paginationRequestDto common.PaginationRequestDto) ([]domain.Discount, int64, error)
	GetAllBySiteID(siteID int64, paginationRequestDto common.PaginationRequestDto) ([]domain.Discount, int64, error)
	GetAllByProductID(productID int64, paginationRequestDto common.PaginationRequestDto) ([]domain.Discount, int64, error)
	GetByID(id int64) (domain.Discount, error)
	GetByCode(code string) (domain.Discount, error)
	Create(discount domain.Discount) error
	Update(discount domain.Discount) error
	Delete(id int64) error
}
