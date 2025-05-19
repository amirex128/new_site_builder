package v1

import (
	"github.com/amirex128/new_site_builder/src/internal/api/utils"
	"github.com/amirex128/new_site_builder/src/internal/api/utils/resp"
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

func (h *CustomerTicketHandler) CreateCustomerTicket(c *gin.Context) {
	var params customer_ticket.CreateCustomerTicketCommand
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.CreateCustomerTicketCommand(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Created(c, result)
}

func (h *CustomerTicketHandler) ReplayCustomerTicket(c *gin.Context) {
	var params customer_ticket.ReplayCustomerTicketCommand
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.ReplayCustomerTicketCommand(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Updated(c, result)
}

func (h *CustomerTicketHandler) AdminReplayCustomerTicket(c *gin.Context) {
	var params customer_ticket.AdminReplayCustomerTicketCommand
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.AdminReplayCustomerTicketCommand(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Updated(c, result)
}

func (h *CustomerTicketHandler) GetByIdCustomerTicket(c *gin.Context) {
	var params customer_ticket.GetByIdCustomerTicketQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.GetByIdCustomerTicketQuery(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Retrieved(c, result)
}

func (h *CustomerTicketHandler) GetAllCustomerTicket(c *gin.Context) {
	var params customer_ticket.GetAllCustomerTicketQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.GetAllCustomerTicketQuery(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Retrieved(c, result)
}

func (h *CustomerTicketHandler) AdminGetAllCustomerTicket(c *gin.Context) {
	var params customer_ticket.AdminGetAllCustomerTicketQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.AdminGetAllCustomerTicketQuery(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Retrieved(c, result)
}
