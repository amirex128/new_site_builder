package v1

import (
	utils2 "github.com/amirex128/new_site_builder/internal/api/utils"
	payment2 "github.com/amirex128/new_site_builder/internal/application/dto/payment"
	"github.com/amirex128/new_site_builder/internal/application/usecase/payment"
	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	usecase   *paymentusecase.PaymentUsecase
	validator *utils2.ValidationHelper
}

func NewPaymentHandler(usc *paymentusecase.PaymentUsecase) *PaymentHandler {
	return &PaymentHandler{
		usecase:   usc,
		validator: utils2.NewValidationHelper(),
	}
}

// VerifyPayment godoc
// @Summary      Verify payment
// @Description  Verifies a payment transaction after it's processed by the payment gateway
// @Tags         payment
// @Accept       json
// @Produce      json
// @Param        request  body      payment.VerifyPaymentCommand  true  "Payment verification information"
// @success      200      {object}  utils.Result                   "Verified payment"
// @Failure      400      {object}  utils.Result                   "Validation error"
// @Failure      401      {object}  utils.Result                   "unauthorized"
// @Failure      404      {object}  utils.Result                   "Payment not found"
// @Failure      500      {object}  utils.Result                   "Internal server error"
// @Router       /payment/verify [post]
// @Security BearerAuth
func (h *PaymentHandler) VerifyPayment(c *gin.Context) {
	var params payment2.VerifyPaymentCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.VerifyPaymentCommand(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// RequestGateway godoc
// @Summary      Request payment gateway
// @Description  Initiates a payment request to a payment gateway
// @Tags         payment
// @Accept       json
// @Produce      json
// @Param        request  body      payment.RequestGatewayCommand  true  "Payment request information"
// @success      200      {object}  utils.Result                    "Payment gateway request result"
// @Failure      400      {object}  utils.Result                    "Validation error"
// @Failure      401      {object}  utils.Result                    "unauthorized"
// @Failure      500      {object}  utils.Result                    "Internal server error"
// @Router       /payment/request [post]
// @Security BearerAuth
func (h *PaymentHandler) RequestGateway(c *gin.Context) {
	var params payment2.RequestGatewayCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.RequestGatewayCommand(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// CreateOrUpdateGateway godoc
// @Summary      Create or update payment gateway
// @Description  Creates a new payment gateway or updates an existing one
// @Tags         gateway
// @Accept       json
// @Produce      json
// @Param        request  body      payment.CreateOrUpdateGatewayCommand  true  "Gateway information"
// @success      200      {object}  utils.Result                           "Created or updated gateway"
// @Failure      400      {object}  utils.Result                           "Validation error"
// @Failure      401      {object}  utils.Result                           "unauthorized"
// @Failure      500      {object}  utils.Result                           "Internal server error"
// @Router       /gateway [post]
// @Security BearerAuth
func (h *PaymentHandler) CreateOrUpdateGateway(c *gin.Context) {
	var params payment2.CreateOrUpdateGatewayCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.CreateOrUpdateGatewayCommand(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// GetByIdGateway godoc
// @Summary      Get payment gateway by ID
// @Description  Retrieves a specific payment gateway by its ID
// @Tags         gateway
// @Accept       json
// @Produce      json
// @Param        request  query     payment.GetByIdGatewayQuery  true  "Gateway ID to retrieve"
// @success      200      {object}  utils.Result                  "Gateway details"
// @Failure      400      {object}  utils.Result                  "Validation error"
// @Failure      401      {object}  utils.Result                  "unauthorized"
// @Failure      404      {object}  utils.Result                  "Gateway not found"
// @Failure      500      {object}  utils.Result                  "Internal server error"
// @Router       /gateway [get]
// @Security BearerAuth
func (h *PaymentHandler) GetByIdGateway(c *gin.Context) {
	var params payment2.GetByIdGatewayQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.GetByIdGatewayQuery(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// AdminGetAllGateway godoc
// @Summary      Admin: Get all payment gateways
// @Description  Admin endpoint to retrieve all payment gateways
// @Tags         gateway
// @Accept       json
// @Produce      json
// @Param        request  query     payment.AdminGetAllGatewayQuery  true  "Query parameters"
// @success      200      {object}  utils.Result                      "List of all gateways"
// @Failure      400      {object}  utils.Result                      "Validation error"
// @Failure      401      {object}  utils.Result                      "unauthorized"
// @Failure      403      {object}  utils.Result                      "Forbidden - Admin access required"
// @Failure      500      {object}  utils.Result                      "Internal server error"
// @Router       /gateway/admin/all [get]
// @Security BearerAuth
func (h *PaymentHandler) AdminGetAllGateway(c *gin.Context) {
	var params payment2.AdminGetAllGatewayQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.AdminGetAllGatewayQuery(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// AdminGetAllPayment godoc
// @Summary      Admin: Get all payments
// @Description  Admin endpoint to retrieve all payment transactions
// @Tags         payment
// @Accept       json
// @Produce      json
// @Param        request  query     payment.AdminGetAllPaymentQuery  true  "Query parameters"
// @success      200      {object}  utils.Result                      "List of all payments"
// @Failure      400      {object}  utils.Result                      "Validation error"
// @Failure      401      {object}  utils.Result                      "unauthorized"
// @Failure      403      {object}  utils.Result                      "Forbidden - Admin access required"
// @Failure      500      {object}  utils.Result                      "Internal server error"
// @Router       /payment/admin/all [get]
// @Security BearerAuth
func (h *PaymentHandler) AdminGetAllPayment(c *gin.Context) {
	var params payment2.AdminGetAllPaymentQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.AdminGetAllPaymentQuery(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}
