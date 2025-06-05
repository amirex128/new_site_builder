package repository

import (
	"github.com/amirex128/new_site_builder/internal/contract/common"
	"github.com/amirex128/new_site_builder/internal/domain"
)

type ICouponRepository interface {
	GetAll(paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain.Coupon], error)
	GetByProductID(productID int64) (*domain.Coupon, error)
	GetByID(id int64) (*domain.Coupon, error)
	Create(coupon *domain.Coupon) error
	Update(coupon *domain.Coupon) error
	Delete(id int64) error
	DecreaseQuantity(couponID int64) error
}
