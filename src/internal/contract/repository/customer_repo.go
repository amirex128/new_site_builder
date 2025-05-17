package repository

import (
	common "github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/domain"
)

type ICustomerRepository interface {
	GetAll(paginationRequestDto common.PaginationRequestDto) ([]domain.Customer, int64, error)
	GetAllBySiteID(siteID int64, paginationRequestDto common.PaginationRequestDto) ([]domain.Customer, int64, error)
	GetByID(id int64) (domain.Customer, error)
	GetByEmail(email string) (domain.Customer, error)
	GetByPhone(phone string) (domain.Customer, error)
	Create(customer domain.Customer) error
	Update(customer domain.Customer) error
	Delete(id int64) error
}
