package ticketusecase

import (
	"errors"
	"time"

	"github.com/amirex128/new_site_builder/src/internal/domain/enums"

	"github.com/amirex128/new_site_builder/src/internal/application/dto/ticket"
	"github.com/amirex128/new_site_builder/src/internal/application/usecase"
	"github.com/amirex128/new_site_builder/src/internal/application/utils/resp"
	"github.com/amirex128/new_site_builder/src/internal/contract"
	"github.com/amirex128/new_site_builder/src/internal/contract/repository"
	"github.com/amirex128/new_site_builder/src/internal/domain"
	"gorm.io/gorm"
)

type TicketUsecase struct {
	*usecase.BaseUsecase
	repo            repository.ITicketRepository
	commentRepo     repository.ICommentRepository
	ticketMediaRepo repository.ITicketMediaRepository
	mediaRepo       repository.IMediaRepository
}

func NewTicketUsecase(c contract.IContainer) *TicketUsecase {
	return &TicketUsecase{
		BaseUsecase: &usecase.BaseUsecase{
			Logger:      c.GetLogger(),
			AuthContext: c.GetAuthTransientService(),
		},
		repo:            c.GetTicketRepo(),
		commentRepo:     c.GetCommentRepo(),
		ticketMediaRepo: c.GetTicketMediaRepo(),
		mediaRepo:       c.GetMediaRepo(),
	}
}

func (u *TicketUsecase) CreateTicketCommand(params *ticket.CreateTicketCommand) (*resp.Response, error) {
	userId, err := u.AuthContext(u.Ctx).GetUserID()
	if err != nil {
		return nil, err
	}

	// Create the ticket
	newTicket := domain.Ticket{
		Title:     *params.Title,
		Status:    enums.TicketNewStatus, // Default to New status
		Category:  *params.Category,
		Priority:  *params.Priority,
		UserID:    *userId,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		IsDeleted: false,
	}

	// Save ticket to repository
	err = u.repo.Create(&newTicket)
	if err != nil {
		return nil, resp.NewError(resp.Internal, "خطا در ایجاد تیکت")
	}

	// Create the first comment
	comment := domain.Comment{
		TicketID:     newTicket.ID,
		Content:      *params.Comment.Content,
		RespondentID: *userId,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		IsDeleted:    false,
	}

	// Save comment
	err = u.commentRepo.Create(&comment)
	if err != nil {
		return nil, resp.NewError(resp.Internal, "خطا در ایجاد پیام تیکت")
	}

	// Handle media attachments
	if len(params.MediaIDs) > 0 {
		err = u.ticketMediaRepo.AddMediaToTicket(newTicket.ID, params.MediaIDs)
		if err != nil {
			// continue
		}
	}

	// Get the created ticket with relations
	createdTicket, err := u.repo.GetByIDWithRelations(newTicket.ID)
	if err != nil {
		return resp.NewResponseData(resp.Created, resp.Data{"ticket": newTicket}, "تیکت با موفقیت ایجاد شد"), nil
	}

	return resp.NewResponseData(resp.Created, enhanceTicketResponse(*createdTicket), "تیکت با موفقیت ایجاد شد"), nil
}

func (u *TicketUsecase) ReplayTicketCommand(params *ticket.ReplayTicketCommand) (*resp.Response, error) {
	userID, err := u.AuthContext(u.Ctx).GetUserID()
	if err != nil {
		return nil, err
	}
	existingTicket, err := u.repo.GetByID(*params.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, resp.NewError(resp.NotFound, "تیکت مورد نظر یافت نشد")
		}
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	if userID != nil && existingTicket.UserID != *userID {
		isAdmin, err := u.AuthContext(u.Ctx).IsAdmin()
		if err != nil || !isAdmin {
			return nil, resp.NewError(resp.Unauthorized, "شما دسترسی به این تیکت ندارید")
		}
	}
	existingTicket.Status = *params.Status
	existingTicket.Category = *params.Category
	existingTicket.AssignedTo = params.AssignedTo
	existingTicket.Priority = *params.Priority
	existingTicket.UpdatedAt = time.Now()
	if *params.Status == enums.TicketClosedStatus {
		now := time.Now()
		existingTicket.ClosedBy = userID
		existingTicket.ClosedAt = &now
	}
	err = u.repo.Update(existingTicket)
	if err != nil {
		return nil, resp.NewError(resp.Internal, "خطا در بروزرسانی تیکت")
	}
	comment := domain.Comment{
		TicketID:     existingTicket.ID,
		Content:      *params.Comment.Content,
		RespondentID: *userID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		IsDeleted:    false,
	}
	err = u.commentRepo.Create(&comment)
	if err != nil {
		return nil, resp.NewError(resp.Internal, "خطا در ایجاد پیام تیکت")
	}
	if len(params.MediaIDs) > 0 {
		err = u.ticketMediaRepo.AddMediaToTicket(existingTicket.ID, params.MediaIDs)
		if err != nil {
			// continue
		}
	}
	updatedTicket, err := u.repo.GetByIDWithRelations(existingTicket.ID)
	if err != nil {
		return resp.NewResponseData(resp.Updated, resp.Data{"ticket": existingTicket}, "تیکت با موفقیت بروزرسانی شد"), nil
	}
	return resp.NewResponseData(resp.Updated, enhanceTicketResponse(*updatedTicket), "تیکت با موفقیت بروزرسانی شد"), nil
}

func (u *TicketUsecase) AdminReplayTicketCommand(params *ticket.AdminReplayTicketCommand) (*resp.Response, error) {
	isAdmin, err := u.AuthContext(u.Ctx).IsAdmin()
	if err != nil || !isAdmin {
		return nil, resp.NewError(resp.Unauthorized, "شما دسترسی به این عملیات را ندارید")
	}
	userID, err := u.AuthContext(u.Ctx).GetUserID()
	if err != nil {
		return nil, err
	}
	existingTicket, err := u.repo.GetByID(*params.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, resp.NewError(resp.NotFound, "تیکت مورد نظر یافت نشد")
		}
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	existingTicket.Status = *params.Status
	existingTicket.Category = *params.Category
	if params.AssignedTo != nil {
		existingTicket.AssignedTo = params.AssignedTo
	}
	existingTicket.Priority = *params.Priority
	existingTicket.UpdatedAt = time.Now()
	if *params.Status == enums.TicketClosedStatus {
		now := time.Now()
		existingTicket.ClosedBy = userID
		existingTicket.ClosedAt = &now
	}
	err = u.repo.Update(existingTicket)
	if err != nil {
		return nil, resp.NewError(resp.Internal, "خطا در بروزرسانی تیکت")
	}
	comment := domain.Comment{
		TicketID:     existingTicket.ID,
		Content:      *params.Comment.Content,
		RespondentID: *userID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		IsDeleted:    false,
	}
	err = u.commentRepo.Create(&comment)
	if err != nil {
		return nil, resp.NewError(resp.Internal, "خطا در ایجاد پیام تیکت")
	}
	if len(params.MediaIDs) > 0 {
		err = u.ticketMediaRepo.AddMediaToTicket(existingTicket.ID, params.MediaIDs)
		if err != nil {
			// continue
		}
	}
	updatedTicket, err := u.repo.GetByIDWithRelations(existingTicket.ID)
	if err != nil {
		return resp.NewResponseData(resp.Updated, resp.Data{"ticket": existingTicket}, "تیکت با موفقیت بروزرسانی شد"), nil
	}
	return resp.NewResponseData(resp.Updated, enhanceTicketResponse(*updatedTicket), "تیکت با موفقیت بروزرسانی شد"), nil
}

func (u *TicketUsecase) GetByIdTicketQuery(params *ticket.GetByIdTicketQuery) (*resp.Response, error) {
	userID, err := u.AuthContext(u.Ctx).GetUserID()
	if err != nil {
		return nil, err
	}
	result, err := u.repo.GetByIDWithRelations(*params.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, resp.NewError(resp.NotFound, "تیکت مورد نظر یافت نشد")
		}
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	if userID != nil && result.UserID != *userID {
		isAdmin, err := u.AuthContext(u.Ctx).IsAdmin()
		if err != nil || !isAdmin {
			return nil, resp.NewError(resp.Unauthorized, "شما دسترسی به این تیکت ندارید")
		}
	}
	return resp.NewResponseData(resp.Retrieved, enhanceTicketResponse(*result), "تیکت با موفقیت دریافت شد"), nil
}

func (u *TicketUsecase) GetAllTicketQuery(params *ticket.GetAllTicketQuery) (*resp.Response, error) {
	userID, err := u.AuthContext(u.Ctx).GetUserID()
	if err != nil {
		return nil, err
	}
	ticketsResult, err := u.repo.GetAllByUserID(*userID, params.PaginationRequestDto)
	if err != nil {
		return nil, resp.NewError(resp.Internal, "خطا در دریافت تیکت ها")
	}
	enhancedTickets := make([]map[string]interface{}, 0, len(ticketsResult.Items))
	for _, t := range ticketsResult.Items {
		enhancedTickets = append(enhancedTickets, enhanceTicketResponse(t))
	}
	return resp.NewResponseData(resp.Retrieved, map[string]interface{}{
		"items":     enhancedTickets,
		"total":     ticketsResult.TotalCount,
		"page":      ticketsResult.PageNumber,
		"pageSize":  params.PageSize,
		"totalPage": ticketsResult.TotalPages,
	}, "لیست تیکت ها با موفقیت دریافت شد"), nil
}

func (u *TicketUsecase) AdminGetAllTicketQuery(params *ticket.AdminGetAllTicketQuery) (*resp.Response, error) {
	isAdmin, err := u.AuthContext(u.Ctx).IsAdmin()
	if err != nil || !isAdmin {
		return nil, resp.NewError(resp.Unauthorized, "شما دسترسی به این عملیات را ندارید")
	}
	ticketsResult, err := u.repo.GetAll(params.PaginationRequestDto)
	if err != nil {
		return nil, resp.NewError(resp.Internal, "خطا در دریافت تیکت ها")
	}
	enhancedTickets := make([]map[string]interface{}, 0, len(ticketsResult.Items))
	for _, t := range ticketsResult.Items {
		enhancedTickets = append(enhancedTickets, enhanceTicketResponse(t))
	}
	return resp.NewResponseData(resp.Retrieved, map[string]interface{}{
		"items":     enhancedTickets,
		"total":     ticketsResult.TotalCount,
		"page":      ticketsResult.PageNumber,
		"pageSize":  params.PageSize,
		"totalPage": ticketsResult.TotalPages,
	}, "لیست تیکت ها با موفقیت دریافت شد"), nil
}

// Helper function to enhance ticket response with structured data
func enhanceTicketResponse(t domain.Ticket) map[string]interface{} {
	response := map[string]interface{}{
		"id":         t.ID,
		"title":      t.Title,
		"status":     t.Status,
		"category":   t.Category,
		"assignedTo": t.AssignedTo,
		"priority":   t.Priority,
		"userId":     t.UserID,
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
