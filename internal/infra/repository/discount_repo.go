package repository

import (
	"github.com/amirex128/new_site_builder/internal/contract/common"
	"github.com/amirex128/new_site_builder/internal/domain"
	"time"

	"gorm.io/gorm"
)

type DiscountRepo struct {
	database *gorm.DB
}

func NewDiscountRepository(db *gorm.DB) *DiscountRepo {
	return &DiscountRepo{
		database: db,
	}
}

func (r *DiscountRepo) GetAll(paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain.Discount], error) {
	var discounts []domain.Discount
	var count int64

	query := r.database.Model(&domain.Discount{})
	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&discounts)
	if result.Error != nil {
		return nil, result.Error
	}

	return buildPaginationResponse(discounts, paginationRequestDto, count)
}

func (r *DiscountRepo) GetAllBySiteID(siteID int64, paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain.Discount], error) {
	var discounts []domain.Discount
	var count int64

	query := r.database.Model(&domain.Discount{}).Where("site_id = ?", siteID)
	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&discounts)
	if result.Error != nil {
		return nil, result.Error
	}

	return buildPaginationResponse(discounts, paginationRequestDto, count)
}

func (r *DiscountRepo) GetAllByProductID(productID int64, paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain.Discount], error) {
	var discounts []domain.Discount
	var count int64

	query := r.database.Model(&domain.Discount{}).
		Joins("JOIN product_discount ON product_discount.discount_id = discounts.id").
		Where("product_discount.product_id = ?", productID)

	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&discounts)
	if result.Error != nil {
		return nil, result.Error
	}

	return buildPaginationResponse(discounts, paginationRequestDto, count)
}

func (r *DiscountRepo) GetByID(id int64) (*domain.Discount, error) {
	var discount *domain.Discount
	result := r.database.First(&discount, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return discount, nil
}

func (r *DiscountRepo) GetByCode(code string) (*domain.Discount, error) {
	var discount *domain.Discount
	result := r.database.Where("code = ?", code).First(&discount)
	if result.Error != nil {
		return nil, result.Error
	}
	return discount, nil
}

func (r *DiscountRepo) Create(discount *domain.Discount) error {
	result := r.database.Create(discount)
	return result.Error
}

func (r *DiscountRepo) Update(discount *domain.Discount) error {
	result := r.database.Save(discount)
	return result.Error
}

func (r *DiscountRepo) Delete(id int64) error {
	result := r.database.Delete(&domain.Discount{}, id)
	return result.Error
}

func (r *DiscountRepo) DecreaseQuantity(discountID int64) error {
	// Get the current discount
	var discount domain.Discount
	if err := r.database.First(&discount, discountID).Error; err != nil {
		return err
	}

	// Check if there are available quantities
	if discount.Quantity <= 0 {
		return gorm.ErrInvalidData
	}

	// Decrease the quantity
	discount.Quantity -= 1
	discount.UpdatedAt = time.Now()

	// Update the discount
	return r.database.Save(&discount).Error
}

func (r *DiscountRepo) AddCustomerUsage(discountID int64, customerID int64) error {
	// Create a new customer discount record
	customerDiscount := domain.CustomerDiscount{
		DiscountID: discountID,
		CustomerID: customerID,
	}

	// Insert the record
	return r.database.Create(&customerDiscount).Error
}

func (r *DiscountRepo) HasCustomerUsedDiscount(discountID int64, customerID int64) (bool, error) {
	var count int64

	// Check if a record exists
	err := r.database.Model(&domain.CustomerDiscount{}).
		Where("discount_id = ? AND customer_id = ?", discountID, customerID).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
