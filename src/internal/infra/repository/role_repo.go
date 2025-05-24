package repository

import (
	common "github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/domain"

	"gorm.io/gorm"
)

type RoleRepo struct {
	database *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *RoleRepo {
	return &RoleRepo{
		database: db,
	}
}

func (r *RoleRepo) GetAll(paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain.Role], error) {
	var roles []domain.Role
	var count int64

	query := r.database.Model(&domain.Role{})
	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&roles)
	if result.Error != nil {
		return nil, result.Error
	}

	return buildPaginationResponse(roles, paginationRequestDto, count)
}

func (r *RoleRepo) GetByID(id int64) (domain.Role, error) {
	var role domain.Role
	result := r.database.First(&role, id)
	if result.Error != nil {
		return role, result.Error
	}
	return role, nil
}

func (r *RoleRepo) GetByName(name string) (domain.Role, error) {
	var role domain.Role
	result := r.database.Where("name = ?", name).First(&role)
	if result.Error != nil {
		return role, result.Error
	}
	return role, nil
}

func (r *RoleRepo) Create(role domain.Role) (int64, error) {
	result := r.database.Create(&role)
	if result.Error != nil {
		return 0, result.Error
	}
	return role.ID, nil
}

func (r *RoleRepo) Update(role domain.Role) error {
	result := r.database.Save(&role)
	return result.Error
}

func (r *RoleRepo) Delete(id int64) error {
	result := r.database.Delete(&domain.Role{}, id)
	return result.Error
}

// Role-Permission operations
func (r *RoleRepo) AddPermissionToRole(roleID int64, permissionID int64) error {
	rolePermission := domain.PermissionRole{
		RoleID:       roleID,
		PermissionID: permissionID,
	}
	result := r.database.Create(&rolePermission)
	return result.Error
}

func (r *RoleRepo) RemovePermissionFromRole(roleID int64, permissionID int64) error {
	result := r.database.Where("role_id = ? AND permission_id = ?", roleID, permissionID).Delete(&domain.PermissionRole{})
	return result.Error
}

func (r *RoleRepo) RemoveAllPermissionsFromRole(roleID int64) error {
	result := r.database.Where("role_id = ?", roleID).Delete(&domain.PermissionRole{})
	return result.Error
}

// GetAllPermissions retrieves all permissions with pagination
func (r *RoleRepo) GetAllPermissions() ([]domain.Permission, error) {
	var permissions []domain.Permission
	var count int64

	query := r.database.Model(&domain.Permission{})
	query.Count(&count)

	result := query.Find(&permissions)
	if result.Error != nil {
		return nil, result.Error
	}

	return permissions, nil
}

// GetRolePermissions retrieves permissions associated with roles
func (r *RoleRepo) GetRolePermissions(roleId int64) ([]domain.Permission, error) {
	var permissions []domain.Permission

	result := r.database.Model(&domain.Permission{}).
		Joins("JOIN permission_roles ON permission_roles.permission_id = permissions.id").
		Where("permission_roles.role_id = ?", roleId).
		Find(&permissions)

	if result.Error != nil {
		return nil, result.Error
	}
	return permissions, nil
}

// Role-User operations
func (r *RoleRepo) AddRoleToUser(roleID int64, userID int64) error {
	userRole := domain.RoleUser{
		RoleID: roleID,
		UserID: userID,
	}
	result := r.database.Create(&userRole)
	return result.Error
}

func (r *RoleRepo) RemoveRoleFromUser(roleID int64, userID int64) error {
	result := r.database.Where("role_id = ? AND user_id = ?", roleID, userID).Delete(&domain.RoleUser{})
	return result.Error
}

func (r *RoleRepo) RemoveAllRolesFromUser(userID int64) error {
	result := r.database.Where("user_id = ?", userID).Delete(&domain.RoleUser{})
	return result.Error
}

// Role-Customer operations
func (r *RoleRepo) AddRoleToCustomer(roleID int64, customerID int64) error {
	customerRole := domain.CustomerRole{
		RoleID:     roleID,
		CustomerID: customerID,
	}
	result := r.database.Create(&customerRole)
	return result.Error
}

func (r *RoleRepo) RemoveRoleFromCustomer(roleID int64, customerID int64) error {
	result := r.database.Where("role_id = ? AND customer_id = ?", roleID, customerID).Delete(&domain.CustomerRole{})
	return result.Error
}

func (r *RoleRepo) RemoveAllRolesFromCustomer(customerID int64) error {
	result := r.database.Where("customer_id = ?", customerID).Delete(&domain.CustomerRole{})
	return result.Error
}

// Role-Plan operations
func (r *RoleRepo) AddRoleToPlan(roleID int64, planID int64) error {
	planRole := domain.RolePlan{
		RoleID: roleID,
		PlanID: planID,
	}
	result := r.database.Create(&planRole)
	return result.Error
}

func (r *RoleRepo) RemoveRoleFromPlan(roleID int64, planID int64) error {
	result := r.database.Where("role_id = ? AND plan_id = ?", roleID, planID).Delete(&domain.RolePlan{})
	return result.Error
}

func (r *RoleRepo) RemoveAllRolesFromPlan(planID int64) error {
	result := r.database.Where("plan_id = ?", planID).Delete(&domain.RolePlan{})
	return result.Error
}
