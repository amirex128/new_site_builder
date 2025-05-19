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

func (h *UserHandler) UpdateProfileUser(c *gin.Context) {
	var params user.UpdateProfileUserCommand
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.SetContext(c).UpdateProfileUserCommand(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Updated(c, result)
}

func (h *UserHandler) GetProfileUser(c *gin.Context) {
	var params user.GetProfileUserQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.SetContext(c).GetProfileUserQuery(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Retrieved(c, result)
}

func (h *UserHandler) ChargeCreditRequestUser(c *gin.Context) {
	var params user.ChargeCreditRequestUserCommand
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.SetContext(c).ChargeCreditRequestUserCommand(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Created(c, result)
}

func (h *UserHandler) UpgradePlanRequestUser(c *gin.Context) {
	var params user.UpgradePlanRequestUserCommand
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.SetContext(c).UpgradePlanRequestUserCommand(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Created(c, result)
}

func (h *UserHandler) RegisterUser(c *gin.Context) {
	var params user.RegisterUserCommand
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.SetContext(c).RegisterUserCommand(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Created(c, result)
}

func (h *UserHandler) LoginUser(c *gin.Context) {
	var params user.LoginUserCommand
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.SetContext(c).LoginUserCommand(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Retrieved(c, result)
}

func (h *UserHandler) RequestVerifyAndForgetUser(c *gin.Context) {
	var params user.RequestVerifyAndForgetUserCommand
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.SetContext(c).RequestVerifyAndForgetUserCommand(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Created(c, result)
}

func (h *UserHandler) VerifyUser(c *gin.Context) {
	var params user.VerifyUserQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.SetContext(c).VerifyUserQuery(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Updated(c, result)
}

func (h *UserHandler) AdminGetAllUser(c *gin.Context) {
	var params user.AdminGetAllUserQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.SetContext(c).AdminGetAllUserQuery(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Retrieved(c, result)
}
