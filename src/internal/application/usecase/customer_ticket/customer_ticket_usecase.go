package customerticketusecase

import (
	"errors"
	"github.com/amirex128/new_site_builder/src/internal/contract/service"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	"github.com/amirex128/new_site_builder/src/internal/application/dto/customer_ticket"
	"github.com/amirex128/new_site_builder/src/internal/contract"
	"github.com/amirex128/new_site_builder/src/internal/contract/repository"
	"github.com/amirex128/new_site_builder/src/internal/domain"
	"gorm.io/gorm"
)

type CustomerTicketUsecase struct {
	ctx                     *gin.Context
	logger                  sflogger.Logger
	repo                    repository.ICustomerTicketRepository
	customerCommentRepo     repository.ICustomerCommentRepository
	customerTicketMediaRepo repository.ICustomerTicketMediaRepository
	mediaRepo               repository.IMediaRepository
	authContext             func(c *gin.Context) service.IAuthService
}

func NewCustomerTicketUsecase(c contract.IContainer) *CustomerTicketUsecase {
	return &CustomerTicketUsecase{
		logger:                  c.GetLogger(),
		repo:                    c.GetCustomerTicketRepo(),
		customerCommentRepo:     c.GetCustomerCommentRepo(),
		customerTicketMediaRepo: c.GetCustomerTicketMediaRepo(),
		mediaRepo:               c.GetMediaRepo(),
		authContext:             c.GetAuthTransientService(),
	}
}

func (u *CustomerTicketUsecase) SetContext(c *gin.Context) *CustomerTicketUsecase {
	u.ctx = c
	return u
}

func (u *CustomerTicketUsecase) CreateCustomerTicketCommand(params *customer_ticket.CreateCustomerTicketCommand) (any, error) {
	u.logger.Info("CreateCustomerTicketCommand called", map[string]interface{}{
		"title":    *params.Title,
		"category": *params.Category,
		"priority": *params.Priority,
	})

	// Get current user ID from auth context
	userID, err := u.authContext(u.ctx).GetUserID()
	if err != nil {
		return nil, errors.New("خطا در احراز هویت کاربر")
	}

	// Create the customer ticket
	newTicket := domain.CustomerTicket{
		Title:      *params.Title,
		CustomerID: *params.OwnerUserID,
		Status:     strconv.Itoa(int(customer_ticket.New)), // Default to New status
		Category:   strconv.Itoa(int(*params.Category)),
		Priority:   strconv.Itoa(int(*params.Priority)),
		UserID:     userID,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		IsDeleted:  false,
	}

	// Save ticket to repository
	err = u.repo.Create(newTicket)
	if err != nil {
		u.logger.Error("Error creating customer ticket", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, errors.New("خطا در ایجاد تیکت مشتری")
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
		u.logger.Error("Error creating comment for customer ticket", map[string]interface{}{
			"error":    err.Error(),
			"ticketId": newTicket.ID,
		})
		return nil, errors.New("خطا در ایجاد پیام تیکت مشتری")
	}

	// Handle media attachments
	if len(params.MediaIDs) > 0 {
		err = u.customerTicketMediaRepo.AddMediaToCustomerTicket(newTicket.ID, params.MediaIDs)
		if err != nil {
			u.logger.Error("Error adding media to customer ticket", map[string]interface{}{
				"error":    err.Error(),
				"ticketId": newTicket.ID,
			})
			// Continue despite media error
		}
	}

	// Get the created ticket with relations
	createdTicket, err := u.repo.GetByIDWithRelations(newTicket.ID)
	if err != nil {
		return newTicket, nil // Return basic ticket if relations can't be loaded
	}

	return enhanceCustomerTicketResponse(createdTicket), nil
}

func (u *CustomerTicketUsecase) ReplayCustomerTicketCommand(params *customer_ticket.ReplayCustomerTicketCommand) (any, error) {
	u.logger.Info("ReplayCustomerTicketCommand called", map[string]interface{}{
		"id":      *params.ID,
		"status":  *params.Status,
		"content": *params.Comment.Content,
	})

	// Get customer ID from auth context
	customerID, err := u.authContext(u.ctx).GetCustomerID()
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
	existingTicket.Status = strconv.Itoa(int(*params.Status))
	existingTicket.Category = strconv.Itoa(int(*params.Category))
	existingTicket.AssignedTo = params.AssignedTo
	existingTicket.Priority = strconv.Itoa(int(*params.Priority))
	existingTicket.UpdatedAt = time.Now()

	// Set closed info if status is closed
	if *params.Status == customer_ticket.Closed {
		respondentID := *params.Comment.RespondentID
		now := time.Now()
		existingTicket.ClosedBy = &respondentID
		existingTicket.ClosedAt = &now
	}

	// Update ticket in repository
	err = u.repo.Update(existingTicket)
	if err != nil {
		u.logger.Error("Error updating customer ticket", map[string]interface{}{
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
		u.logger.Error("Error creating comment for customer ticket", map[string]interface{}{
			"error":    err.Error(),
			"ticketId": existingTicket.ID,
		})
		return nil, errors.New("خطا در ایجاد پیام تیکت مشتری")
	}

	// Handle media attachments
	if len(params.MediaIDs) > 0 {
		err = u.customerTicketMediaRepo.AddMediaToCustomerTicket(existingTicket.ID, params.MediaIDs)
		if err != nil {
			u.logger.Error("Error adding media to customer ticket", map[string]interface{}{
				"error":    err.Error(),
				"ticketId": existingTicket.ID,
			})
			// Continue despite media error
		}
	}

	// Get the updated ticket with relations
	updatedTicket, err := u.repo.GetByIDWithRelations(existingTicket.ID)
	if err != nil {
		return existingTicket, nil // Return basic ticket if relations can't be loaded
	}

	return enhanceCustomerTicketResponse(updatedTicket), nil
}

func (u *CustomerTicketUsecase) AdminReplayCustomerTicketCommand(params *customer_ticket.AdminReplayCustomerTicketCommand) (any, error) {
	u.logger.Info("AdminReplayCustomerTicketCommand called", map[string]interface{}{
		"id":      *params.ID,
		"status":  *params.Status,
		"content": *params.Comment.Content,
	})

	// Check if user is admin
	isAdmin, err := u.authContext(u.ctx).IsAdmin()
	if err != nil {
		return nil, errors.New("خطا در بررسی دسترسی کاربر")
	}

	if !isAdmin {
		return nil, errors.New("شما دسترسی به این عملیات را ندارید")
	}

	// Get current user ID
	userID, err := u.authContext(u.ctx).GetUserID()
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
	existingTicket.Status = strconv.Itoa(int(*params.Status))
	existingTicket.Category = strconv.Itoa(int(*params.Category))
	existingTicket.AssignedTo = params.AssignedTo
	existingTicket.Priority = strconv.Itoa(int(*params.Priority))
	existingTicket.UpdatedAt = time.Now()

	// Set closed info if status is closed
	if *params.Status == customer_ticket.Closed {
		now := time.Now()
		existingTicket.ClosedBy = &userID
		existingTicket.ClosedAt = &now
	}

	// Update ticket in repository
	err = u.repo.Update(existingTicket)
	if err != nil {
		u.logger.Error("Error updating customer ticket", map[string]interface{}{
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
		u.logger.Error("Error creating comment for customer ticket", map[string]interface{}{
			"error":    err.Error(),
			"ticketId": existingTicket.ID,
		})
		return nil, errors.New("خطا در ایجاد پیام تیکت مشتری")
	}

	// Handle media attachments
	if len(params.MediaIDs) > 0 {
		err = u.customerTicketMediaRepo.AddMediaToCustomerTicket(existingTicket.ID, params.MediaIDs)
		if err != nil {
			u.logger.Error("Error adding media to customer ticket", map[string]interface{}{
				"error":    err.Error(),
				"ticketId": existingTicket.ID,
			})
			// Continue despite media error
		}
	}

	// Get the updated ticket with relations
	updatedTicket, err := u.repo.GetByIDWithRelations(existingTicket.ID)
	if err != nil {
		return existingTicket, nil // Return basic ticket if relations can't be loaded
	}

	return enhanceCustomerTicketResponse(updatedTicket), nil
}

func (u *CustomerTicketUsecase) GetByIdCustomerTicketQuery(params *customer_ticket.GetByIdCustomerTicketQuery) (any, error) {
	u.logger.Info("GetByIdCustomerTicketQuery called", map[string]interface{}{
		"id": *params.ID,
	})

	// Get user ID and determine if admin
	userID, err := u.authContext(u.ctx).GetUserID()
	if err != nil {
		return nil, errors.New("خطا در احراز هویت کاربر")
	}

	isAdmin, err := u.authContext(u.ctx).IsAdmin()
	if err != nil {
		return nil, errors.New("خطا در بررسی دسترسی کاربر")
	}

	// Get customer ID if available
	customerID, _ := u.authContext(u.ctx).GetCustomerID()

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

	return enhanceCustomerTicketResponse(result), nil
}

func (u *CustomerTicketUsecase) GetAllCustomerTicketQuery(params *customer_ticket.GetAllCustomerTicketQuery) (any, error) {
	u.logger.Info("GetAllCustomerTicketQuery called", map[string]interface{}{
		"page":     params.Page,
		"pageSize": params.PageSize,
	})

	// Get customer ID from auth context
	customerID, err := u.authContext(u.ctx).GetCustomerID()
	if err != nil {
		return nil, errors.New("خطا در احراز هویت مشتری")
	}

	// Get all tickets for the current customer with pagination
	tickets, count, err := u.repo.GetAllByCustomerID(customerID, params.PaginationRequestDto)
	if err != nil {
		u.logger.Error("Error getting all tickets for customer", map[string]interface{}{
			"error":      err.Error(),
			"customerId": customerID,
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

func (u *CustomerTicketUsecase) AdminGetAllCustomerTicketQuery(params *customer_ticket.AdminGetAllCustomerTicketQuery) (any, error) {
	u.logger.Info("AdminGetAllCustomerTicketQuery called", map[string]interface{}{
		"page":     params.Page,
		"pageSize": params.PageSize,
	})

	// Check if user is admin
	isAdmin, err := u.authContext(u.ctx).IsAdmin()
	if err != nil {
		return nil, errors.New("خطا در بررسی دسترسی کاربر")
	}

	if !isAdmin {
		return nil, errors.New("شما دسترسی به این عملیات را ندارید")
	}

	// Get all tickets with pagination
	tickets, count, err := u.repo.GetAll(params.PaginationRequestDto)
	if err != nil {
		u.logger.Error("Error getting all customer tickets for admin", map[string]interface{}{
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

// Helper function to enhance customer ticket response with structured data
func enhanceCustomerTicketResponse(t domain.CustomerTicket) map[string]interface{} {
	statusEnum, _ := strconv.Atoi(t.Status)
	categoryEnum, _ := strconv.Atoi(t.Category)
	priorityEnum, _ := strconv.Atoi(t.Priority)

	response := map[string]interface{}{
		"id":         t.ID,
		"title":      t.Title,
		"status":     customer_ticket.CustomerTicketStatusEnum(statusEnum),
		"category":   customer_ticket.CustomerTicketCategoryEnum(categoryEnum),
		"assignedTo": t.AssignedTo,
		"priority":   customer_ticket.CustomerTicketPriorityEnum(priorityEnum),
		"userId":     t.UserID,
		"customerId": t.CustomerID,
		"createdAt":  t.CreatedAt,
		"updatedAt":  t.UpdatedAt,
	}

	// Add comments if available
	if len(t.Comments) > 0 {
		comments := make([]map[string]interface{}, 0, len(t.Comments))
		for _, c := range t.Comments {
			commentData := map[string]interface{}{
				"id":           c.ID,
				"content":      c.Content,
				"respondentId": c.RespondentID,
				"createdAt":    c.CreatedAt,
			}
			if c.Respondent != nil {
				commentData["respondent"] = map[string]interface{}{
					"id":        c.Respondent.ID,
					"firstName": c.Respondent.FirstName,
					"lastName":  c.Respondent.LastName,
					"email":     c.Respondent.Email,
				}
			}
			comments = append(comments, commentData)
		}
		response["comments"] = comments
	}

	// Add media if available
	if len(t.Media) > 0 {
		mediaItems := make([]map[string]interface{}, 0, len(t.Media))
		for _, m := range t.Media {
			// Since the Media struct doesn't have many fields directly,
			// we're just adding the ID. In a real implementation, you would
			// join with file_items or another table with detailed media info
			mediaItem := map[string]interface{}{
				"id": m.ID,
			}
			mediaItems = append(mediaItems, mediaItem)
		}
		response["media"] = mediaItems
	}

	// Add user info if available
	if t.User != nil {
		response["user"] = map[string]interface{}{
			"id":        t.User.ID,
			"firstName": t.User.FirstName,
			"lastName":  t.User.LastName,
			"email":     t.User.Email,
		}
	}

	// Add customer info if available
	if t.Customer != nil {
		response["customer"] = map[string]interface{}{
			"id":        t.Customer.ID,
			"firstName": t.Customer.FirstName,
			"lastName":  t.Customer.LastName,
			"email":     t.Customer.Email,
		}
	}

	// Add assigned user info if available
	if t.AssignedTo != nil && t.Assigned != nil {
		response["assignedUser"] = map[string]interface{}{
			"id":        t.Assigned.ID,
			"firstName": t.Assigned.FirstName,
			"lastName":  t.Assigned.LastName,
			"email":     t.Assigned.Email,
		}
	}

	// Add closed info if available
	if t.ClosedAt != nil {
		response["closedAt"] = t.ClosedAt
		response["closedBy"] = t.ClosedBy
		if t.Closer != nil {
			response["closingUser"] = map[string]interface{}{
				"id":        t.Closer.ID,
				"firstName": t.Closer.FirstName,
				"lastName":  t.Closer.LastName,
				"email":     t.Closer.Email,
			}
		}
	}

	return response
}
