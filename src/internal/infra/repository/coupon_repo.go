package repository

import (
	common "github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/domain"

	"time"

	"gorm.io/gorm"
)

type CouponRepo struct {
	database *gorm.DB
}

func NewCouponRepository(db *gorm.DB) *CouponRepo {
	return &CouponRepo{
		database: db,
	}
}

func (r *CouponRepo) GetAll(paginationRequestDto common.PaginationRequestDto) ([]domain.Coupon, int64, error) {
	var coupons []domain.Coupon
	var count int64

	query := r.database.Model(&domain.Coupon{})
	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&coupons)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return coupons, count, nil
}

func (r *CouponRepo) GetByProductID(productID int64) (domain.Coupon, error) {
	var coupon domain.Coupon
	result := r.database.Where("product_id = ?", productID).First(&coupon)
	if result.Error != nil {
		return coupon, result.Error
	}
	return coupon, nil
}

func (r *CouponRepo) GetByID(id int64) (domain.Coupon, error) {
	var coupon domain.Coupon
	result := r.database.First(&coupon, id)
	if result.Error != nil {
		return coupon, result.Error
	}
	return coupon, nil
}

func (r *CouponRepo) Create(coupon domain.Coupon) error {
	result := r.database.Create(&coupon)
	return result.Error
}

func (r *CouponRepo) Update(coupon domain.Coupon) error {
	result := r.database.Save(&coupon)
	return result.Error
}

func (r *CouponRepo) Delete(id int64) error {
	result := r.database.Delete(&domain.Coupon{}, id)
	return result.Error
}

func (r *CouponRepo) DecreaseQuantity(couponID int64) error {
	// Get the current coupon
	var coupon domain.Coupon
	if err := r.database.First(&coupon, couponID).Error; err != nil {
		return err
	}

	// Check if there are available quantities
	if coupon.Quantity <= 0 {
		return gorm.ErrInvalidData
	}

	// Decrease the quantity
	coupon.Quantity -= 1
	coupon.UpdatedAt = time.Now()

	// Update the coupon
	return r.database.Save(&coupon).Error
}
