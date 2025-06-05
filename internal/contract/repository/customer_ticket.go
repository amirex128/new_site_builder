package repository

import (
	"github.com/amirex128/new_site_builder/internal/contract/common"
	domain2 "github.com/amirex128/new_site_builder/internal/domain"
)

type ICustomerTicketRepository interface {
	GetAll(paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain2.CustomerTicket], error)
	GetAllBySiteID(siteID int64, paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain2.CustomerTicket], error)
	GetAllByCustomerID(customerID int64, paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain2.CustomerTicket], error)
	GetByID(id int64) (*domain2.CustomerTicket, error)
	GetByIDWithRelations(id int64) (*domain2.CustomerTicket, error)
	Create(ticket *domain2.CustomerTicket) error
	Update(ticket *domain2.CustomerTicket) error
	Delete(id int64) error
}

type ICustomerCommentRepository interface {
	GetAllByCustomerTicketID(customerTicketID int64) ([]domain2.CustomerComment, error)
	GetByID(id int64) (*domain2.CustomerComment, error)
	Create(comment *domain2.CustomerComment) error
	Update(comment *domain2.CustomerComment) error
	Delete(id int64) error
}

type ICustomerTicketMediaRepository interface {
	AddMediaToCustomerTicket(customerTicketID int64, mediaIDs []int64) error
	RemoveMediaFromCustomerTicket(customerTicketID int64, mediaIDs []int64) error
	GetMediaByCustomerTicketID(customerTicketID int64) ([]domain2.Media, error)
}
