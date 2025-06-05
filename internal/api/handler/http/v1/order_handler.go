package v1

import (
	utils2 "github.com/amirex128/new_site_builder/internal/api/utils"
	order2 "github.com/amirex128/new_site_builder/internal/application/dto/order"
	"github.com/amirex128/new_site_builder/internal/application/usecase/order"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	usecase   *orderusecase.OrderUsecase
	validator *utils2.ValidationHelper
}

func NewOrderHandler(usc *orderusecase.OrderUsecase) *OrderHandler {
	return &OrderHandler{
		usecase:   usc,
		validator: utils2.NewValidationHelper(),
	}
}

// CreateOrderRequest godoc
// @Summary      Create an order request
// @Description  Creates a new order request from the user's basket
// @Tags         order
// @Accept       json
// @Produce      json
// @Param        request  body      order.CreateOrderRequestCommand  true  "Order request information"
// @success      201      {object}  utils.Result                      "Created order request"
// @Failure      400      {object}  utils.Result                      "Validation error"
// @Failure      401      {object}  utils.Result                      "unauthorized"
// @Failure      500      {object}  utils.Result                      "Internal server error"
// @Router       /order [post]
// @Security BearerAuth
func (h *OrderHandler) CreateOrderRequest(c *gin.Context) {
	var params order2.CreateOrderRequestCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.CreateOrderRequestCommand(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// CreateOrderVerify godoc
// @Summary      Verify an order
// @Description  Verifies and finalizes an order after payment
// @Tags         order
// @Accept       json
// @Produce      json
// @Param        request  body      order.CreateOrderVerifyCommand  true  "Order verification information"
// @success      200      {object}  utils.Result                     "Verified order"
// @Failure      400      {object}  utils.Result                     "Validation error"
// @Failure      401      {object}  utils.Result                     "unauthorized"
// @Failure      404      {object}  utils.Result                     "Order not found"
// @Failure      500      {object}  utils.Result                     "Internal server error"
// @Router       /order/verify [post]
// @Security BearerAuth
func (h *OrderHandler) CreateOrderVerify(c *gin.Context) {
	var params order2.CreateOrderVerifyCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.CreateOrderVerifyCommand(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// GetAllOrderCustomer godoc
// @Summary      Get all customer orders
// @Description  Retrieves all orders for the authenticated customer
// @Tags         order
// @Accept       json
// @Produce      json
// @Param        request  query     order.GetAllOrderCustomerQuery  true  "Query parameters"
// @success      200      {object}  utils.Result                     "List of customer orders"
// @Failure      400      {object}  utils.Result                     "Validation error"
// @Failure      401      {object}  utils.Result                     "unauthorized"
// @Failure      500      {object}  utils.Result                     "Internal server error"
// @Router       /order/customer/all [get]
// @Security BearerAuth
func (h *OrderHandler) GetAllOrderCustomer(c *gin.Context) {
	var params order2.GetAllOrderCustomerQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.GetAllOrderCustomerQuery(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// GetOrderCustomerDetails godoc
// @Summary      Get customer order details
// @Description  Retrieves detailed information about a specific customer order
// @Tags         order
// @Accept       json
// @Produce      json
// @Param        request  query     order.GetOrderCustomerDetailsQuery  true  "Order ID to retrieve"
// @success      200      {object}  utils.Result                         "Order details"
// @Failure      400      {object}  utils.Result                         "Validation error"
// @Failure      401      {object}  utils.Result                         "unauthorized"
// @Failure      404      {object}  utils.Result                         "Order not found"
// @Failure      500      {object}  utils.Result                         "Internal server error"
// @Router       /order/customer/details [get]
// @Security BearerAuth
func (h *OrderHandler) GetOrderCustomerDetails(c *gin.Context) {
	var params order2.GetOrderCustomerDetailsQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.GetOrderCustomerDetailsQuery(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// GetAllOrderUser godoc
// @Summary      Get all user orders
// @Description  Retrieves all orders for the authenticated user (seller)
// @Tags         order
// @Accept       json
// @Produce      json
// @Param        request  query     order.GetAllOrderUserQuery  true  "Query parameters"
// @success      200      {object}  utils.Result                 "List of user orders"
// @Failure      400      {object}  utils.Result                 "Validation error"
// @Failure      401      {object}  utils.Result                 "unauthorized"
// @Failure      500      {object}  utils.Result                 "Internal server error"
// @Router       /order/user/all [get]
// @Security BearerAuth
func (h *OrderHandler) GetAllOrderUser(c *gin.Context) {
	var params order2.GetAllOrderUserQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.GetAllOrderUserQuery(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// GetOrderUserDetails godoc
// @Summary      Get user order details
// @Description  Retrieves detailed information about a specific user (seller) order
// @Tags         order
// @Accept       json
// @Produce      json
// @Param        request  query     order.GetOrderUserDetailsQuery  true  "Order ID to retrieve"
// @success      200      {object}  utils.Result                     "Order details"
// @Failure      400      {object}  utils.Result                     "Validation error"
// @Failure      401      {object}  utils.Result                     "unauthorized"
// @Failure      404      {object}  utils.Result                     "Order not found"
// @Failure      500      {object}  utils.Result                     "Internal server error"
// @Router       /order/user/details [get]
// @Security BearerAuth
func (h *OrderHandler) GetOrderUserDetails(c *gin.Context) {
	var params order2.GetOrderUserDetailsQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.GetOrderUserDetailsQuery(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// AdminGetAllOrderUser godoc
// @Summary      Admin: Get all orders
// @Description  Admin endpoint to retrieve all orders across all users
// @Tags         order
// @Accept       json
// @Produce      json
// @Param        request  query     order.AdminGetAllOrderUserQuery  true  "Query parameters"
// @success      200      {object}  utils.Result                      "List of all orders"
// @Failure      400      {object}  utils.Result                      "Validation error"
// @Failure      401      {object}  utils.Result                      "unauthorized"
// @Failure      403      {object}  utils.Result                      "Forbidden - Admin access required"
// @Failure      500      {object}  utils.Result                      "Internal server error"
// @Router       /order/admin/all [get]
// @Security BearerAuth
func (h *OrderHandler) AdminGetAllOrderUser(c *gin.Context) {
	var params order2.AdminGetAllOrderUserQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.AdminGetAllOrderUserQuery(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}
