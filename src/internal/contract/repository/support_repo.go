package repository

import (
	common "github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/domain"
)

type ITicketRepository interface {
	GetAll(paginationRequestDto common.PaginationRequestDto) ([]domain.Ticket, int64, error)
	GetAllBySiteID(siteID int64, paginationRequestDto common.PaginationRequestDto) ([]domain.Ticket, int64, error)
	GetAllByUserID(userID int64, paginationRequestDto common.PaginationRequestDto) ([]domain.Ticket, int64, error)
	GetByID(id int64) (domain.Ticket, error)
	Create(ticket domain.Ticket) error
	Update(ticket domain.Ticket) error
	Delete(id int64) error
}

type ICustomerTicketRepository interface {
	GetAll(paginationRequestDto common.PaginationRequestDto) ([]domain.CustomerTicket, int64, error)
	GetAllBySiteID(siteID int64, paginationRequestDto common.PaginationRequestDto) ([]domain.CustomerTicket, int64, error)
	GetAllByCustomerID(customerID int64, paginationRequestDto common.PaginationRequestDto) ([]domain.CustomerTicket, int64, error)
	GetByID(id int64) (domain.CustomerTicket, error)
	Create(ticket domain.CustomerTicket) error
	Update(ticket domain.CustomerTicket) error
	Delete(id int64) error
}
