package repository

import (
	common_contract "go-boilerplate/src/internal/contract/common"
	"go-boilerplate/src/internal/domain"
)

type IOrderRepository interface {
	GetAll(paginationRequestDto common_contract.PaginationRequestDto) ([]domain.Order, int64, error)
	GetAllBySiteID(siteID int64, paginationRequestDto common_contract.PaginationRequestDto) ([]domain.Order, int64, error)
	GetAllByCustomerID(customerID int64, paginationRequestDto common_contract.PaginationRequestDto) ([]domain.Order, int64, error)
	GetByID(id int64) (domain.Order, error)
	GetByOrderNumber(orderNumber string) (domain.Order, error)
	Create(order domain.Order) error
	Update(order domain.Order) error
	Delete(id int64) error
}
