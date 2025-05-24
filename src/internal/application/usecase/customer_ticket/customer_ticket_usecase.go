package customerticketusecase

import (
	"errors"
	"time"

	"github.com/amirex128/new_site_builder/src/internal/domain/enums"

	"github.com/amirex128/new_site_builder/src/internal/application/usecase"
	"github.com/amirex128/new_site_builder/src/internal/application/utils/resp"
	"github.com/amirex128/new_site_builder/src/internal/contract/service"

	"github.com/gin-gonic/gin"

	"github.com/amirex128/new_site_builder/src/internal/application/dto/customer_ticket"
	"github.com/amirex128/new_site_builder/src/internal/contract"
	"github.com/amirex128/new_site_builder/src/internal/contract/repository"
	"github.com/amirex128/new_site_builder/src/internal/domain"
	"gorm.io/gorm"
)

type CustomerTicketUsecase struct {
	*usecase.BaseUsecase
	repo                    repository.ICustomerTicketRepository
	customerCommentRepo     repository.ICustomerCommentRepository
	customerTicketMediaRepo repository.ICustomerTicketMediaRepository
	mediaRepo               repository.IMediaRepository
	authContext             func(c *gin.Context) service.IAuthService
}

func NewCustomerTicketUsecase(c contract.IContainer) *CustomerTicketUsecase {
	return &CustomerTicketUsecase{
		BaseUsecase: &usecase.BaseUsecase{
			Logger: c.GetLogger(),
		},
		repo:                    c.GetCustomerTicketRepo(),
		customerCommentRepo:     c.GetCustomerCommentRepo(),
		customerTicketMediaRepo: c.GetCustomerTicketMediaRepo(),
		mediaRepo:               c.GetMediaRepo(),
		authContext:             c.GetAuthTransientService(),
	}
}

func (u *CustomerTicketUsecase) CreateCustomerTicketCommand(params *customer_ticket.CreateCustomerTicketCommand) (*resp.Response, error) {
	// Get current user ID from auth context
	userID, err := u.authContext(u.Ctx).GetUserID()
	if err != nil {
		return nil, errors.New("خطا در احراز هویت کاربر")
	}

	// Create the customer ticket
	newTicket := domain.CustomerTicket{
		Title:      *params.Title,
		CustomerID: *params.OwnerUserID,
		Status:     enums.TicketNewStatus,
		Category:   *params.Category,
		Priority:   *params.Priority,
		UserID:     userID,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		IsDeleted:  false,
	}

	// Save ticket to repository
	err = u.repo.Create(newTicket)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}

	// Create the first comment
	comment := domain.CustomerComment{
		CustomerTicketID: newTicket.ID,
		Content:          *params.Comment.Content,
		RespondentID:     *params.Comment.RespondentID,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		IsDeleted:        false,
	}

	// Save comment
	err = u.customerCommentRepo.Create(comment)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}

	// Handle media attachments
	if len(params.MediaIDs) > 0 {
		err = u.customerTicketMediaRepo.AddMediaToCustomerTicket(newTicket.ID, params.MediaIDs)
		if err != nil {
			return nil, resp.NewError(resp.Internal, err.Error())
		}
	}

	// Get the created ticket with relations
	createdTicket, err := u.repo.GetByIDWithRelations(newTicket.ID)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}

	return resp.NewResponseData(
		resp.Created,
		resp.Data{
			"ticket": createdTicket,
		},
		"تیکت مشتری با موفقیت ایجاد شد",
	), nil
}

func (u *CustomerTicketUsecase) ReplayCustomerTicketCommand(params *customer_ticket.ReplayCustomerTicketCommand) (*resp.Response, error) {
	u.Logger.Info("ReplayCustomerTicketCommand called", map[string]interface{}{
		"id":      *params.ID,
		"status":  *params.Status,
		"content": *params.Comment.Content,
	})

	// Get customer ID from auth context
	customerID, err := u.authContext(u.Ctx).GetCustomerID()
	if err != nil {
		return nil, errors.New("خطا در احراز هویت مشتری")
	}

	// Get the existing ticket
	existingTicket, err := u.repo.GetByID(*params.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("تیکت مشتری مورد نظر یافت نشد")
		}
		return nil, err
	}

	// Check if customer owns the ticket
	if existingTicket.CustomerID != customerID {
		return nil, errors.New("شما دسترسی به این تیکت مشتری ندارید")
	}

	// Update the ticket
	existingTicket.Status = *params.Status
	existingTicket.Category = *params.Category
	existingTicket.AssignedTo = params.AssignedTo
	existingTicket.Priority = *params.Priority
	existingTicket.UpdatedAt = time.Now()

	// Set closed info if status is closed
	if *params.Status == enums.TicketClosedStatus {
		respondentID := *params.Comment.RespondentID
		now := time.Now()
		existingTicket.ClosedBy = &respondentID
		existingTicket.ClosedAt = &now
	}

	// Update ticket in repository
	err = u.repo.Update(existingTicket)
	if err != nil {
		u.Logger.Error("Error updating customer ticket", map[string]interface{}{
			"error":    err.Error(),
			"ticketId": existingTicket.ID,
		})
		return nil, errors.New("خطا در بروزرسانی تیکت مشتری")
	}

	// Create a new comment
	comment := domain.CustomerComment{
		CustomerTicketID: existingTicket.ID,
		Content:          *params.Comment.Content,
		RespondentID:     *params.Comment.RespondentID,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		IsDeleted:        false,
	}

	// Save comment
	err = u.customerCommentRepo.Create(comment)
	if err != nil {
		u.Logger.Error("Error creating comment for customer ticket", map[string]interface{}{
			"error":    err.Error(),
			"ticketId": existingTicket.ID,
		})
		return nil, errors.New("خطا در ایجاد پیام تیکت مشتری")
	}

	// Handle media attachments
	if len(params.MediaIDs) > 0 {
		err = u.customerTicketMediaRepo.AddMediaToCustomerTicket(existingTicket.ID, params.MediaIDs)
		if err != nil {
			u.Logger.Error("Error adding media to customer ticket", map[string]interface{}{
				"error":    err.Error(),
				"ticketId": existingTicket.ID,
			})
			// Continue despite media error
		}
	}

	// Get the updated ticket with relations
	updatedTicket, err := u.repo.GetByIDWithRelations(existingTicket.ID)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}

	return resp.NewResponseData(
		resp.Created,
		resp.Data{
			"ticket": updatedTicket,
		},
		"پاسخ به تیکت مشتری با موفقیت ایجاد شد",
	), nil
}

func (u *CustomerTicketUsecase) AdminReplayCustomerTicketCommand(params *customer_ticket.AdminReplayCustomerTicketCommand) (*resp.Response, error) {
	u.Logger.Info("AdminReplayCustomerTicketCommand called", map[string]interface{}{
		"id":      *params.ID,
		"status":  *params.Status,
		"content": *params.Comment.Content,
	})

	// Check if user is admin
	isAdmin, err := u.authContext(u.Ctx).IsAdmin()
	if err != nil {
		return nil, errors.New("خطا در بررسی دسترسی کاربر")
	}

	if !isAdmin {
		return nil, errors.New("شما دسترسی به این عملیات را ندارید")
	}

	// Get current user ID
	userID, err := u.authContext(u.Ctx).GetUserID()
	if err != nil {
		return nil, errors.New("خطا در احراز هویت کاربر")
	}

	// Get the existing ticket
	existingTicket, err := u.repo.GetByID(*params.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("تیکت مشتری مورد نظر یافت نشد")
		}
		return nil, err
	}

	// Update the ticket
	existingTicket.Status = *params.Status
	existingTicket.Category = *params.Category
	existingTicket.AssignedTo = params.AssignedTo
	existingTicket.Priority = *params.Priority
	existingTicket.UpdatedAt = time.Now()

	// Set closed info if status is closed
	if *params.Status == enums.TicketClosedStatus {
		now := time.Now()
		existingTicket.ClosedBy = &userID
		existingTicket.ClosedAt = &now
	}

	// Update ticket in repository
	err = u.repo.Update(existingTicket)
	if err != nil {
		u.Logger.Error("Error updating customer ticket", map[string]interface{}{
			"error":    err.Error(),
			"ticketId": existingTicket.ID,
		})
		return nil, errors.New("خطا در بروزرسانی تیکت مشتری")
	}

	// Create a new comment
	comment := domain.CustomerComment{
		CustomerTicketID: existingTicket.ID,
		Content:          *params.Comment.Content,
		RespondentID:     *params.Comment.RespondentID,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		IsDeleted:        false,
	}

	// Save comment
	err = u.customerCommentRepo.Create(comment)
	if err != nil {
		u.Logger.Error("Error creating comment for customer ticket", map[string]interface{}{
			"error":    err.Error(),
			"ticketId": existingTicket.ID,
		})
		return nil, errors.New("خطا در ایجاد پیام تیکت مشتری")
	}

	// Handle media attachments
	if len(params.MediaIDs) > 0 {
		err = u.customerTicketMediaRepo.AddMediaToCustomerTicket(existingTicket.ID, params.MediaIDs)
		if err != nil {
			u.Logger.Error("Error adding media to customer ticket", map[string]interface{}{
				"error":    err.Error(),
				"ticketId": existingTicket.ID,
			})
			// Continue despite media error
		}
	}

	// Get the updated ticket with relations
	updatedTicket, err := u.repo.GetByIDWithRelations(existingTicket.ID)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}

	return resp.NewResponseData(
		resp.Created,
		resp.Data{
			"ticket": updatedTicket,
		},
		"پاسخ به تیکت مشتری با موفقیت ایجاد شد",
	), nil
}

func (u *CustomerTicketUsecase) GetByIdCustomerTicketQuery(params *customer_ticket.GetByIdCustomerTicketQuery) (*resp.Response, error) {
	// Get user ID and determine if admin
	userID, err := u.authContext(u.Ctx).GetUserID()
	if err != nil {
		return nil, errors.New("خطا در احراز هویت کاربر")
	}

	isAdmin, err := u.authContext(u.Ctx).IsAdmin()
	if err != nil {
		return nil, errors.New("خطا در بررسی دسترسی کاربر")
	}

	// Get customer ID if available
	customerID, _ := u.authContext(u.Ctx).GetCustomerID()

	// Get ticket by ID with relations
	result, err := u.repo.GetByIDWithRelations(*params.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("تیکت مشتری مورد نظر یافت نشد")
		}
		return nil, err
	}

	// Check access
	if result.UserID != userID && result.CustomerID != customerID && !isAdmin {
		return nil, errors.New("شما دسترسی به این تیکت مشتری ندارید")
	}

	return resp.NewResponseData(
		resp.Retrieved,
		resp.Data{
			"ticket": result,
		},
		"تیکت مشتری با موفقیت دریافت شد",
	), nil
}

func (u *CustomerTicketUsecase) GetAllCustomerTicketQuery(params *customer_ticket.GetAllCustomerTicketQuery) (*resp.Response, error) {
	u.Logger.Info("GetAllCustomerTicketQuery called", map[string]interface{}{
		"page":     params.Page,
		"pageSize": params.PageSize,
	})

	// Get customer ID from auth context
	customerID, err := u.authContext(u.Ctx).GetCustomerID()
	if err != nil {
		return nil, errors.New("خطا در احراز هویت مشتری")
	}

	// Get all tickets for the current customer with pagination
	tickets, count, err := u.repo.GetAllByCustomerID(customerID, params.PaginationRequestDto)
	if err != nil {
		u.Logger.Error("Error getting all tickets for customer", map[string]interface{}{
			"error":      err.Error(),
			"customerId": customerID,
		})
		return nil, errors.New("خطا در دریافت تیکت های مشتری")
	}

	return map[string]interface{}{
		"items":     enhancedTickets,
		"total":     count,
		"page":      params.Page,
		"pageSize":  params.PageSize,
		"totalPage": (count + int64(params.PageSize) - 1) / int64(params.PageSize),
	}, nil
}

func (u *CustomerTicketUsecase) AdminGetAllCustomerTicketQuery(params *customer_ticket.AdminGetAllCustomerTicketQuery) (*resp.Response, error) {
	u.Logger.Info("AdminGetAllCustomerTicketQuery called", map[string]interface{}{
		"page":     params.Page,
		"pageSize": params.PageSize,
	})

	// Check if user is admin
	isAdmin, err := u.authContext(u.Ctx).IsAdmin()
	if err != nil {
		return nil, errors.New("خطا در بررسی دسترسی کاربر")
	}

	if !isAdmin {
		return nil, errors.New("شما دسترسی به این عملیات را ندارید")
	}

	// Get all tickets with pagination
	tickets, count, err := u.repo.GetAll(params.PaginationRequestDto)
	if err != nil {
		u.Logger.Error("Error getting all customer tickets for admin", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, errors.New("خطا در دریافت تیکت های مشتری")
	}

	// Enhance ticket responses
	enhancedTickets := make([]map[string]interface{}, 0, len(tickets))
	for _, t := range tickets {
		enhancedTickets = append(enhancedTickets, enhanceCustomerTicketResponse(t))
	}

	return map[string]interface{}{
		"items":     enhancedTickets,
		"total":     count,
		"page":      params.Page,
		"pageSize":  params.PageSize,
		"totalPage": (count + int64(params.PageSize) - 1) / int64(params.PageSize),
	}, nil
}
