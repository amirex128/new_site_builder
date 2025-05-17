package v1

import (
	"net/http"

	"github.com/amirex128/new_site_builder/src/internal/api/utils"
	"github.com/amirex128/new_site_builder/src/internal/api/utils/resp"
	"github.com/amirex128/new_site_builder/src/internal/application/dto/payment"
	paymentusecase "github.com/amirex128/new_site_builder/src/internal/application/usecase/payment"
	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	usecase   *paymentusecase.PaymentUsecase
	validator *utils.ValidationHelper
}

func NewPaymentHandler(usc *paymentusecase.PaymentUsecase) *PaymentHandler {
	return &PaymentHandler{
		usecase:   usc,
		validator: utils.NewValidationHelper(),
	}
}

func (h *PaymentHandler) VerifyPayment(c *gin.Context) {
	var params payment.VerifyPaymentCommand
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.VerifyPaymentCommand(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.InternalError().WithSystemMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp.Updated().WithData(result))
}

func (h *PaymentHandler) CreateOrUpdateGateway(c *gin.Context) {
	var params payment.CreateOrUpdateGatewayCommand
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.CreateOrUpdateGatewayCommand(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.InternalError().WithSystemMessage(err.Error()))
		return
	}

	// Since there's no clear ID field to determine if this is an update or a create,
	// we'll always respond with a generic success response
	c.JSON(http.StatusOK, resp.Success().WithData(result))
}

func (h *PaymentHandler) GetByIdGateway(c *gin.Context) {
	var params payment.GetByIdGatewayQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.GetByIdGatewayQuery(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.InternalError().WithSystemMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp.Retrieved().WithData(result))
}

func (h *PaymentHandler) AdminGetAllGateway(c *gin.Context) {
	var params payment.AdminGetAllGatewayQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.AdminGetAllGatewayQuery(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.InternalError().WithSystemMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp.Retrieved().WithData(result))
}

func (h *PaymentHandler) AdminGetAllPayment(c *gin.Context) {
	var params payment.AdminGetAllPaymentQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.AdminGetAllPaymentQuery(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.InternalError().WithSystemMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp.Retrieved().WithData(result))
}
