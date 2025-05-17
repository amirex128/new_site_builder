package repository

import (
	common "github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/domain"
)

type ICreditRepository interface {
	GetAll(paginationRequestDto common.PaginationRequestDto) ([]domain.Credit, int64, error)
	GetAllByUserID(userID int64, paginationRequestDto common.PaginationRequestDto) ([]domain.Credit, int64, error)
	GetAllByCustomerID(customerID int64, paginationRequestDto common.PaginationRequestDto) ([]domain.Credit, int64, error)
	GetByID(id int64) (domain.Credit, error)
	Create(credit domain.Credit) error
	Update(credit domain.Credit) error
	Delete(id int64) error
}
