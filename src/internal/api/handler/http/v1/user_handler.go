package v1

import (
	"github.com/amirex128/new_site_builder/src/internal/api/utils"
	"github.com/amirex128/new_site_builder/src/internal/api/utils/resp"
	"github.com/amirex128/new_site_builder/src/internal/application/dto/user"
	userusecase "github.com/amirex128/new_site_builder/src/internal/application/usecase/user"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	usecase   *userusecase.UserUsecase
	validator *utils.ValidationHelper
}

func NewUserHandler(usc *userusecase.UserUsecase) *UserHandler {
	return &UserHandler{
		usecase:   usc,
		validator: utils.NewValidationHelper(),
	}
}

// UpdateProfileUser godoc
// @Summary      Update user profile
// @Description  Updates the profile information for the authenticated user
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        request  body      user.UpdateProfileUserCommand  true  "Updated profile information"
// @Success      200      {object}  resp.Result                    "Updated user profile"
// @Failure      400      {object}  resp.Result                    "Validation error"
// @Failure      401      {object}  resp.Result                    "Unauthorized"
// @Failure      500      {object}  resp.Result                    "Internal server error"
// @Router       /user/profile [put]
// @Security     BearerAuth
func (h *UserHandler) UpdateProfileUser(c *gin.Context) {
	var params user.UpdateProfileUserCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}

	result, err := h.usecase.UpdateProfileUserCommand(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Updated(c, result)
}

// GetProfileUser godoc
// @Summary      Get user profile
// @Description  Retrieves the profile information for the authenticated user
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        request  query     user.GetProfileUserQuery  true  "Query parameters"
// @Success      200      {object}  resp.Result               "User profile details"
// @Failure      400      {object}  resp.Result               "Validation error"
// @Failure      401      {object}  resp.Result               "Unauthorized"
// @Failure      404      {object}  resp.Result               "User not found"
// @Failure      500      {object}  resp.Result               "Internal server error"
// @Router       /user/profile [get]
// @Security     BearerAuth
func (h *UserHandler) GetProfileUser(c *gin.Context) {
	var params user.GetProfileUserQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}

	result, err := h.usecase.GetProfileUserQuery(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Retrieved(c, result)
}

// ChargeCreditRequestUser godoc
// @Summary      Request credit charge
// @Description  Creates a request to charge the user's account credit
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        request  body      user.ChargeCreditRequestUserCommand  true  "Credit charge request details"
// @Success      201      {object}  resp.Result                          "Created charge request"
// @Failure      400      {object}  resp.Result                          "Validation error"
// @Failure      401      {object}  resp.Result                          "Unauthorized"
// @Failure      500      {object}  resp.Result                          "Internal server error"
// @Router       /user/credit/charge [post]
// @Security     BearerAuth
func (h *UserHandler) ChargeCreditRequestUser(c *gin.Context) {
	var params user.ChargeCreditRequestUserCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}

	result, err := h.usecase.ChargeCreditRequestUserCommand(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Created(c, result)
}

// UpgradePlanRequestUser godoc
// @Summary      Request plan upgrade
// @Description  Creates a request to upgrade the user's subscription plan
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        request  body      user.UpgradePlanRequestUserCommand  true  "Plan upgrade request details"
// @Success      201      {object}  resp.Result                         "Created plan upgrade request"
// @Failure      400      {object}  resp.Result                         "Validation error"
// @Failure      401      {object}  resp.Result                         "Unauthorized"
// @Failure      500      {object}  resp.Result                         "Internal server error"
// @Router       /user/plan/upgrade [post]
// @Security     BearerAuth
func (h *UserHandler) UpgradePlanRequestUser(c *gin.Context) {
	var params user.UpgradePlanRequestUserCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}

	result, err := h.usecase.UpgradePlanRequestUserCommand(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Created(c, result)
}

// RegisterUser godoc
// @Summary      Register new user
// @Description  Creates a new user account
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        request  body      user.RegisterUserCommand  true  "User registration details"
// @Success      201      {object}  resp.Result              "Created user account"
// @Failure      400      {object}  resp.Result              "Validation error"
// @Failure      500      {object}  resp.Result              "Internal server error"
// @Router       /user/register [post]
func (h *UserHandler) RegisterUser(c *gin.Context) {
	var params user.RegisterUserCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}

	result, err := h.usecase.RegisterUserCommand(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Created(c, result)
}

// LoginUser godoc
// @Summary      User login
// @Description  Authenticates a user and returns a token
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        request  body      user.LoginUserCommand  true  "User login credentials"
// @Success      200      {object}  resp.Result           "Authentication token and user details"
// @Failure      400      {object}  resp.Result           "Validation error"
// @Failure      401      {object}  resp.Result           "Invalid credentials"
// @Failure      500      {object}  resp.Result           "Internal server error"
// @Router       /user/login [post]
func (h *UserHandler) LoginUser(c *gin.Context) {
	var params user.LoginUserCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}

	result, err := h.usecase.LoginUserCommand(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Retrieved(c, result)
}

// RequestVerifyAndForgetUser godoc
// @Summary      Request verification or password reset
// @Description  Sends a verification code or password reset link to the user's email
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        request  body      user.RequestVerifyAndForgetUserCommand  true  "Email verification or password reset request"
// @Success      201      {object}  resp.Result                             "Verification or reset request created"
// @Failure      400      {object}  resp.Result                             "Validation error"
// @Failure      404      {object}  resp.Result                             "User not found"
// @Failure      500      {object}  resp.Result                             "Internal server error"
// @Router       /user/verify-forget/request [post]
func (h *UserHandler) RequestVerifyAndForgetUser(c *gin.Context) {
	var params user.RequestVerifyAndForgetUserCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}

	result, err := h.usecase.RequestVerifyAndForgetUserCommand(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Created(c, result)
}

// VerifyUser godoc
// @Summary      Verify user
// @Description  Verifies a user's email or resets password using a verification code
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        request  query     user.VerifyUserQuery  true  "Verification details"
// @Success      200      {object}  resp.Result           "User verified or password reset"
// @Failure      400      {object}  resp.Result           "Validation error"
// @Failure      404      {object}  resp.Result           "User not found or invalid code"
// @Failure      500      {object}  resp.Result           "Internal server error"
// @Router       /user/verify [get]
func (h *UserHandler) VerifyUser(c *gin.Context) {
	var params user.VerifyUserQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}

	result, err := h.usecase.VerifyUserQuery(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Updated(c, result)
}

// AdminGetAllUser godoc
// @Summary      Admin: Get all users
// @Description  Admin endpoint to retrieve all users in the system
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        request  query     user.AdminGetAllUserQuery  true  "Query parameters"
// @Success      200      {object}  resp.Result               "List of all users"
// @Failure      400      {object}  resp.Result               "Validation error"
// @Failure      401      {object}  resp.Result               "Unauthorized"
// @Failure      403      {object}  resp.Result               "Forbidden - Admin access required"
// @Failure      500      {object}  resp.Result               "Internal server error"
// @Router       /user/admin/all [get]
// @Security     BearerAuth
func (h *UserHandler) AdminGetAllUser(c *gin.Context) {
	var params user.AdminGetAllUserQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}

	result, err := h.usecase.AdminGetAllUserQuery(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Retrieved(c, result)
}
