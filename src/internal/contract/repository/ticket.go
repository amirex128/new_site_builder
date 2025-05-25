package repository

import (
	common "github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/domain"
)

type ITicketRepository interface {
	GetAll(paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain.Ticket], error)
	GetAllBySiteID(siteID int64, paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain.Ticket], error)
	GetAllByUserID(userID int64, paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain.Ticket], error)
	GetByID(id int64) (*domain.Ticket, error)
	GetByIDWithRelations(id int64) (*domain.Ticket, error)
	Create(ticket *domain.Ticket) error
	Update(ticket *domain.Ticket) error
	Delete(id int64) error
}

type ICommentRepository interface {
	GetAllByTicketID(ticketID int64) ([]domain.Comment, error)
	GetByID(id int64) (*domain.Comment, error)
	Create(comment *domain.Comment) error
	Update(comment *domain.Comment) error
	Delete(id int64) error
}

type ITicketMediaRepository interface {
	AddMediaToTicket(ticketID int64, mediaIDs []int64) error
	RemoveMediaFromTicket(ticketID int64, mediaIDs []int64) error
	GetMediaByTicketID(ticketID int64) ([]domain.Media, error)
}
