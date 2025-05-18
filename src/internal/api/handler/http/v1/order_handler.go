package v1

import (
	"net/http"

	"github.com/amirex128/new_site_builder/src/internal/api/utils"
	"github.com/amirex128/new_site_builder/src/internal/api/utils/resp"
	"github.com/amirex128/new_site_builder/src/internal/application/dto/order"
	orderusecase "github.com/amirex128/new_site_builder/src/internal/application/usecase/order"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	usecase   *orderusecase.OrderUsecase
	validator *utils.ValidationHelper
}

func NewOrderHandler(usc *orderusecase.OrderUsecase) *OrderHandler {
	return &OrderHandler{
		usecase:   usc,
		validator: utils.NewValidationHelper(),
	}
}

func (h *OrderHandler) CreateOrderRequest(c *gin.Context) {
	var params order.CreateOrderRequestCommand
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.CreateOrderRequestCommand(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.InternalError().WithSystemMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp.Created().WithData(result))
}

func (h *OrderHandler) CreateOrderVerify(c *gin.Context) {
	var params order.CreateOrderVerifyCommand
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.CreateOrderVerifyCommand(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.InternalError().WithSystemMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp.Success().WithData(result))
}

func (h *OrderHandler) GetAllOrderCustomer(c *gin.Context) {
	var params order.GetAllOrderCustomerQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.GetAllOrderCustomerQuery(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.InternalError().WithSystemMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp.Retrieved().WithData(result))
}

func (h *OrderHandler) GetOrderCustomerDetails(c *gin.Context) {
	var params order.GetOrderCustomerDetailsQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.GetOrderCustomerDetailsQuery(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.InternalError().WithSystemMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp.Retrieved().WithData(result))
}

func (h *OrderHandler) GetAllOrderUser(c *gin.Context) {
	var params order.GetAllOrderUserQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.GetAllOrderUserQuery(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.InternalError().WithSystemMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp.Retrieved().WithData(result))
}

func (h *OrderHandler) GetOrderUserDetails(c *gin.Context) {
	var params order.GetOrderUserDetailsQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.GetOrderUserDetailsQuery(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.InternalError().WithSystemMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp.Retrieved().WithData(result))
}

func (h *OrderHandler) AdminGetAllOrderUser(c *gin.Context) {
	var params order.AdminGetAllOrderUserQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.AdminGetAllOrderUserQuery(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.InternalError().WithSystemMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp.Retrieved().WithData(result))
}
