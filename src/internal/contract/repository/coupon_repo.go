package repository

import (
	common "github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/domain"
)

type ICouponRepository interface {
	GetAll(paginationRequestDto common.PaginationRequestDto) ([]domain.Coupon, int64, error)
	GetByProductID(productID int64) (domain.Coupon, error)
	GetByID(id int64) (domain.Coupon, error)
	Create(coupon domain.Coupon) error
	Update(coupon domain.Coupon) error
	Delete(id int64) error
	DecreaseQuantity(couponID int64) error
}
