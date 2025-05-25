package repository

import (
	"github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/domain"
)

type ICustomerTicketRepository interface {
	GetAll(paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain.CustomerTicket], error)
	GetAllBySiteID(siteID int64, paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain.CustomerTicket], error)
	GetAllByCustomerID(customerID int64, paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain.CustomerTicket], error)
	GetByID(id int64) (*domain.CustomerTicket, error)
	GetByIDWithRelations(id int64) (*domain.CustomerTicket, error)
	Create(ticket *domain.CustomerTicket) error
	Update(ticket *domain.CustomerTicket) error
	Delete(id int64) error
}

type ICustomerCommentRepository interface {
	GetAllByCustomerTicketID(customerTicketID int64) ([]domain.CustomerComment, error)
	GetByID(id int64) (*domain.CustomerComment, error)
	Create(comment *domain.CustomerComment) error
	Update(comment *domain.CustomerComment) error
	Delete(id int64) error
}

type ICustomerTicketMediaRepository interface {
	AddMediaToCustomerTicket(customerTicketID int64, mediaIDs []int64) error
	RemoveMediaFromCustomerTicket(customerTicketID int64, mediaIDs []int64) error
	GetMediaByCustomerTicketID(customerTicketID int64) ([]domain.Media, error)
}
