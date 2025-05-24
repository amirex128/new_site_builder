package repository

import (
	common "github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/domain"
)

type IBasketRepository interface {
	GetAll(paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain.Basket], error)
	GetAllBySiteID(siteID int64, paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain.Basket], error)
	GetAllByCustomerID(customerID int64, paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain.Basket], error)
	GetByID(id int64) (domain.Basket, error)
	GetBasketByCustomerIDAndSiteID(customerID, siteID int64) (domain.Basket, error)
	GetBasketWithItemsByCustomerIDAndSiteID(customerID, siteID int64) (domain.Basket, error)
	Create(basket domain.Basket) error
	Update(basket domain.Basket) error
	UpsertBasket(basket domain.Basket) error
	Delete(id int64) error
	DeleteBasketItems(basketID int64) error
}
