package repository

import (
	"github.com/amirex128/new_site_builder/internal/contract/common"
	"github.com/amirex128/new_site_builder/internal/domain"
)

type ICreditRepository interface {
	GetAll(paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain.Credit], error)
	GetAllByUserID(userID int64, paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain.Credit], error)
	GetAllByCustomerID(customerID int64, paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain.Credit], error)
	GetByID(id int64) (domain.Credit, error)
	Create(credit *domain.Credit) error
	Update(credit *domain.Credit) error
	Delete(id int64) error
}
