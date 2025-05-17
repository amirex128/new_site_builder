package customerticketusecase

import (
	"fmt"
	"strconv"
	"time"

	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	"github.com/amirex128/new_site_builder/src/internal/application/dto/customer_ticket"
	"github.com/amirex128/new_site_builder/src/internal/contract"
	"github.com/amirex128/new_site_builder/src/internal/contract/repository"
	"github.com/amirex128/new_site_builder/src/internal/domain"
)

type CustomerTicketUsecase struct {
	logger sflogger.Logger
	repo   repository.ICustomerTicketRepository
}

func NewCustomerTicketUsecase(c contract.IContainer) *CustomerTicketUsecase {
	return &CustomerTicketUsecase{
		logger: c.GetLogger(),
		repo:   c.GetCustomerTicketRepo(),
	}
}

func (u *CustomerTicketUsecase) CreateCustomerTicketCommand(params *customer_ticket.CreateCustomerTicketCommand) (any, error) {
	// Implementation for creating a customer ticket
	fmt.Println(params)

	// Create the customer ticket
	newTicket := domain.CustomerTicket{
		Title:      *params.Title,
		CustomerID: *params.OwnerUserID,
		Status:     strconv.Itoa(int(customer_ticket.New)), // Default to New status
		Category:   strconv.Itoa(int(*params.Category)),
		Priority:   strconv.Itoa(int(*params.Priority)),
		UserID:     1, // Should come from auth context
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		IsDeleted:  false,
	}

	err := u.repo.Create(newTicket)
	if err != nil {
		return nil, err
	}

	// Create the first comment
	_ = domain.CustomerComment{
		CustomerTicketID: newTicket.ID,
		Content:          *params.Comment.Content,
		RespondentID:     *params.Comment.RespondentID,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		IsDeleted:        false,
	}

	// TODO: Save the comment using a comment repository
	// TODO: Handle media attachments using CustomerTicketMedia join table

	return newTicket, nil
}

func (u *CustomerTicketUsecase) ReplayCustomerTicketCommand(params *customer_ticket.ReplayCustomerTicketCommand) (any, error) {
	// Implementation for replying to a customer ticket
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
	if *params.Status == customer_ticket.Closed {
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
	_ = domain.CustomerComment{
		CustomerTicketID: existingTicket.ID,
		Content:          *params.Comment.Content,
		RespondentID:     *params.Comment.RespondentID,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		IsDeleted:        false,
	}

	// TODO: Save the comment using a comment repository
	// TODO: Handle media attachments using CustomerTicketMedia join table

	return existingTicket, nil
}

func (u *CustomerTicketUsecase) AdminReplayCustomerTicketCommand(params *customer_ticket.AdminReplayCustomerTicketCommand) (any, error) {
	// Implementation for admin replying to a customer ticket
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
	if *params.Status == customer_ticket.Closed {
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
	_ = domain.CustomerComment{
		CustomerTicketID: existingTicket.ID,
		Content:          *params.Comment.Content,
		RespondentID:     *params.Comment.RespondentID,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		IsDeleted:        false,
	}

	// TODO: Save the comment using a comment repository
	// TODO: Handle media attachments using CustomerTicketMedia join table

	return existingTicket, nil
}

func (u *CustomerTicketUsecase) GetByIdCustomerTicketQuery(params *customer_ticket.GetByIdCustomerTicketQuery) (any, error) {
	// Implementation to get customer ticket by ID
	fmt.Println(params)

	result, err := u.repo.GetByID(*params.ID)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (u *CustomerTicketUsecase) GetAllCustomerTicketQuery(params *customer_ticket.GetAllCustomerTicketQuery) (any, error) {
	// Implementation to get all customer tickets for the current customer
	fmt.Println(params)

	// In a real implementation, get the customer ID from the auth context
	customerID := int64(1)

	result, count, err := u.repo.GetAllByCustomerID(customerID, params.PaginationRequestDto)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"items": result,
		"total": count,
	}, nil
}

func (u *CustomerTicketUsecase) AdminGetAllCustomerTicketQuery(params *customer_ticket.AdminGetAllCustomerTicketQuery) (any, error) {
	// Implementation to get all customer tickets for admin
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
