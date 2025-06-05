package repository

import (
	"github.com/amirex128/new_site_builder/internal/contract/common"
	"github.com/amirex128/new_site_builder/internal/domain"
)

type IReturnItemRepository interface {
	GetAll(paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain.ReturnItem], error)
	GetAllByOrderItemID(orderItemID int64, paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain.ReturnItem], error)
	GetAllByCustomerID(customerID int64, paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain.ReturnItem], error)
	GetByID(id int64) (*domain.ReturnItem, error)
	Create(returnItem *domain.ReturnItem) error
	Update(returnItem *domain.ReturnItem) error
	Delete(id int64) error
}
