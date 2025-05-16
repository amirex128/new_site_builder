package repository

import (
	common_contract "go-boilerplate/src/internal/contract/common"
	"go-boilerplate/src/internal/domain"
)

type IBasketRepository interface {
	GetAll(paginationRequestDto common_contract.PaginationRequestDto) ([]domain.Basket, int64, error)
	GetAllBySiteID(siteID int64, paginationRequestDto common_contract.PaginationRequestDto) ([]domain.Basket, int64, error)
	GetAllByCustomerID(customerID int64, paginationRequestDto common_contract.PaginationRequestDto) ([]domain.Basket, int64, error)
	GetByID(id int64) (domain.Basket, error)
	GetActiveBasketByCustomerID(customerID int64) (domain.Basket, error)
	Create(basket domain.Basket) error
	Update(basket domain.Basket) error
	Delete(id int64) error
}
