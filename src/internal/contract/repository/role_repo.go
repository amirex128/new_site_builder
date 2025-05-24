package repository

import (
	common "github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/domain"
)

type IRoleRepository interface {
	GetAll(paginationRequestDto common.PaginationRequestDto) (common.PaginationResponseDto[domain.Role], int64, error)
	GetByID(id int64) (domain.Role, error)
	GetByName(name string) (domain.Role, error)
	Create(role domain.Role) (int64, error)
	Update(role domain.Role) error
	Delete(id int64) error

	// Role-Permission operations
	AddPermissionToRole(roleID int64, permissionID int64) error
	RemovePermissionFromRole(roleID int64, permissionID int64) error
	RemoveAllPermissionsFromRole(roleID int64) error
	GetAllPermissions(paginationRequestDto common.PaginationRequestDto) ([]domain.Permission, int64, error)
	GetRolePermissions(paginationRequestDto common.PaginationRequestDto) ([]domain.Permission, int64, error)

	// Role-User operations
	AddRoleToUser(roleID int64, userID int64) error
	RemoveRoleFromUser(roleID int64, userID int64) error
	RemoveAllRolesFromUser(userID int64) error

	// Role-Customer operations
	AddRoleToCustomer(roleID int64, customerID int64) error
	RemoveRoleFromCustomer(roleID int64, customerID int64) error
	RemoveAllRolesFromCustomer(customerID int64) error

	// Role-Plan operations
	AddRoleToPlan(roleID int64, planID int64) error
	RemoveRoleFromPlan(roleID int64, planID int64) error
	RemoveAllRolesFromPlan(planID int64) error
}
