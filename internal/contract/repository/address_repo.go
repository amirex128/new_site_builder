package repository

import (
	"github.com/amirex128/new_site_builder/internal/contract/common"
	"github.com/amirex128/new_site_builder/internal/domain"
)

type IAddressRepository interface {
	GetAll(paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain.Address], error)
	GetByID(id int64) (*domain.Address, error)
	GetAllByUserID(userID int64) ([]domain.Address, error)
	GetAllByCustomerID(customerID int64) ([]domain.Address, error)
	Create(address *domain.Address) error
	Update(address *domain.Address) error
	Delete(id int64) error
	AddAddressToUser(addressID int64, userID int64) error
	RemoveAddressFromUser(addressID int64, userID int64) error
	RemoveAllAddressesFromUser(userID int64) error
	AddAddressToCustomer(addressID int64, customerID int64) error
	RemoveAddressFromCustomer(addressID int64, customerID int64) error
	RemoveAllAddressesFromCustomer(customerID int64) error
}
