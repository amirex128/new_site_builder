package v1

import (
	"github.com/amirex128/new_site_builder/src/internal/api/utils"
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

// CreateOrderRequest godoc
// @Summary      Create an order request
// @Description  Creates a new order request from the user's basket
// @Tags         order
// @Accept       json
// @Produce      json
// @Param        request  body      order.CreateOrderRequestCommand  true  "Order request information"
// @Success      201      {object}  utils.Result                      "Created order request"
// @Failure      400      {object}  utils.Result                      "Validation error"
// @Failure      401      {object}  utils.Result                      "Unauthorized"
// @Failure      500      {object}  utils.Result                      "Internal server error"
// @Router       /order [post]
// @Security BearerAuth
func (h *OrderHandler) CreateOrderRequest(c *gin.Context) {
	var params order.CreateOrderRequestCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}

	result, err := h.usecase.CreateOrderRequestCommand(&params)
	utils.HandleError(c, err)
	utils.HandleResponse(c, result)
}

// CreateOrderVerify godoc
// @Summary      Verify an order
// @Description  Verifies and finalizes an order after payment
// @Tags         order
// @Accept       json
// @Produce      json
// @Param        request  body      order.CreateOrderVerifyCommand  true  "Order verification information"
// @Success      200      {object}  utils.Result                     "Verified order"
// @Failure      400      {object}  utils.Result                     "Validation error"
// @Failure      401      {object}  utils.Result                     "Unauthorized"
// @Failure      404      {object}  utils.Result                     "Order not found"
// @Failure      500      {object}  utils.Result                     "Internal server error"
// @Router       /order/verify [post]
// @Security BearerAuth
func (h *OrderHandler) CreateOrderVerify(c *gin.Context) {
	var params order.CreateOrderVerifyCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}

	result, err := h.usecase.CreateOrderVerifyCommand(&params)
	utils.HandleError(c, err)
	utils.HandleResponse(c, result)
}

// GetAllOrderCustomer godoc
// @Summary      Get all customer orders
// @Description  Retrieves all orders for the authenticated customer
// @Tags         order
// @Accept       json
// @Produce      json
// @Param        request  query     order.GetAllOrderCustomerQuery  true  "Query parameters"
// @Success      200      {object}  utils.Result                     "List of customer orders"
// @Failure      400      {object}  utils.Result                     "Validation error"
// @Failure      401      {object}  utils.Result                     "Unauthorized"
// @Failure      500      {object}  utils.Result                     "Internal server error"
// @Router       /order/customer/all [get]
// @Security BearerAuth
func (h *OrderHandler) GetAllOrderCustomer(c *gin.Context) {
	var params order.GetAllOrderCustomerQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}

	result, err := h.usecase.GetAllOrderCustomerQuery(&params)
	utils.HandleError(c, err)
	utils.HandleResponse(c, result)
}

// GetOrderCustomerDetails godoc
// @Summary      Get customer order details
// @Description  Retrieves detailed information about a specific customer order
// @Tags         order
// @Accept       json
// @Produce      json
// @Param        request  query     order.GetOrderCustomerDetailsQuery  true  "Order ID to retrieve"
// @Success      200      {object}  utils.Result                         "Order details"
// @Failure      400      {object}  utils.Result                         "Validation error"
// @Failure      401      {object}  utils.Result                         "Unauthorized"
// @Failure      404      {object}  utils.Result                         "Order not found"
// @Failure      500      {object}  utils.Result                         "Internal server error"
// @Router       /order/customer/details [get]
// @Security BearerAuth
func (h *OrderHandler) GetOrderCustomerDetails(c *gin.Context) {
	var params order.GetOrderCustomerDetailsQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}

	result, err := h.usecase.GetOrderCustomerDetailsQuery(&params)
	utils.HandleError(c, err)
	utils.HandleResponse(c, result)
}

// GetAllOrderUser godoc
// @Summary      Get all user orders
// @Description  Retrieves all orders for the authenticated user (seller)
// @Tags         order
// @Accept       json
// @Produce      json
// @Param        request  query     order.GetAllOrderUserQuery  true  "Query parameters"
// @Success      200      {object}  utils.Result                 "List of user orders"
// @Failure      400      {object}  utils.Result                 "Validation error"
// @Failure      401      {object}  utils.Result                 "Unauthorized"
// @Failure      500      {object}  utils.Result                 "Internal server error"
// @Router       /order/user/all [get]
// @Security BearerAuth
func (h *OrderHandler) GetAllOrderUser(c *gin.Context) {
	var params order.GetAllOrderUserQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}

	result, err := h.usecase.GetAllOrderUserQuery(&params)
	utils.HandleError(c, err)
	utils.HandleResponse(c, result)
}

// GetOrderUserDetails godoc
// @Summary      Get user order details
// @Description  Retrieves detailed information about a specific user (seller) order
// @Tags         order
// @Accept       json
// @Produce      json
// @Param        request  query     order.GetOrderUserDetailsQuery  true  "Order ID to retrieve"
// @Success      200      {object}  utils.Result                     "Order details"
// @Failure      400      {object}  utils.Result                     "Validation error"
// @Failure      401      {object}  utils.Result                     "Unauthorized"
// @Failure      404      {object}  utils.Result                     "Order not found"
// @Failure      500      {object}  utils.Result                     "Internal server error"
// @Router       /order/user/details [get]
// @Security BearerAuth
func (h *OrderHandler) GetOrderUserDetails(c *gin.Context) {
	var params order.GetOrderUserDetailsQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}

	result, err := h.usecase.GetOrderUserDetailsQuery(&params)
	utils.HandleError(c, err)
	utils.HandleResponse(c, result)
}

// AdminGetAllOrderUser godoc
// @Summary      Admin: Get all orders
// @Description  Admin endpoint to retrieve all orders across all users
// @Tags         order
// @Accept       json
// @Produce      json
// @Param        request  query     order.AdminGetAllOrderUserQuery  true  "Query parameters"
// @Success      200      {object}  utils.Result                      "List of all orders"
// @Failure      400      {object}  utils.Result                      "Validation error"
// @Failure      401      {object}  utils.Result                      "Unauthorized"
// @Failure      403      {object}  utils.Result                      "Forbidden - Admin access required"
// @Failure      500      {object}  utils.Result                      "Internal server error"
// @Router       /order/admin/all [get]
// @Security BearerAuth
func (h *OrderHandler) AdminGetAllOrderUser(c *gin.Context) {
	var params order.AdminGetAllOrderUserQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}

	result, err := h.usecase.AdminGetAllOrderUserQuery(&params)
	utils.HandleError(c, err)
	utils.HandleResponse(c, result)
}
