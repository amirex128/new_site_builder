package repository

import (
	"github.com/amirex128/new_site_builder/internal/contract/common"
	domain2 "github.com/amirex128/new_site_builder/internal/domain"
)

type ITicketRepository interface {
	GetAll(paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain2.Ticket], error)
	GetAllBySiteID(siteID int64, paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain2.Ticket], error)
	GetAllByUserID(userID int64, paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain2.Ticket], error)
	GetByID(id int64) (*domain2.Ticket, error)
	GetByIDWithRelations(id int64) (*domain2.Ticket, error)
	Create(ticket *domain2.Ticket) error
	Update(ticket *domain2.Ticket) error
	Delete(id int64) error
}

type ICommentRepository interface {
	GetAllByTicketID(ticketID int64) ([]domain2.Comment, error)
	GetByID(id int64) (*domain2.Comment, error)
	Create(comment *domain2.Comment) error
	Update(comment *domain2.Comment) error
	Delete(id int64) error
}

type ITicketMediaRepository interface {
	AddMediaToTicket(ticketID int64, mediaIDs []int64) error
	RemoveMediaFromTicket(ticketID int64, mediaIDs []int64) error
	GetMediaByTicketID(ticketID int64) ([]domain2.Media, error)
}
