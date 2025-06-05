package v1

import (
	utils2 "github.com/amirex128/new_site_builder/internal/api/utils"
	customer2 "github.com/amirex128/new_site_builder/internal/application/dto/customer"
	"github.com/amirex128/new_site_builder/internal/application/usecase/customer"
	"github.com/gin-gonic/gin"
)

type CustomerHandler struct {
	usecase   *customerusecase.CustomerUsecase
	validator *utils2.ValidationHelper
}

func NewCustomerHandler(usc *customerusecase.CustomerUsecase) *CustomerHandler {
	return &CustomerHandler{
		usecase:   usc,
		validator: utils2.NewValidationHelper(),
	}
}

// UpdateProfileCustomer godoc
// @Summary      Update customer profile
// @Description  Updates the authenticated customer's profile information
// @Tags         customer
// @Accept       json
// @Produce      json
// @Param        request  body      customer.UpdateProfileCustomerCommand  true  "Updated profile information"
// @success      200      {object}  utils.Result                            "Updated profile"
// @Failure      400      {object}  utils.Result                            "Validation error"
// @Failure      401      {object}  utils.Result                            "unauthorized"
// @Failure      500      {object}  utils.Result                            "Internal server error"
// @Router       /customer/profile [put]
// @Security BearerAuth
func (h *CustomerHandler) UpdateProfileCustomer(c *gin.Context) {
	var params customer2.UpdateProfileCustomerCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.UpdateProfileCustomerCommand(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// GetProfileCustomer godoc
// @Summary      Get customer profile
// @Description  Retrieves the authenticated customer's profile information
// @Tags         customer
// @Accept       json
// @Produce      json
// @Param        request  query     customer.GetProfileCustomerQuery  true  "Query parameters"
// @success      200      {object}  utils.Result                       "Customer profile"
// @Failure      400      {object}  utils.Result                       "Validation error"
// @Failure      401      {object}  utils.Result                       "unauthorized"
// @Failure      404      {object}  utils.Result                       "Profile not found"
// @Failure      500      {object}  utils.Result                       "Internal server error"
// @Router       /customer/profile [get]
// @Security BearerAuth
func (h *CustomerHandler) GetProfileCustomer(c *gin.Context) {
	var params customer2.GetProfileCustomerQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.GetProfileCustomerQuery(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// RegisterCustomer godoc
// @Summary      Register new customer
// @Description  Registers a new customer account with the provided information
// @Tags         customer
// @Accept       json
// @Produce      json
// @Param        request  body      customer.RegisterCustomerCommand  true  "Customer registration information"
// @success      201      {object}  utils.Result                       "Registered customer"
// @Failure      400      {object}  utils.Result                       "Validation error"
// @Failure      409      {object}  utils.Result                       "Email already exists"
// @Failure      500      {object}  utils.Result                       "Internal server error"
// @Router       /customer/register [post]
// @Security BearerAuth
func (h *CustomerHandler) RegisterCustomer(c *gin.Context) {
	var params customer2.RegisterCustomerCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.RegisterCustomerCommand(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// LoginCustomer godoc
// @Summary      Customer login
// @Description  Authenticates a customer and returns an access token
// @Tags         customer
// @Accept       json
// @Produce      json
// @Param        request  body      customer.LoginCustomerCommand  true  "Login credentials"
// @success      200      {object}  utils.Result                    "Authentication token"
// @Failure      400      {object}  utils.Result                    "Validation error"
// @Failure      401      {object}  utils.Result                    "Invalid credentials"
// @Failure      500      {object}  utils.Result                    "Internal server error"
// @Router       /customer/login [post]
// @Security BearerAuth
func (h *CustomerHandler) LoginCustomer(c *gin.Context) {
	var params customer2.LoginCustomerCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.LoginCustomerCommand(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// RequestVerifyAndForgetCustomer godoc
// @Summary      Request verification or password reset
// @Description  Sends verification email or password reset link to the customer's email
// @Tags         customer
// @Accept       json
// @Produce      json
// @Param        request  body      customer.RequestVerifyAndForgetCustomerCommand  true  "Email information"
// @success      201      {object}  utils.Result                                     "Email sent confirmation"
// @Failure      400      {object}  utils.Result                                     "Validation error"
// @Failure      404      {object}  utils.Result                                     "Email not found"
// @Failure      500      {object}  utils.Result                                     "Internal server error"
// @Router       /customer/verify-forget [post]
// @Security BearerAuth
func (h *CustomerHandler) RequestVerifyAndForgetCustomer(c *gin.Context) {
	var params customer2.RequestVerifyAndForgetCustomerCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.RequestVerifyAndForgetCustomerCommand(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// VerifyCustomer godoc
// @Summary      Verify customer email
// @Description  Verifies a customer's email address using the verification token
// @Tags         customer
// @Accept       json
// @Produce      json
// @Param        request  query     customer.VerifyCustomerQuery  true  "Verification token"
// @success      200      {object}  utils.Result                   "Verification successful"
// @Failure      400      {object}  utils.Result                   "Validation error"
// @Failure      401      {object}  utils.Result                   "Invalid token"
// @Failure      500      {object}  utils.Result                   "Internal server error"
// @Router       /customer/verify [get]
// @Security BearerAuth
func (h *CustomerHandler) VerifyCustomer(c *gin.Context) {
	var params customer2.VerifyCustomerQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.VerifyCustomerQuery(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}

// AdminGetAllCustomer godoc
// @Summary      Admin: Get all customers
// @Description  Admin endpoint to retrieve all customers with additional information
// @Tags         customer
// @Accept       json
// @Produce      json
// @Param        request  query     customer.AdminGetAllCustomerQuery  true  "Query parameters"
// @success      200      {object}  utils.Result                        "List of all customers"
// @Failure      400      {object}  utils.Result                        "Validation error"
// @Failure      401      {object}  utils.Result                        "unauthorized"
// @Failure      403      {object}  utils.Result                        "Forbidden - Admin access required"
// @Failure      500      {object}  utils.Result                        "Internal server error"
// @Router       /customer/admin/all [get]
// @Security BearerAuth
func (h *CustomerHandler) AdminGetAllCustomer(c *gin.Context) {
	var params customer2.AdminGetAllCustomerQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}
	h.usecase.SetContext(c)
	result, err := h.usecase.AdminGetAllCustomerQuery(&params)
	utils2.HandleError(c, err)
	utils2.HandleResponse(c, result)
}
