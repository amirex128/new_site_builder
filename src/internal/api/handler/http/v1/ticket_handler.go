package v1

import (
	"github.com/amirex128/new_site_builder/src/internal/api/utils"
	"github.com/amirex128/new_site_builder/src/internal/application/dto/ticket"
	ticketusecase "github.com/amirex128/new_site_builder/src/internal/application/usecase/ticket"
	"github.com/gin-gonic/gin"
)

type TicketHandler struct {
	usecase   *ticketusecase.TicketUsecase
	validator *utils.ValidationHelper
}

func NewTicketHandler(usc *ticketusecase.TicketUsecase) *TicketHandler {
	return &TicketHandler{
		usecase:   usc,
		validator: utils.NewValidationHelper(),
	}
}

// CreateTicket godoc
// @Summary      Create a new ticket
// @Description  Creates a new support ticket for the authenticated user
// @Tags         ticket
// @Accept       json
// @Produce      json
// @Param        request  body      ticket.CreateTicketCommand  true  "Ticket information"
// @Success      201      {object}  utils.Result                "Created ticket"
// @Failure      400      {object}  utils.Result                "Validation error"
// @Failure      401      {object}  utils.Result                "Unauthorized"
// @Failure      500      {object}  utils.Result                "Internal server error"
// @Router       /ticket [post]
// @Security BearerAuth
func (h *TicketHandler) CreateTicket(c *gin.Context) {
	var params ticket.CreateTicketCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}

	result, err := h.usecase.CreateTicketCommand(&params)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Created(c, result)
}

// ReplayTicket godoc
// @Summary      Reply to a ticket
// @Description  Adds a user reply to an existing support ticket
// @Tags         ticket
// @Accept       json
// @Produce      json
// @Param        request  body      ticket.ReplayTicketCommand  true  "Reply information"
// @Success      200      {object}  utils.Result                "Updated ticket with reply"
// @Failure      400      {object}  utils.Result                "Validation error"
// @Failure      401      {object}  utils.Result                "Unauthorized"
// @Failure      404      {object}  utils.Result                "Ticket not found"
// @Failure      500      {object}  utils.Result                "Internal server error"
// @Router       /ticket/reply [post]
// @Security BearerAuth
func (h *TicketHandler) ReplayTicket(c *gin.Context) {
	var params ticket.ReplayTicketCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}

	result, err := h.usecase.ReplayTicketCommand(&params)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Updated(c, result)
}

// AdminReplayTicket godoc
// @Summary      Admin: Reply to a ticket
// @Description  Admin endpoint to add a reply to any support ticket
// @Tags         ticket
// @Accept       json
// @Produce      json
// @Param        request  body      ticket.AdminReplayTicketCommand  true  "Admin reply information"
// @Success      200      {object}  utils.Result                     "Updated ticket with admin reply"
// @Failure      400      {object}  utils.Result                     "Validation error"
// @Failure      401      {object}  utils.Result                     "Unauthorized"
// @Failure      403      {object}  utils.Result                     "Forbidden - Admin access required"
// @Failure      404      {object}  utils.Result                     "Ticket not found"
// @Failure      500      {object}  utils.Result                     "Internal server error"
// @Router       /ticket/admin/reply [post]
// @Security BearerAuth
func (h *TicketHandler) AdminReplayTicket(c *gin.Context) {
	var params ticket.AdminReplayTicketCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}

	result, err := h.usecase.AdminReplayTicketCommand(&params)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Updated(c, result)
}

// GetByIdTicket godoc
// @Summary      Get ticket by ID
// @Description  Retrieves a specific ticket by its ID for the authenticated user
// @Tags         ticket
// @Accept       json
// @Produce      json
// @Param        request  query     ticket.GetByIdTicketQuery  true  "Ticket ID to retrieve"
// @Success      200      {object}  utils.Result               "Ticket details"
// @Failure      400      {object}  utils.Result               "Validation error"
// @Failure      401      {object}  utils.Result               "Unauthorized"
// @Failure      404      {object}  utils.Result               "Ticket not found"
// @Failure      500      {object}  utils.Result               "Internal server error"
// @Router       /ticket [get]
// @Security BearerAuth
func (h *TicketHandler) GetByIdTicket(c *gin.Context) {
	var params ticket.GetByIdTicketQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}

	result, err := h.usecase.GetByIdTicketQuery(&params)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Retrieved(c, result)
}

// GetAllTicket godoc
// @Summary      Get all tickets
// @Description  Retrieves all tickets for the authenticated user
// @Tags         ticket
// @Accept       json
// @Produce      json
// @Param        request  query     ticket.GetAllTicketQuery  true  "Query parameters"
// @Success      200      {object}  utils.Result              "List of tickets"
// @Failure      400      {object}  utils.Result              "Validation error"
// @Failure      401      {object}  utils.Result              "Unauthorized"
// @Failure      500      {object}  utils.Result              "Internal server error"
// @Router       /ticket/all [get]
// @Security BearerAuth
func (h *TicketHandler) GetAllTicket(c *gin.Context) {
	var params ticket.GetAllTicketQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}

	result, err := h.usecase.GetAllTicketQuery(&params)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Retrieved(c, result)
}

// AdminGetAllTicket godoc
// @Summary      Admin: Get all tickets
// @Description  Admin endpoint to retrieve all tickets in the system
// @Tags         ticket
// @Accept       json
// @Produce      json
// @Param        request  query     ticket.AdminGetAllTicketQuery  true  "Query parameters"
// @Success      200      {object}  utils.Result                   "List of all tickets"
// @Failure      400      {object}  utils.Result                   "Validation error"
// @Failure      401      {object}  utils.Result                   "Unauthorized"
// @Failure      403      {object}  utils.Result                   "Forbidden - Admin access required"
// @Failure      500      {object}  utils.Result                   "Internal server error"
// @Router       /ticket/admin/all [get]
// @Security BearerAuth
func (h *TicketHandler) AdminGetAllTicket(c *gin.Context) {
	var params ticket.AdminGetAllTicketQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}

	result, err := h.usecase.AdminGetAllTicketQuery(&params)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Retrieved(c, result)
}
