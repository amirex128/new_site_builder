package repository

import (
	common "github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/domain"

	"gorm.io/gorm"
)

type BasketItemRepo struct {
	database *gorm.DB
}

func NewBasketItemRepository(db *gorm.DB) *BasketItemRepo {
	return &BasketItemRepo{
		database: db,
	}
}

func (r *BasketItemRepo) GetAll(paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain.BasketItem], error) {
	var basketItems []domain.BasketItem
	var count int64

	query := r.database.Model(&domain.BasketItem{})
	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&basketItems)
	if result.Error != nil {
		return nil, result.Error
	}

	return buildPaginationResponse(basketItems, paginationRequestDto, count)
}

func (r *BasketItemRepo) GetAllByBasketID(basketID int64, paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain.BasketItem], error) {
	var basketItems []domain.BasketItem
	var count int64

	query := r.database.Model(&domain.BasketItem{}).Where("basket_id = ?", basketID)
	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&basketItems)
	if result.Error != nil {
		return nil, result.Error
	}

	return buildPaginationResponse(basketItems, paginationRequestDto, count)
}

func (r *BasketItemRepo) GetByID(id int64) (*domain.BasketItem, error) {
	var basketItem domain.BasketItem
	result := r.database.First(&basketItem, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &basketItem, nil
}

func (r *BasketItemRepo) Create(basketItem *domain.BasketItem) error {
	result := r.database.Create(basketItem)
	return result.Error
}

func (r *BasketItemRepo) Update(basketItem *domain.BasketItem) error {
	result := r.database.Save(basketItem)
	return result.Error
}

func (r *BasketItemRepo) Delete(id int64) error {
	result := r.database.Delete(&domain.BasketItem{}, id)
	return result.Error
}
