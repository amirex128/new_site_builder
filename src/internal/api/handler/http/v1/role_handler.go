package v1

import (
	"github.com/amirex128/new_site_builder/src/internal/api/utils"
	"github.com/amirex128/new_site_builder/src/internal/api/utils/resp"
	"github.com/amirex128/new_site_builder/src/internal/application/dto/role"
	roleusecase "github.com/amirex128/new_site_builder/src/internal/application/usecase/role"
	"github.com/gin-gonic/gin"
)

type RoleHandler struct {
	usecase   *roleusecase.RoleUsecase
	validator *utils.ValidationHelper
}

func NewRoleHandler(usc *roleusecase.RoleUsecase) *RoleHandler {
	return &RoleHandler{
		usecase:   usc,
		validator: utils.NewValidationHelper(),
	}
}

func (h *RoleHandler) CreateRole(c *gin.Context) {
	var params role.CreateRoleCommand
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.CreateRoleCommand(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Created(c, result)
}

func (h *RoleHandler) UpdateRole(c *gin.Context) {
	var params role.UpdateRoleCommand
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.UpdateRoleCommand(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Updated(c, result)
}

func (h *RoleHandler) SetRoleToCustomer(c *gin.Context) {
	var params role.SetRoleToCustomerCommand
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.SetRoleToCustomerCommand(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Updated(c, result)
}

func (h *RoleHandler) SetRoleToUser(c *gin.Context) {
	var params role.SetRoleToUserCommand
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.SetRoleToUserCommand(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Updated(c, result)
}

func (h *RoleHandler) SetRoleToPlan(c *gin.Context) {
	var params role.SetRoleToPlanCommand
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.SetRoleToPlanCommand(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Updated(c, result)
}

func (h *RoleHandler) GetAllPermission(c *gin.Context) {
	var params role.GetAllPermissionQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.GetAllPermissionQuery(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Retrieved(c, result)
}

func (h *RoleHandler) GetAllRole(c *gin.Context) {
	var params role.GetAllRoleQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.GetAllRoleQuery(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Retrieved(c, result)
}

func (h *RoleHandler) GetRolePermissions(c *gin.Context) {
	var params role.GetRolePermissionsQuery
	if !h.validator.ValidateRequest(c, &params) {
		return
	}

	result, err := h.usecase.GetRolePermissionsQuery(&params)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Retrieved(c, result)
}
