package customerticketusecase

import (
	"errors"
	"time"

	"github.com/amirex128/new_site_builder/src/internal/domain/enums"

	"github.com/amirex128/new_site_builder/src/internal/application/dto/customer_ticket"
	"github.com/amirex128/new_site_builder/src/internal/application/usecase"
	"github.com/amirex128/new_site_builder/src/internal/application/utils/resp"
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
}

func NewCustomerTicketUsecase(c contract.IContainer) *CustomerTicketUsecase {
	return &CustomerTicketUsecase{
		BaseUsecase: &usecase.BaseUsecase{
			Logger:      c.GetLogger(),
			AuthContext: c.GetAuthTransientService(),
		},
		repo:                    c.GetCustomerTicketRepo(),
		customerCommentRepo:     c.GetCustomerCommentRepo(),
		customerTicketMediaRepo: c.GetCustomerTicketMediaRepo(),
		mediaRepo:               c.GetMediaRepo(),
	}
}

func (u *CustomerTicketUsecase) CreateCustomerTicketCommand(params *customer_ticket.CreateCustomerTicketCommand) (*resp.Response, error) {
	customerID, err := u.AuthContext(u.Ctx).GetCustomerID()
	if err != nil {
		return nil, resp.NewError(resp.Unauthorized, err.Error())
	}

	newTicket := &domain.CustomerTicket{
		Title:      *params.Title,
		CustomerID: *customerID,
		Status:     enums.TicketNewStatus,
		Category:   *params.Category,
		Priority:   *params.Priority,
		UserID:     *params.OwnerUserID,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		IsDeleted:  false,
	}

	err = u.repo.Create(newTicket)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}

	comment := &domain.CustomerComment{
		CustomerTicketID: newTicket.ID,
		Content:          *params.Comment.Content,
		RespondentID:     *params.Comment.RespondentID,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		IsDeleted:        false,
	}
	err = u.customerCommentRepo.Create(comment)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}

	if len(params.MediaIDs) > 0 {
		err = u.customerTicketMediaRepo.AddMediaToCustomerTicket(newTicket.ID, params.MediaIDs)
		if err != nil {
			return nil, resp.NewError(resp.Internal, "خطا در اضافه کردن فایل ها به تیکت مشتری")
		}
	}

	createdTicket, err := u.repo.GetByIDWithRelations(newTicket.ID)
	if err != nil {
		return nil, resp.NewError(resp.NotFound, "تیکت مشتری مورد نظر یافت نشد")
	}

	return resp.NewResponseData(resp.Created, createdTicket, "تیکت مشتری با موفقیت ایجاد شد"), nil
}

func (u *CustomerTicketUsecase) ReplayCustomerTicketCommand(params *customer_ticket.ReplayCustomerTicketCommand) (*resp.Response, error) {
	customerID, err := u.AuthContext(u.Ctx).GetCustomerID()
	if err != nil {
		return nil, resp.NewError(resp.Unauthorized, err.Error())
	}

	existingTicket, err := u.repo.GetByID(*params.ID)
	if err != nil {
		return nil, resp.NewError(resp.NotFound, "تیکت مشتری مورد نظر یافت نشد")
	}

	err = u.CheckAccessCustomerModel(existingTicket, customerID)
	if err != nil {
		return nil, err
	}

	existingTicket.Status = *params.Status
	existingTicket.Category = *params.Category
	existingTicket.AssignedTo = params.AssignedTo
	existingTicket.Priority = *params.Priority
	existingTicket.UpdatedAt = time.Now()

	if *params.Status == enums.TicketClosedStatus {
		respondentID := *params.Comment.RespondentID
		now := time.Now()
		existingTicket.ClosedBy = &respondentID
		existingTicket.ClosedAt = &now
	}

	err = u.repo.Update(existingTicket)
	if err != nil {
		return nil, resp.NewError(resp.Internal, "خطا در بروزرسانی تیکت مشتری")
	}

	comment := &domain.CustomerComment{
		CustomerTicketID: existingTicket.ID,
		Content:          *params.Comment.Content,
		RespondentID:     *params.Comment.RespondentID,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		IsDeleted:        false,
	}
	err = u.customerCommentRepo.Create(comment)
	if err != nil {
		return nil, resp.NewError(resp.Internal, "خطا در ایجاد پیام تیکت مشتری")
	}

	if len(params.MediaIDs) > 0 {
		err = u.customerTicketMediaRepo.AddMediaToCustomerTicket(existingTicket.ID, params.MediaIDs)
		if err != nil {
			return nil, resp.NewError(resp.Internal, "خطا در اضافه کردن فایل ها به تیکت مشتری")
		}
	}

	updatedTicket, err := u.repo.GetByIDWithRelations(existingTicket.ID)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}

	return resp.NewResponseData(resp.Created, updatedTicket, "پاسخ به تیکت مشتری با موفقیت ایجاد شد"), nil
}

func (u *CustomerTicketUsecase) AdminReplayCustomerTicketCommand(params *customer_ticket.AdminReplayCustomerTicketCommand) (*resp.Response, error) {
	userID, err := u.AuthContext(u.Ctx).GetUserID()
	if err != nil {
		return nil, resp.NewError(resp.Unauthorized, err.Error())
	}

	err = u.CheckAccessAdmin()
	if err != nil {
		return nil, err
	}

	existingTicket, err := u.repo.GetByID(*params.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, resp.NewError(resp.NotFound, "تیکت مشتری مورد نظر یافت نشد")
		}
		return nil, resp.NewError(resp.Internal, err.Error())
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
		u.Logger.Error("Error updating customer ticket", map[string]interface{}{
			"error":    err.Error(),
			"ticketId": existingTicket.ID,
		})
		return nil, resp.NewError(resp.Internal, "خطا در بروزرسانی تیکت مشتری")
	}

	comment := &domain.CustomerComment{
		CustomerTicketID: existingTicket.ID,
		Content:          *params.Comment.Content,
		RespondentID:     *params.Comment.RespondentID,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		IsDeleted:        false,
	}
	err = u.customerCommentRepo.Create(comment)
	if err != nil {
		u.Logger.Error("Error creating comment for customer ticket", map[string]interface{}{
			"error":    err.Error(),
			"ticketId": existingTicket.ID,
		})
		return nil, resp.NewError(resp.Internal, "خطا در ایجاد پیام تیکت مشتری")
	}

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

	updatedTicket, err := u.repo.GetByIDWithRelations(existingTicket.ID)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}

	return resp.NewResponseData(resp.Created, map[string]interface{}{
		"ticket": updatedTicket,
	}, "پاسخ به تیکت مشتری با موفقیت ایجاد شد"), nil
}

func (u *CustomerTicketUsecase) GetByIdCustomerTicketQuery(params *customer_ticket.GetByIdCustomerTicketQuery) (*resp.Response, error) {
	userID, customerID, _, err := u.AuthContext(u.Ctx).GetUserOrCustomerID()
	if err != nil {
		return nil, resp.NewError(resp.Unauthorized, err.Error())
	}

	isAdmin, err := u.AuthContext(u.Ctx).IsAdmin()
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}

	result, err := u.repo.GetByIDWithRelations(*params.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, resp.NewError(resp.NotFound, "تیکت مشتری مورد نظر یافت نشد")
		}
		return nil, resp.NewError(resp.Internal, err.Error())
	}

	if (result.UserID != 0 && userID != nil && result.UserID != *userID) && (result.CustomerID != 0 && customerID != nil && result.CustomerID != *customerID) && !isAdmin {
		return nil, resp.NewError(resp.Unauthorized, "شما دسترسی به این تیکت مشتری ندارید")
	}

	return resp.NewResponseData(resp.Retrieved, map[string]interface{}{
		"ticket": result,
	}, "تیکت مشتری با موفقیت دریافت شد"), nil
}

func (u *CustomerTicketUsecase) GetAllCustomerTicketQuery(params *customer_ticket.GetAllCustomerTicketQuery) (*resp.Response, error) {
	u.Logger.Info("GetAllCustomerTicketQuery called", map[string]interface{}{
		"page":     params.Page,
		"pageSize": params.PageSize,
	})

	_, customerID, _, err := u.AuthContext(u.Ctx).GetUserOrCustomerID()
	if err != nil {
		return nil, resp.NewError(resp.Unauthorized, err.Error())
	}
	if customerID == nil {
		return nil, resp.NewError(resp.Unauthorized, "شناسه مشتری یافت نشد")
	}

	results, err := u.repo.GetAllByCustomerID(*customerID, params.PaginationRequestDto)
	if err != nil {
		u.Logger.Error("Error getting all tickets for customer", map[string]interface{}{
			"error":      err.Error(),
			"customerId": *customerID,
		})
		return nil, resp.NewError(resp.Internal, "خطا در دریافت تیکت های مشتری")
	}

	return resp.NewResponseData(resp.Retrieved, map[string]interface{}{
		"items":     results.Items,
		"total":     results.TotalCount,
		"page":      results.PageNumber,
		"pageSize":  params.PageSize,
		"totalPage": results.TotalPages,
	}, "تیکت های مشتری با موفقیت دریافت شدند"), nil
}

func (u *CustomerTicketUsecase) AdminGetAllCustomerTicketQuery(params *customer_ticket.AdminGetAllCustomerTicketQuery) (*resp.Response, error) {
	u.Logger.Info("AdminGetAllCustomerTicketQuery called", map[string]interface{}{
		"page":     params.Page,
		"pageSize": params.PageSize,
	})

	isAdmin, err := u.AuthContext(u.Ctx).IsAdmin()
	if err != nil || !isAdmin {
		return nil, resp.NewError(resp.Unauthorized, "شما دسترسی به این عملیات را ندارید")
	}

	results, err := u.repo.GetAll(params.PaginationRequestDto)
	if err != nil {
		u.Logger.Error("Error getting all customer tickets for admin", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, resp.NewError(resp.Internal, "خطا در دریافت تیکت های مشتری")
	}

	return resp.NewResponseData(resp.Retrieved, map[string]interface{}{
		"items":     results.Items,
		"total":     results.TotalCount,
		"page":      results.PageNumber,
		"pageSize":  params.PageSize,
		"totalPage": results.TotalPages,
	}, "تیکت های مشتری با موفقیت دریافت شدند (ادمین)"), nil
}
