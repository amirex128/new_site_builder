package repository

import (
	common "github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/domain"
)

type IBasketItemRepository interface {
	GetAll(paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain.BasketItem], error)
	GetAllByBasketID(basketID int64, paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain.BasketItem], error)
	GetByID(id int64) (domain.BasketItem, error)
	Create(basketItem domain.BasketItem) error
	Update(basketItem domain.BasketItem) error
	Delete(id int64) error
}
