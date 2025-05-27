package v1

import (
	"github.com/amirex128/new_site_builder/src/internal/api/utils"
	"github.com/amirex128/new_site_builder/src/internal/application/dto/customer_ticket"
	customerticketusecase "github.com/amirex128/new_site_builder/src/internal/application/usecase/customer_ticket"
	"github.com/gin-gonic/gin"
)

type CustomerTicketHandler struct {
	usecase   *customerticketusecase.CustomerTicketUsecase
	validator *utils.ValidationHelper
}

func NewCustomerTicketHandler(usc *customerticketusecase.CustomerTicketUsecase) *CustomerTicketHandler {
	return &CustomerTicketHandler{
		usecase:   usc,
		validator: utils.NewValidationHelper(),
	}
}

// CreateCustomerTicket godoc
// @Summary      Create a customer support ticket
// @Description  Creates a new support ticket for a customer
// @Tags         customer-ticket
// @Accept       json
// @Produce      json
// @Param        request  body      customer_ticket.CreateCustomerTicketCommand  true  "Ticket information"
// @success      201      {object}  utils.Result                                  "Created ticket"
// @Failure      400      {object}  utils.Result                                  "Validation error"
// @Failure      401      {object}  utils.Result                                  "unauthorized"
// @Failure      500      {object}  utils.Result                                  "Internal server error"
// @Router       /customer-ticket [post]
// @Security BearerAuth
func (h *CustomerTicketHandler) CreateCustomerTicket(c *gin.Context) {
	var params customer_ticket.CreateCustomerTicketCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}

	result, err := h.usecase.CreateCustomerTicketCommand(&params)
	utils.HandleError(c, err)
	utils.HandleResponse(c, result)
}

// ReplayCustomerTicket godoc
// @Summary      Reply to a customer ticket
// @Description  Adds a customer reply to an existing support ticket
// @Tags         customer-ticket
// @Accept       json
// @Produce      json
// @Param        request  body      customer_ticket.ReplayCustomerTicketCommand  true  "Reply information"
// @success      200      {object}  utils.Result                                  "Updated ticket with reply"
// @Failure      400      {object}  utils.Result                                  "Validation error"
// @Failure      401      {object}  utils.Result                                  "unauthorized"
// @Failure      404      {object}  utils.Result                                  "Ticket not found"
// @Failure      500      {object}  utils.Result                                  "Internal server error"
// @Router       /customer-ticket [put]
// @Security BearerAuth
func (h *CustomerTicketHandler) ReplayCustomerTicket(c *gin.Context) {
	var params customer_ticket.ReplayCustomerTicketCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}

	result, err := h.usecase.ReplayCustomerTicketCommand(&params)
	utils.HandleError(c, err)
	utils.HandleResponse(c, result)
}

// AdminReplayCustomerTicket godoc
// @Summary      Admin: Reply to a customer ticket
// @Description  Adds an admin reply to an existing customer support ticket
// @Tags         customer-ticket
// @Accept       json
// @Produce      json
// @Param        request  body      customer_ticket.AdminReplayCustomerTicketCommand  true  "Admin reply information"
// @success      200      {object}  utils.Result                                       "Updated ticket with admin reply"
// @Failure      400      {object}  utils.Result                                       "Validation error"
// @Failure      401      {object}  utils.Result                                       "unauthorized"
// @Failure      403      {object}  utils.Result                                       "Forbidden - Admin access required"
// @Failure      404      {object}  utils.Result                                       "Ticket not found"
// @Failure      500      {object}  utils.Result                                       "Internal server error"
// @Router       /customer-ticket/admin [put]
// @Security BearerAuth
func (h *CustomerTicketHandler) AdminReplayCustomerTicket(c *gin.Context) {
	var params customer_ticket.AdminReplayCustomerTicketCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}

	result, err := h.usecase.AdminReplayCustomerTicketCommand(&params)
	utils.HandleError(c, err)
	utils.HandleResponse(c, result)
}

// GetByIdCustomerTicket godoc
// @Summary      Get customer ticket by ID
// @Description  Retrieves a specific customer support ticket by its ID
// @Tags         customer-ticket
// @Accept       json
// @Produce      json
// @Param        request  query     customer_ticket.GetByIdCustomerTicketQuery  true  "Ticket ID to retrieve"
// @success      200      {object}  utils.Result                                 "Ticket details"
// @Failure      400      {object}  utils.Result                                 "Validation error"
// @Failure      401      {object}  utils.Result                                 "unauthorized"
// @Failure      404      {object}  utils.Result                                 "Ticket not found"
// @Failure      500      {object}  utils.Result                                 "Internal server error"
// @Router       /customer-ticket [get]
// @Security BearerAuth
func (h *CustomerTicketHandler) GetByIdCustomerTicket(c *gin.Context) {
	var params customer_ticket.GetByIdCustomerTicketQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}

	result, err := h.usecase.GetByIdCustomerTicketQuery(&params)
	utils.HandleError(c, err)
	utils.HandleResponse(c, result)
}

// GetAllCustomerTicket godoc
// @Summary      Get all customer tickets
// @Description  Retrieves all support tickets for the current customer
// @Tags         customer-ticket
// @Accept       json
// @Produce      json
// @Param        request  query     customer_ticket.GetAllCustomerTicketQuery  true  "Query parameters"
// @success      200      {object}  utils.Result                                "List of customer tickets"
// @Failure      400      {object}  utils.Result                                "Validation error"
// @Failure      401      {object}  utils.Result                                "unauthorized"
// @Failure      500      {object}  utils.Result                                "Internal server error"
// @Router       /customer-ticket/all [get]
// @Security BearerAuth
func (h *CustomerTicketHandler) GetAllCustomerTicket(c *gin.Context) {
	var params customer_ticket.GetAllCustomerTicketQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}

	result, err := h.usecase.GetAllCustomerTicketQuery(&params)
	utils.HandleError(c, err)
	utils.HandleResponse(c, result)
}

// AdminGetAllCustomerTicket godoc
// @Summary      Admin: Get all customer tickets
// @Description  Admin endpoint to retrieve all customer support tickets across all customers
// @Tags         customer-ticket
// @Accept       json
// @Produce      json
// @Param        request  query     customer_ticket.AdminGetAllCustomerTicketQuery  true  "Query parameters"
// @success      200      {object}  utils.Result                                     "List of all customer tickets"
// @Failure      400      {object}  utils.Result                                     "Validation error"
// @Failure      401      {object}  utils.Result                                     "unauthorized"
// @Failure      403      {object}  utils.Result                                     "Forbidden - Admin access required"
// @Failure      500      {object}  utils.Result                                     "Internal server error"
// @Router       /customer-ticket/admin/all [get]
// @Security BearerAuth
func (h *CustomerTicketHandler) AdminGetAllCustomerTicket(c *gin.Context) {
	var params customer_ticket.AdminGetAllCustomerTicketQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}

	result, err := h.usecase.AdminGetAllCustomerTicketQuery(&params)
	utils.HandleError(c, err)
	utils.HandleResponse(c, result)
}
