package v1

import (
	"github.com/amirex128/new_site_builder/src/internal/api/utils"
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

// CreateRole godoc
// @Summary      Create a new role
// @Description  Creates a new role with specified permissions
// @Tags         role
// @Accept       json
// @Produce      json
// @Param        request  body      role.CreateRoleCommand  true  "Role information"
// @success      201      {object}  utils.Result            "Created role"
// @Failure      400      {object}  utils.Result            "Validation error"
// @Failure      401      {object}  utils.Result            "unauthorized"
// @Failure      403      {object}  utils.Result            "Forbidden - Admin access required"
// @Failure      500      {object}  utils.Result            "Internal server error"
// @Router       /role [post]
// @Security BearerAuth
func (h *RoleHandler) CreateRole(c *gin.Context) {
	var params role.CreateRoleCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}

	result, err := h.usecase.CreateRoleCommand(&params)
	utils.HandleError(c, err)
	utils.HandleResponse(c, result)
}

// UpdateRole godoc
// @Summary      Update an existing role
// @Description  Updates an existing role's name and permissions
// @Tags         role
// @Accept       json
// @Produce      json
// @Param        request  body      role.UpdateRoleCommand  true  "Updated role information"
// @success      200      {object}  utils.Result            "Updated role"
// @Failure      400      {object}  utils.Result            "Validation error"
// @Failure      401      {object}  utils.Result            "unauthorized"
// @Failure      403      {object}  utils.Result            "Forbidden - Admin access required"
// @Failure      404      {object}  utils.Result            "Role not found"
// @Failure      500      {object}  utils.Result            "Internal server error"
// @Router       /role [put]
// @Security BearerAuth
func (h *RoleHandler) UpdateRole(c *gin.Context) {
	var params role.UpdateRoleCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}

	result, err := h.usecase.UpdateRoleCommand(&params)
	utils.HandleError(c, err)
	utils.HandleResponse(c, result)
}

// SetRoleToCustomer godoc
// @Summary      Assign role to customer
// @Description  Assigns a specific role to a customer
// @Tags         role
// @Accept       json
// @Produce      json
// @Param        request  body      role.SetRoleToCustomerCommand  true  "Role and customer assignment details"
// @success      200      {object}  utils.Result                   "Role assigned to customer"
// @Failure      400      {object}  utils.Result                   "Validation error"
// @Failure      401      {object}  utils.Result                   "unauthorized"
// @Failure      403      {object}  utils.Result                   "Forbidden - Admin access required"
// @Failure      404      {object}  utils.Result                   "Role or customer not found"
// @Failure      500      {object}  utils.Result                   "Internal server error"
// @Router       /role/customer [put]
// @Security BearerAuth
func (h *RoleHandler) SetRoleToCustomer(c *gin.Context) {
	var params role.SetRoleToCustomerCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}

	result, err := h.usecase.SetRoleToCustomerCommand(&params)
	utils.HandleError(c, err)
	utils.HandleResponse(c, result)
}

// SetRoleToUser godoc
// @Summary      Assign role to user
// @Description  Assigns a specific role to a user
// @Tags         role
// @Accept       json
// @Produce      json
// @Param        request  body      role.SetRoleToUserCommand  true  "Role and user assignment details"
// @success      200      {object}  utils.Result               "Role assigned to user"
// @Failure      400      {object}  utils.Result               "Validation error"
// @Failure      401      {object}  utils.Result               "unauthorized"
// @Failure      403      {object}  utils.Result               "Forbidden - Admin access required"
// @Failure      404      {object}  utils.Result               "Role or user not found"
// @Failure      500      {object}  utils.Result               "Internal server error"
// @Router       /role/user [put]
// @Security BearerAuth
func (h *RoleHandler) SetRoleToUser(c *gin.Context) {
	var params role.SetRoleToUserCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}

	result, err := h.usecase.SetRoleToUserCommand(&params)
	utils.HandleError(c, err)
	utils.HandleResponse(c, result)
}

// SetRoleToPlan godoc
// @Summary      Assign role to plan
// @Description  Assigns a specific role to a subscription plan
// @Tags         role
// @Accept       json
// @Produce      json
// @Param        request  body      role.SetRoleToPlanCommand  true  "Role and plan assignment details"
// @success      200      {object}  utils.Result               "Role assigned to plan"
// @Failure      400      {object}  utils.Result               "Validation error"
// @Failure      401      {object}  utils.Result               "unauthorized"
// @Failure      403      {object}  utils.Result               "Forbidden - Admin access required"
// @Failure      404      {object}  utils.Result               "Role or plan not found"
// @Failure      500      {object}  utils.Result               "Internal server error"
// @Router       /role/plan [put]
// @Security BearerAuth
func (h *RoleHandler) SetRoleToPlan(c *gin.Context) {
	var params role.SetRoleToPlanCommand
	if !h.validator.ValidateCommand(c, &params) {
		return
	}

	result, err := h.usecase.SetRoleToPlanCommand(&params)
	utils.HandleError(c, err)
	utils.HandleResponse(c, result)
}

// GetAllPermission godoc
// @Summary      Get all permissions
// @Description  Retrieves all available permissions in the system
// @Tags         role
// @Accept       json
// @Produce      json
// @Param        request  query     role.GetAllPermissionQuery  true  "Query parameters"
// @success      200      {object}  utils.Result                "List of all permissions"
// @Failure      400      {object}  utils.Result                "Validation error"
// @Failure      401      {object}  utils.Result                "unauthorized"
// @Failure      500      {object}  utils.Result                "Internal server error"
// @Router       /role/permission/all [get]
// @Security BearerAuth
func (h *RoleHandler) GetAllPermission(c *gin.Context) {
	var params role.GetAllPermissionQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}

	result, err := h.usecase.GetAllPermissionQuery(&params)
	utils.HandleError(c, err)
	utils.HandleResponse(c, result)
}

// GetAllRole godoc
// @Summary      Get all roles
// @Description  Retrieves all roles in the system
// @Tags         role
// @Accept       json
// @Produce      json
// @Param        request  query     role.GetAllRoleQuery  true  "Query parameters"
// @success      200      {object}  utils.Result          "List of all roles"
// @Failure      400      {object}  utils.Result          "Validation error"
// @Failure      401      {object}  utils.Result          "unauthorized"
// @Failure      500      {object}  utils.Result          "Internal server error"
// @Router       /role/all [get]
// @Security BearerAuth
func (h *RoleHandler) GetAllRole(c *gin.Context) {
	var params role.GetAllRoleQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}

	result, err := h.usecase.GetAllRoleQuery(&params)
	utils.HandleError(c, err)
	utils.HandleResponse(c, result)
}

// GetRolePermissions godoc
// @Summary      Get role permissions
// @Description  Retrieves all permissions assigned to a specific role
// @Tags         role
// @Accept       json
// @Produce      json
// @Param        request  query     role.GetRolePermissionsQuery  true  "Role ID to retrieve permissions for"
// @success      200      {object}  utils.Result                  "List of permissions for the role"
// @Failure      400      {object}  utils.Result                  "Validation error"
// @Failure      401      {object}  utils.Result                  "unauthorized"
// @Failure      404      {object}  utils.Result                  "Role not found"
// @Failure      500      {object}  utils.Result                  "Internal server error"
// @Router       /role/permissions [get]
// @Security BearerAuth
func (h *RoleHandler) GetRolePermissions(c *gin.Context) {
	var params role.GetRolePermissionsQuery
	if !h.validator.ValidateQuery(c, &params) {
		return
	}

	result, err := h.usecase.GetRolePermissionsQuery(&params)
	utils.HandleError(c, err)
	utils.HandleResponse(c, result)
}
