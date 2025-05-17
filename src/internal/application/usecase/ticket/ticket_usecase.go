package ticketusecase

import (
	"fmt"
	"strconv"
	"time"

	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	"github.com/amirex128/new_site_builder/src/internal/application/dto/ticket"
	"github.com/amirex128/new_site_builder/src/internal/contract"
	"github.com/amirex128/new_site_builder/src/internal/contract/repository"
	"github.com/amirex128/new_site_builder/src/internal/domain"
)

type TicketUsecase struct {
	logger sflogger.Logger
	repo   repository.ITicketRepository
}

func NewTicketUsecase(c contract.IContainer) *TicketUsecase {
	return &TicketUsecase{
		logger: c.GetLogger(),
		repo:   c.GetTicketRepo(),
	}
}

func (u *TicketUsecase) CreateTicketCommand(params *ticket.CreateTicketCommand) (any, error) {
	// Implementation for creating a ticket
	fmt.Println(params)

	// Create the ticket
	newTicket := domain.Ticket{
		Title:     *params.Title,
		Status:    strconv.Itoa(int(ticket.New)), // Default to New status
		Category:  strconv.Itoa(int(*params.Category)),
		Priority:  strconv.Itoa(int(*params.Priority)),
		UserID:    1, // Should come from auth context
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		IsDeleted: false,
	}

	err := u.repo.Create(newTicket)
	if err != nil {
		return nil, err
	}

	// Create the first comment
	_ = domain.Comment{
		TicketID:     newTicket.ID,
		Content:      *params.Comment.Content,
		RespondentID: *params.Comment.RespondentID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		IsDeleted:    false,
	}

	// TODO: Save the comment using a comment repository
	// TODO: Handle media attachments using TicketMedia join table

	return newTicket, nil
}

func (u *TicketUsecase) ReplayTicketCommand(params *ticket.ReplayTicketCommand) (any, error) {
	// Implementation for replying to a ticket
	fmt.Println(params)

	// Get the existing ticket
	existingTicket, err := u.repo.GetByID(*params.ID)
	if err != nil {
		return nil, err
	}

	// Update the ticket
	existingTicket.Status = strconv.Itoa(int(*params.Status))
	existingTicket.Category = strconv.Itoa(int(*params.Category))
	existingTicket.AssignedTo = params.AssignedTo
	existingTicket.Priority = strconv.Itoa(int(*params.Priority))
	existingTicket.UpdatedAt = time.Now()

	// If status is Closed, set the closed info
	if *params.Status == ticket.Closed {
		userId := *params.Comment.RespondentID // Usually the current user
		now := time.Now()
		existingTicket.ClosedBy = &userId
		existingTicket.ClosedAt = &now
	}

	err = u.repo.Update(existingTicket)
	if err != nil {
		return nil, err
	}

	// Create a new comment
	_ = domain.Comment{
		TicketID:     existingTicket.ID,
		Content:      *params.Comment.Content,
		RespondentID: *params.Comment.RespondentID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		IsDeleted:    false,
	}

	// TODO: Save the comment using a comment repository
	// TODO: Handle media attachments using TicketMedia join table

	return existingTicket, nil
}

func (u *TicketUsecase) AdminReplayTicketCommand(params *ticket.AdminReplayTicketCommand) (any, error) {
	// Implementation for admin replying to a ticket
	fmt.Println(params)

	// Get the existing ticket
	existingTicket, err := u.repo.GetByID(*params.ID)
	if err != nil {
		return nil, err
	}

	// Update the ticket
	existingTicket.Status = strconv.Itoa(int(*params.Status))
	existingTicket.Category = strconv.Itoa(int(*params.Category))
	if params.AssignedTo != nil {
		existingTicket.AssignedTo = params.AssignedTo
	}
	existingTicket.Priority = strconv.Itoa(int(*params.Priority))
	existingTicket.UpdatedAt = time.Now()

	// If status is Closed, set the closed info
	if *params.Status == ticket.Closed {
		userId := *params.Comment.RespondentID // Usually the admin user
		now := time.Now()
		existingTicket.ClosedBy = &userId
		existingTicket.ClosedAt = &now
	}

	err = u.repo.Update(existingTicket)
	if err != nil {
		return nil, err
	}

	// Create a new comment
	_ = domain.Comment{
		TicketID:     existingTicket.ID,
		Content:      *params.Comment.Content,
		RespondentID: *params.Comment.RespondentID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		IsDeleted:    false,
	}

	// TODO: Save the comment using a comment repository
	// TODO: Handle media attachments using TicketMedia join table

	return existingTicket, nil
}

func (u *TicketUsecase) GetByIdTicketQuery(params *ticket.GetByIdTicketQuery) (any, error) {
	// Implementation to get ticket by ID
	fmt.Println(params)

	result, err := u.repo.GetByID(*params.ID)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (u *TicketUsecase) GetAllTicketQuery(params *ticket.GetAllTicketQuery) (any, error) {
	// Implementation to get all tickets for the current user
	fmt.Println(params)

	// In a real implementation, get the user ID from the auth context
	userID := int64(1)

	result, count, err := u.repo.GetAllByUserID(userID, params.PaginationRequestDto)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"items": result,
		"total": count,
	}, nil
}

func (u *TicketUsecase) AdminGetAllTicketQuery(params *ticket.AdminGetAllTicketQuery) (any, error) {
	// Implementation to get all tickets for admin
	fmt.Println(params)

	result, count, err := u.repo.GetAll(params.PaginationRequestDto)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"items": result,
		"total": count,
	}, nil
}
