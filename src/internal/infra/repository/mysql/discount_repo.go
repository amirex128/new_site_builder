package mysql

import (
	common "github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/domain"

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

func (r *DiscountRepo) GetAll(paginationRequestDto common.PaginationRequestDto) ([]domain.Discount, int64, error) {
	var discounts []domain.Discount
	var count int64

	query := r.database.Model(&domain.Discount{})
	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&discounts)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return discounts, count, nil
}

func (r *DiscountRepo) GetAllBySiteID(siteID int64, paginationRequestDto common.PaginationRequestDto) ([]domain.Discount, int64, error) {
	var discounts []domain.Discount
	var count int64

	query := r.database.Model(&domain.Discount{}).Where("site_id = ?", siteID)
	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&discounts)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return discounts, count, nil
}

func (r *DiscountRepo) GetAllByProductID(productID int64, paginationRequestDto common.PaginationRequestDto) ([]domain.Discount, int64, error) {
	var discounts []domain.Discount
	var count int64

	// For many-to-many relationship using the join table
	query := r.database.Model(&domain.Discount{}).
		Joins("JOIN product_discount ON product_discount.discount_id = discounts.id").
		Where("product_discount.product_id = ?", productID)

	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&discounts)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return discounts, count, nil
}

func (r *DiscountRepo) GetByID(id int64) (domain.Discount, error) {
	var discount domain.Discount
	result := r.database.First(&discount, id)
	if result.Error != nil {
		return discount, result.Error
	}
	return discount, nil
}

func (r *DiscountRepo) GetByCode(code string) (domain.Discount, error) {
	var discount domain.Discount
	result := r.database.Where("code = ?", code).First(&discount)
	if result.Error != nil {
		return discount, result.Error
	}
	return discount, nil
}

func (r *DiscountRepo) Create(discount domain.Discount) error {
	result := r.database.Create(&discount)
	return result.Error
}

func (r *DiscountRepo) Update(discount domain.Discount) error {
	result := r.database.Save(&discount)
	return result.Error
}

func (r *DiscountRepo) Delete(id int64) error {
	result := r.database.Delete(&domain.Discount{}, id)
	return result.Error
}
