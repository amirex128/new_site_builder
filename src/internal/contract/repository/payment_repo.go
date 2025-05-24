package repository

import (
	common "github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/domain"
)

type IPaymentRepository interface {
	GetAll(paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain.Payment], error)
	GetAllBySiteID(siteID int64, paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain.Payment], error)
	GetAllByOrderID(orderID int64, paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain.Payment], error)
	GetAllByCustomerID(customerID int64, paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain.Payment], error)
	GetByID(id int64) (domain.Payment, error)
	GetByTrackingNumber(trackingNumber string) (domain.Payment, error)
	Create(payment domain.Payment) error
	Update(payment domain.Payment) error
	Delete(id int64) error
	RequestPayment(amount int64, orderID int64, userID int64, gateway string, orderData map[string]string) (string, error)
}
