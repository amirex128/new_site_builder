package repository

import (
	common "github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/domain"

	"gorm.io/gorm"
)

// TicketRepo implements ITicketRepository
type TicketRepo struct {
	database *gorm.DB
}

func NewTicketRepository(db *gorm.DB) *TicketRepo {
	return &TicketRepo{
		database: db,
	}
}

func (r *TicketRepo) GetAll(paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain.Ticket], error) {
	var tickets []domain.Ticket
	var count int64

	query := r.database.Model(&domain.Ticket{})
	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&tickets)
	if result.Error != nil {
		return nil, result.Error
	}

	return buildPaginationResponse(tickets, paginationRequestDto, count)
}

func (r *TicketRepo) GetAllBySiteID(siteID int64, paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain.Ticket], error) {
	var tickets []domain.Ticket
	var count int64

	query := r.database.Model(&domain.Ticket{}).Where("site_id = ?", siteID)
	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&tickets)
	if result.Error != nil {
		return nil, result.Error
	}

	return buildPaginationResponse(tickets, paginationRequestDto, count)
}

func (r *TicketRepo) GetAllByUserID(userID int64, paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain.Ticket], error) {
	var tickets []domain.Ticket
	var count int64

	query := r.database.Model(&domain.Ticket{}).Where("user_id = ?", userID)
	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&tickets)
	if result.Error != nil {
		return nil, result.Error
	}

	return buildPaginationResponse(tickets, paginationRequestDto, count)
}

func (r *TicketRepo) GetByID(id int64) (*domain.Ticket, error) {
	var ticket domain.Ticket
	result := r.database.First(&ticket, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &ticket, nil
}

func (r *TicketRepo) GetByIDWithRelations(id int64) (*domain.Ticket, error) {
	var ticket domain.Ticket

	// Get ticket with preloaded relations
	result := r.database.
		Preload("Comments").
		Preload("Media").
		Preload("User").
		First(&ticket, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return &ticket, nil
}

func (r *TicketRepo) Create(ticket *domain.Ticket) error {
	result := r.database.Create(ticket)
	return result.Error
}

func (r *TicketRepo) Update(ticket *domain.Ticket) error {
	result := r.database.Save(ticket)
	return result.Error
}

func (r *TicketRepo) Delete(id int64) error {
	result := r.database.Delete(&domain.Ticket{}, id)
	return result.Error
}

// CustomerTicketRepo implements ICustomerTicketRepository
type CustomerTicketRepo struct {
	database *gorm.DB
}

func NewCustomerTicketRepository(db *gorm.DB) *CustomerTicketRepo {
	return &CustomerTicketRepo{
		database: db,
	}
}

func (r *CustomerTicketRepo) GetAll(paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain.CustomerTicket], error) {
	var tickets []domain.CustomerTicket
	var count int64

	query := r.database.Model(&domain.CustomerTicket{})
	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&tickets)
	if result.Error != nil {
		return nil, result.Error
	}

	return buildPaginationResponse(tickets, paginationRequestDto, count)
}

func (r *CustomerTicketRepo) GetAllBySiteID(siteID int64, paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain.CustomerTicket], error) {
	var tickets []domain.CustomerTicket
	var count int64

	query := r.database.Model(&domain.CustomerTicket{}).Where("site_id = ?", siteID)
	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&tickets)
	if result.Error != nil {
		return nil, result.Error
	}

	return buildPaginationResponse(tickets, paginationRequestDto, count)
}

func (r *CustomerTicketRepo) GetAllByCustomerID(customerID int64, paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain.CustomerTicket], error) {
	var tickets []domain.CustomerTicket
	var count int64

	query := r.database.Model(&domain.CustomerTicket{}).Where("customer_id = ?", customerID)
	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&tickets)
	if result.Error != nil {
		return nil, result.Error
	}

	return buildPaginationResponse(tickets, paginationRequestDto, count)
}

func (r *CustomerTicketRepo) GetByID(id int64) (*domain.CustomerTicket, error) {
	var ticket domain.CustomerTicket
	result := r.database.First(&ticket, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &ticket, nil
}

func (r *CustomerTicketRepo) GetByIDWithRelations(id int64) (*domain.CustomerTicket, error) {
	var ticket domain.CustomerTicket

	// Get customer ticket with preloaded relations
	result := r.database.
		Preload("Comments").
		Preload("Media").
		Preload("Customer").
		First(&ticket, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return &ticket, nil
}

func (r *CustomerTicketRepo) Create(ticket *domain.CustomerTicket) error {
	result := r.database.Create(ticket)
	return result.Error
}

func (r *CustomerTicketRepo) Update(ticket *domain.CustomerTicket) error {
	result := r.database.Save(ticket)
	return result.Error
}

func (r *CustomerTicketRepo) Delete(id int64) error {
	result := r.database.Delete(&domain.CustomerTicket{}, id)
	return result.Error
}
