package repository

import (
	common "github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/domain"
)

type IReturnItemRepository interface {
	GetAll(paginationRequestDto common.PaginationRequestDto) ([]domain.ReturnItem, int64, error)
	GetAllByOrderItemID(orderItemID int64, paginationRequestDto common.PaginationRequestDto) ([]domain.ReturnItem, int64, error)
	GetAllByCustomerID(customerID int64, paginationRequestDto common.PaginationRequestDto) ([]domain.ReturnItem, int64, error)
	GetByID(id int64) (domain.ReturnItem, error)
	Create(returnItem domain.ReturnItem) error
	Update(returnItem domain.ReturnItem) error
	Delete(id int64) error
}
