package v1

import (
	"github.com/amirex128/new_site_builder/src/internal/api/utils"
	"github.com/amirex128/new_site_builder/src/internal/api/utils/resp"
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

func (h *TicketHandler) CreateTicket(c *gin.Context) {
	var params ticket.CreateTicketCommand
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.CreateTicketCommand(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Created(c, result)
}

func (h *TicketHandler) ReplayTicket(c *gin.Context) {
	var params ticket.ReplayTicketCommand
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.ReplayTicketCommand(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Updated(c, result)
}

func (h *TicketHandler) AdminReplayTicket(c *gin.Context) {
	var params ticket.AdminReplayTicketCommand
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.AdminReplayTicketCommand(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Updated(c, result)
}

func (h *TicketHandler) GetByIdTicket(c *gin.Context) {
	var params ticket.GetByIdTicketQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.GetByIdTicketQuery(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Retrieved(c, result)
}

func (h *TicketHandler) GetAllTicket(c *gin.Context) {
	var params ticket.GetAllTicketQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.GetAllTicketQuery(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Retrieved(c, result)
}

func (h *TicketHandler) AdminGetAllTicket(c *gin.Context) {
	var params ticket.AdminGetAllTicketQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.AdminGetAllTicketQuery(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Retrieved(c, result)
}
