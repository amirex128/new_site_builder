package v1

import (
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

// VerifyPayment godoc
// @Summary      Verify payment
// @Description  Verifies a payment transaction after it's processed by the payment gateway
// @Tags         payment
// @Accept       json
// @Produce      json
// @Param        request  body      payment.VerifyPaymentCommand  true  "Payment verification information"
// @Success      200      {object}  resp.Result                   "Verified payment"
// @Failure      400      {object}  resp.Result                   "Validation error"
// @Failure      401      {object}  resp.Result                   "Unauthorized"
// @Failure      404      {object}  resp.Result                   "Payment not found"
// @Failure      500      {object}  resp.Result                   "Internal server error"
// @Router       /payment/verify [post]
// @Security     BearerAuth
func (h *PaymentHandler) VerifyPayment(c *gin.Context) {
	var params payment.VerifyPaymentCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}

	result, err := h.usecase.VerifyPaymentCommand(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Updated(c, result)
}

// RequestGateway godoc
// @Summary      Request payment gateway
// @Description  Initiates a payment request to a payment gateway
// @Tags         payment
// @Accept       json
// @Produce      json
// @Param        request  body      payment.RequestGatewayCommand  true  "Payment request information"
// @Success      200      {object}  resp.Result                    "Payment gateway request result"
// @Failure      400      {object}  resp.Result                    "Validation error"
// @Failure      401      {object}  resp.Result                    "Unauthorized"
// @Failure      500      {object}  resp.Result                    "Internal server error"
// @Router       /payment/request [post]
// @Security     BearerAuth
func (h *PaymentHandler) RequestGateway(c *gin.Context) {
	var params payment.RequestGatewayCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}

	result, err := h.usecase.RequestGatewayCommand(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.OK(c, result)
}

// CreateOrUpdateGateway godoc
// @Summary      Create or update payment gateway
// @Description  Creates a new payment gateway or updates an existing one
// @Tags         gateway
// @Accept       json
// @Produce      json
// @Param        request  body      payment.CreateOrUpdateGatewayCommand  true  "Gateway information"
// @Success      200      {object}  resp.Result                           "Created or updated gateway"
// @Failure      400      {object}  resp.Result                           "Validation error"
// @Failure      401      {object}  resp.Result                           "Unauthorized"
// @Failure      500      {object}  resp.Result                           "Internal server error"
// @Router       /gateway [post]
// @Security     BearerAuth
func (h *PaymentHandler) CreateOrUpdateGateway(c *gin.Context) {
	var params payment.CreateOrUpdateGatewayCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}

	result, err := h.usecase.CreateOrUpdateGatewayCommand(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	// Since there's no clear ID field to determine if this is an update or a create,
	// we'll always respond with a generic success response
	resp.OK(c, result)
}

// GetByIdGateway godoc
// @Summary      Get payment gateway by ID
// @Description  Retrieves a specific payment gateway by its ID
// @Tags         gateway
// @Accept       json
// @Produce      json
// @Param        request  query     payment.GetByIdGatewayQuery  true  "Gateway ID to retrieve"
// @Success      200      {object}  resp.Result                  "Gateway details"
// @Failure      400      {object}  resp.Result                  "Validation error"
// @Failure      401      {object}  resp.Result                  "Unauthorized"
// @Failure      404      {object}  resp.Result                  "Gateway not found"
// @Failure      500      {object}  resp.Result                  "Internal server error"
// @Router       /gateway [get]
// @Security     BearerAuth
func (h *PaymentHandler) GetByIdGateway(c *gin.Context) {
	var params payment.GetByIdGatewayQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}

	result, err := h.usecase.GetByIdGatewayQuery(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Retrieved(c, result)
}

// AdminGetAllGateway godoc
// @Summary      Admin: Get all payment gateways
// @Description  Admin endpoint to retrieve all payment gateways
// @Tags         gateway
// @Accept       json
// @Produce      json
// @Param        request  query     payment.AdminGetAllGatewayQuery  true  "Query parameters"
// @Success      200      {object}  resp.Result                      "List of all gateways"
// @Failure      400      {object}  resp.Result                      "Validation error"
// @Failure      401      {object}  resp.Result                      "Unauthorized"
// @Failure      403      {object}  resp.Result                      "Forbidden - Admin access required"
// @Failure      500      {object}  resp.Result                      "Internal server error"
// @Router       /gateway/admin/all [get]
// @Security     BearerAuth
func (h *PaymentHandler) AdminGetAllGateway(c *gin.Context) {
	var params payment.AdminGetAllGatewayQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}

	result, err := h.usecase.AdminGetAllGatewayQuery(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Retrieved(c, result)
}

// AdminGetAllPayment godoc
// @Summary      Admin: Get all payments
// @Description  Admin endpoint to retrieve all payment transactions
// @Tags         payment
// @Accept       json
// @Produce      json
// @Param        request  query     payment.AdminGetAllPaymentQuery  true  "Query parameters"
// @Success      200      {object}  resp.Result                      "List of all payments"
// @Failure      400      {object}  resp.Result                      "Validation error"
// @Failure      401      {object}  resp.Result                      "Unauthorized"
// @Failure      403      {object}  resp.Result                      "Forbidden - Admin access required"
// @Failure      500      {object}  resp.Result                      "Internal server error"
// @Router       /payment/admin/all [get]
// @Security     BearerAuth
func (h *PaymentHandler) AdminGetAllPayment(c *gin.Context) {
	var params payment.AdminGetAllPaymentQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}

	result, err := h.usecase.AdminGetAllPaymentQuery(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Retrieved(c, result)
}
