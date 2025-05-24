package role

import "github.com/amirex128/new_site_builder/src/internal/contract/common"

// GetAllPermissionQuery represents a query to get all permissions with pagination
type GetAllPermissionQuery struct {
	common.PaginationRequestDto
}

// GetAllRoleQuery represents a query to get all roles with pagination
type GetAllRoleQuery struct {
	common.PaginationRequestDto
}

// GetRolePermissionsQuery represents a query to get role permissions
type GetRolePermissionsQuery struct {
	common.PaginationRequestDto
	RoleID int64 `json:"roleId" form:"roleId" nameFa:"شناسه نقش" validate:"required"`
}
