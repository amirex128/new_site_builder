package repository

import (
	common_contract "go-boilerplate/src/internal/contract/common"
	"go-boilerplate/src/internal/domain"
)

type IPaymentRepository interface {
	GetAll(paginationRequestDto common_contract.PaginationRequestDto) ([]domain.Payment, int64, error)
	GetAllBySiteID(siteID int64, paginationRequestDto common_contract.PaginationRequestDto) ([]domain.Payment, int64, error)
	GetAllByOrderID(orderID int64, paginationRequestDto common_contract.PaginationRequestDto) ([]domain.Payment, int64, error)
	GetAllByCustomerID(customerID int64, paginationRequestDto common_contract.PaginationRequestDto) ([]domain.Payment, int64, error)
	GetByID(id int64) (domain.Payment, error)
	GetByTrackingNumber(trackingNumber string) (domain.Payment, error)
	Create(payment domain.Payment) error
	Update(payment domain.Payment) error
	Delete(id int64) error
}
