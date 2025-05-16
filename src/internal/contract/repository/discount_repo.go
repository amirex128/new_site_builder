package repository

import (
	common_contract "go-boilerplate/src/internal/contract/common"
	"go-boilerplate/src/internal/domain"
)

type IDiscountRepository interface {
	GetAll(paginationRequestDto common_contract.PaginationRequestDto) ([]domain.Discount, int64, error)
	GetAllBySiteID(siteID int64, paginationRequestDto common_contract.PaginationRequestDto) ([]domain.Discount, int64, error)
	GetAllByProductID(productID int64, paginationRequestDto common_contract.PaginationRequestDto) ([]domain.Discount, int64, error)
	GetByID(id int64) (domain.Discount, error)
	GetByCode(code string) (domain.Discount, error)
	Create(discount domain.Discount) error
	Update(discount domain.Discount) error
	Delete(id int64) error
}
