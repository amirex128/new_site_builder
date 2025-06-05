package repository

import (
	"github.com/amirex128/new_site_builder/internal/contract/common"
	"github.com/amirex128/new_site_builder/internal/domain"
	"gorm.io/gorm"
)

type ReturnItemRepo struct {
	database *gorm.DB
}

func NewReturnItemRepository(db *gorm.DB) *ReturnItemRepo {
	return &ReturnItemRepo{
		database: db,
	}
}

func (r *ReturnItemRepo) GetAll(paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain.ReturnItem], error) {
	var returnItems []domain.ReturnItem
	var count int64

	query := r.database.Model(&domain.ReturnItem{})
	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&returnItems)
	if result.Error != nil {
		return nil, result.Error
	}

	return buildPaginationResponse(returnItems, paginationRequestDto, count)
}

func (r *ReturnItemRepo) GetAllByOrderItemID(orderItemID int64, paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain.ReturnItem], error) {
	var returnItems []domain.ReturnItem
	var count int64

	query := r.database.Model(&domain.ReturnItem{}).Where("order_item_id = ?", orderItemID)
	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&returnItems)
	if result.Error != nil {
		return nil, result.Error
	}

	return buildPaginationResponse(returnItems, paginationRequestDto, count)
}

func (r *ReturnItemRepo) GetAllByCustomerID(customerID int64, paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain.ReturnItem], error) {
	var returnItems []domain.ReturnItem
	var count int64

	query := r.database.Model(&domain.ReturnItem{}).Where("customer_id = ?", customerID)
	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&returnItems)
	if result.Error != nil {
		return nil, result.Error
	}

	return buildPaginationResponse(returnItems, paginationRequestDto, count)
}

func (r *ReturnItemRepo) GetByID(id int64) (*domain.ReturnItem, error) {
	var returnItem *domain.ReturnItem
	result := r.database.First(&returnItem, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return returnItem, nil
}

func (r *ReturnItemRepo) Create(returnItem *domain.ReturnItem) error {
	result := r.database.Create(returnItem)
	return result.Error
}

func (r *ReturnItemRepo) Update(returnItem *domain.ReturnItem) error {
	result := r.database.Save(returnItem)
	return result.Error
}

func (r *ReturnItemRepo) Delete(id int64) error {
	result := r.database.Delete(&domain.ReturnItem{}, id)
	return result.Error
}
