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
	GetByIDWithRelations(id int64) (domain.Ticket, error)
	Create(ticket domain.Ticket) error
	Update(ticket domain.Ticket) error
	Delete(id int64) error
}

type ICommentRepository interface {
	GetAllByTicketID(ticketID int64) ([]domain.Comment, error)
	GetByID(id int64) (domain.Comment, error)
	Create(comment domain.Comment) error
	Update(comment domain.Comment) error
	Delete(id int64) error
}

type ITicketMediaRepository interface {
	AddMediaToTicket(ticketID int64, mediaIDs []int64) error
	RemoveMediaFromTicket(ticketID int64, mediaIDs []int64) error
	GetMediaByTicketID(ticketID int64) ([]domain.Media, error)
}

type ICustomerTicketRepository interface {
	GetAll(paginationRequestDto common.PaginationRequestDto) ([]domain.CustomerTicket, int64, error)
	GetAllBySiteID(siteID int64, paginationRequestDto common.PaginationRequestDto) ([]domain.CustomerTicket, int64, error)
	GetAllByCustomerID(customerID int64, paginationRequestDto common.PaginationRequestDto) ([]domain.CustomerTicket, int64, error)
	GetByID(id int64) (domain.CustomerTicket, error)
	GetByIDWithRelations(id int64) (domain.CustomerTicket, error)
	Create(ticket domain.CustomerTicket) error
	Update(ticket domain.CustomerTicket) error
	Delete(id int64) error
}

type ICustomerCommentRepository interface {
	GetAllByCustomerTicketID(customerTicketID int64) ([]domain.CustomerComment, error)
	GetByID(id int64) (domain.CustomerComment, error)
	Create(comment domain.CustomerComment) error
	Update(comment domain.CustomerComment) error
	Delete(id int64) error
}

type ICustomerTicketMediaRepository interface {
	AddMediaToCustomerTicket(customerTicketID int64, mediaIDs []int64) error
	RemoveMediaFromCustomerTicket(customerTicketID int64, mediaIDs []int64) error
	GetMediaByCustomerTicketID(customerTicketID int64) ([]domain.Media, error)
}
