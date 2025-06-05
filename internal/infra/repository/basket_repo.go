package repository

import (
	"github.com/amirex128/new_site_builder/internal/contract/common"
	domain2 "github.com/amirex128/new_site_builder/internal/domain"
	"gorm.io/gorm"
)

type BasketRepo struct {
	database *gorm.DB
}

func NewBasketRepository(db *gorm.DB) *BasketRepo {
	return &BasketRepo{
		database: db,
	}
}

func (r *BasketRepo) GetAll(paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain2.Basket], error) {
	var baskets []domain2.Basket
	var count int64

	query := r.database.Model(&domain2.Basket{})
	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&baskets)
	if result.Error != nil {
		return nil, result.Error
	}

	return buildPaginationResponse(baskets, paginationRequestDto, count)
}

func (r *BasketRepo) GetAllBySiteID(siteID int64, paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain2.Basket], error) {
	var baskets []domain2.Basket
	var count int64

	query := r.database.Model(&domain2.Basket{}).Where("site_id = ?", siteID)
	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&baskets)
	if result.Error != nil {
		return nil, result.Error
	}

	return buildPaginationResponse(baskets, paginationRequestDto, count)
}

func (r *BasketRepo) GetAllByCustomerID(customerID int64, paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain2.Basket], error) {
	var baskets []domain2.Basket
	var count int64

	query := r.database.Model(&domain2.Basket{}).Where("customer_id = ?", customerID)
	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&baskets)
	if result.Error != nil {
		return nil, result.Error
	}

	return buildPaginationResponse(baskets, paginationRequestDto, count)
}

func (r *BasketRepo) GetByID(id int64) (*domain2.Basket, error) {
	var basket domain2.Basket
	result := r.database.First(&basket, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &basket, nil
}

func (r *BasketRepo) GetBasketByCustomerIDAndSiteID(customerID, siteID int64) (*domain2.Basket, error) {
	var basket domain2.Basket
	result := r.database.Where("customer_id = ? AND site_id = ?", customerID, siteID).First(&basket)
	if result.Error != nil {
		return nil, result.Error
	}
	return &basket, nil
}

func (r *BasketRepo) GetBasketWithItemsByCustomerIDAndSiteID(customerID, siteID int64) (*domain2.Basket, error) {
	var basket domain2.Basket

	// Get the basket with items preloaded
	result := r.database.
		Preload("Items").
		Preload("Items.Product").
		Preload("Items.ProductVariant").
		Where("customer_id = ? AND site_id = ?", customerID, siteID).
		First(&basket)

	if result.Error != nil {
		return nil, result.Error
	}
	return &basket, nil
}

func (r *BasketRepo) GetActiveBasketByCustomerID(customerID int64) (domain2.Basket, error) {
	var basket domain2.Basket
	result := r.database.Where("customer_id = ? AND status = 'active'", customerID).First(&basket)
	if result.Error != nil {
		return basket, result.Error
	}
	return basket, nil
}

func (r *BasketRepo) Create(basket *domain2.Basket) error {
	result := r.database.Create(basket)
	return result.Error
}

func (r *BasketRepo) Update(basket *domain2.Basket) error {
	result := r.database.Save(basket)
	return result.Error
}

func (r *BasketRepo) UpsertBasket(basket *domain2.Basket) error {
	// Check if basket exists
	var existingBasket domain2.Basket
	result := r.database.Where("customer_id = ? AND site_id = ?", basket.CustomerID, basket.SiteID).First(&existingBasket)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// Create new basket
			return r.Create(basket)
		}
		return result.Error
	}

	// Update existing basket
	basket.ID = existingBasket.ID
	return r.Update(basket)
}

func (r *BasketRepo) Delete(id int64) error {
	result := r.database.Delete(&domain2.Basket{}, id)
	return result.Error
}

func (r *BasketRepo) DeleteBasketItems(basketID int64) error {
	result := r.database.Where("basket_id = ?", basketID).Delete(&domain2.BasketItem{})
	return result.Error
}
