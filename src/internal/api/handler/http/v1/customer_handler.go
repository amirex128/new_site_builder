package v1

import (
	"net/http"

	"github.com/amirex128/new_site_builder/src/internal/api/utils"
	"github.com/amirex128/new_site_builder/src/internal/api/utils/resp"
	"github.com/amirex128/new_site_builder/src/internal/application/dto/customer"
	customerusecase "github.com/amirex128/new_site_builder/src/internal/application/usecase/customer"
	"github.com/gin-gonic/gin"
)

type CustomerHandler struct {
	usecase   *customerusecase.CustomerUsecase
	validator *utils.ValidationHelper
}

func NewCustomerHandler(usc *customerusecase.CustomerUsecase) *CustomerHandler {
	return &CustomerHandler{
		usecase:   usc,
		validator: utils.NewValidationHelper(),
	}
}

func (h *CustomerHandler) UpdateProfileCustomer(c *gin.Context) {
	var params customer.UpdateProfileCustomerCommand
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.UpdateProfileCustomerCommand(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.InternalError().WithSystemMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp.Updated().WithData(result))
}

func (h *CustomerHandler) GetProfileCustomer(c *gin.Context) {
	var params customer.GetProfileCustomerQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.GetProfileCustomerQuery(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.InternalError().WithSystemMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp.Retrieved().WithData(result))
}

func (h *CustomerHandler) RegisterCustomer(c *gin.Context) {
	var params customer.RegisterCustomerCommand
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.RegisterCustomerCommand(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.InternalError().WithSystemMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp.Created().WithData(result))
}

func (h *CustomerHandler) LoginCustomer(c *gin.Context) {
	var params customer.LoginCustomerCommand
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.LoginCustomerCommand(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.InternalError().WithSystemMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp.Retrieved().WithData(result))
}

func (h *CustomerHandler) RequestVerifyAndForgetCustomer(c *gin.Context) {
	var params customer.RequestVerifyAndForgetCustomerCommand
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.RequestVerifyAndForgetCustomerCommand(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.InternalError().WithSystemMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp.Created().WithData(result))
}

func (h *CustomerHandler) VerifyCustomer(c *gin.Context) {
	var params customer.VerifyCustomerQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.VerifyCustomerQuery(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.InternalError().WithSystemMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp.Updated().WithData(result))
}

func (h *CustomerHandler) AdminGetAllCustomer(c *gin.Context) {
	var params customer.AdminGetAllCustomerQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.AdminGetAllCustomerQuery(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.InternalError().WithSystemMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp.Retrieved().WithData(result))
}
