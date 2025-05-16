package repository

import (
	common_contract "go-boilerplate/src/internal/contract/common"
	"go-boilerplate/src/internal/domain"
)

type ITicketRepository interface {
	GetAll(paginationRequestDto common_contract.PaginationRequestDto) ([]domain.Ticket, int64, error)
	GetAllBySiteID(siteID int64, paginationRequestDto common_contract.PaginationRequestDto) ([]domain.Ticket, int64, error)
	GetAllByUserID(userID int64, paginationRequestDto common_contract.PaginationRequestDto) ([]domain.Ticket, int64, error)
	GetByID(id int64) (domain.Ticket, error)
	Create(ticket domain.Ticket) error
	Update(ticket domain.Ticket) error
	Delete(id int64) error
}

type ICustomerTicketRepository interface {
	GetAll(paginationRequestDto common_contract.PaginationRequestDto) ([]domain.CustomerTicket, int64, error)
	GetAllBySiteID(siteID int64, paginationRequestDto common_contract.PaginationRequestDto) ([]domain.CustomerTicket, int64, error)
	GetAllByCustomerID(customerID int64, paginationRequestDto common_contract.PaginationRequestDto) ([]domain.CustomerTicket, int64, error)
	GetByID(id int64) (domain.CustomerTicket, error)
	Create(ticket domain.CustomerTicket) error
	Update(ticket domain.CustomerTicket) error
	Delete(id int64) error
}
