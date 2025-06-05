package repository

import (
	"github.com/amirex128/new_site_builder/internal/contract/common"
	"github.com/amirex128/new_site_builder/internal/domain"
)

type IDiscountRepository interface {
	GetAll(paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain.Discount], error)
	GetAllBySiteID(siteID int64, paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain.Discount], error)
	GetAllByProductID(productID int64, paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain.Discount], error)
	GetByID(id int64) (*domain.Discount, error)
	GetByCode(code string) (*domain.Discount, error)
	Create(discount *domain.Discount) error
	Update(discount *domain.Discount) error
	Delete(id int64) error
	DecreaseQuantity(discountID int64) error
	AddCustomerUsage(discountID int64, customerID int64) error
	HasCustomerUsedDiscount(discountID int64, customerID int64) (bool, error)
}
