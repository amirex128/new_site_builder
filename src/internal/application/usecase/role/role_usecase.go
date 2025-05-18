package roleusecase

import (
	"time"

	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	"github.com/amirex128/new_site_builder/src/internal/application/dto/role"
	"github.com/amirex128/new_site_builder/src/internal/contract"
	"github.com/amirex128/new_site_builder/src/internal/contract/repository"
	"github.com/amirex128/new_site_builder/src/internal/domain"
)

type RoleUsecase struct {
	logger       sflogger.Logger
	roleRepo     repository.IRoleRepository
	customerRepo repository.ICustomerRepository
	userRepo     repository.IUserRepository
	planRepo     repository.IPlanRepository
}

func NewRoleUsecase(c contract.IContainer) *RoleUsecase {
	return &RoleUsecase{
		logger:       c.GetLogger(),
		roleRepo:     c.GetRoleRepo(),
		customerRepo: c.GetCustomerRepo(),
		userRepo:     c.GetUserRepo(),
		planRepo:     c.GetPlanRepo(),
	}
}

// CreateRoleCommand creates a new role
func (u *RoleUsecase) CreateRoleCommand(params *role.CreateRoleCommand) (any, error) {
	// Check admin access
	// Note: In .NET this was done with gate.IsAdminAccess()

	// Create the role entity
	newRole := domain.Role{
		Name:      *params.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		IsDeleted: false,
	}

	// Create the role
	roleID, err := u.roleRepo.Create(newRole)
	if err != nil {
		return nil, err
	}

	// Assign permissions if specified
	if len(params.PermissionIDs) > 0 {
		for _, permissionID := range params.PermissionIDs {
			err = u.roleRepo.AddPermissionToRole(roleID, permissionID)
			if err != nil {
				// Consider what to do if permission assignment fails
				// For now, we'll continue and try to assign the rest
				u.logger.Errorf("Failed to assign permission %d to role %d: %v", permissionID, roleID, err)
			}
		}
	}

	return map[string]interface{}{
		"id":   roleID,
		"name": newRole.Name,
	}, nil
}

// UpdateRoleCommand updates an existing role
func (u *RoleUsecase) UpdateRoleCommand(params *role.UpdateRoleCommand) (any, error) {
	// Check admin access
	// Note: In .NET this was done with gate.IsAdminAccess()

	// Get the existing role
	existingRole, err := u.roleRepo.GetByID(*params.ID)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if params.Name != nil {
		existingRole.Name = *params.Name
	}

	existingRole.UpdatedAt = time.Now()

	// Update the role
	err = u.roleRepo.Update(existingRole)
	if err != nil {
		return nil, err
	}

	// Update permissions if specified
	if params.PermissionIDs != nil {
		// First, remove all existing permissions
		err = u.roleRepo.RemoveAllPermissionsFromRole(*params.ID)
		if err != nil {
			return nil, err
		}

		// Then add the new permissions
		for _, permissionID := range params.PermissionIDs {
			err = u.roleRepo.AddPermissionToRole(*params.ID, permissionID)
			if err != nil {
				// Log error but continue
				u.logger.Errorf("Failed to assign permission %d to role %d: %v", permissionID, *params.ID, err)
			}
		}
	}

	return existingRole, nil
}

// SetRoleToCustomerCommand assigns roles to a customer
func (u *RoleUsecase) SetRoleToCustomerCommand(params *role.SetRoleToCustomerCommand) (any, error) {
	// Check admin access
	// Note: In .NET this was done with gate.IsAdminAccess()

	// Check if customer exists
	_, err := u.customerRepo.GetByID(*params.CustomerID)
	if err != nil {
		return nil, err
	}

	// First, remove all existing roles from customer
	err = u.roleRepo.RemoveAllRolesFromCustomer(*params.CustomerID)
	if err != nil {
		return nil, err
	}

	// Then add the new roles
	for _, roleName := range params.Role {
		role, err := u.roleRepo.GetByName(roleName)
		if err != nil {
			// Log error and continue
			u.logger.Errorf("Failed to find role with name %s: %v", roleName, err)
			continue
		}

		err = u.roleRepo.AddRoleToCustomer(role.ID, *params.CustomerID)
		if err != nil {
			// Log error and continue
			u.logger.Errorf("Failed to assign role %s to customer %d: %v", roleName, *params.CustomerID, err)
		}
	}

	return map[string]interface{}{
		"customerId": params.CustomerID,
		"roles":      params.Role,
	}, nil
}

// SetRoleToUserCommand assigns roles to a user
func (u *RoleUsecase) SetRoleToUserCommand(params *role.SetRoleToUserCommand) (any, error) {
	// Check admin access
	// Note: In .NET this was done with gate.IsAdminAccess()

	// Check if user exists
	_, err := u.userRepo.GetByID(*params.UserID)
	if err != nil {
		return nil, err
	}

	// First, remove all existing roles from user
	err = u.roleRepo.RemoveAllRolesFromUser(*params.UserID)
	if err != nil {
		return nil, err
	}

	// Then add the new roles
	for _, roleName := range params.Roles {
		role, err := u.roleRepo.GetByName(roleName)
		if err != nil {
			// Log error and continue
			u.logger.Errorf("Failed to find role with name %s: %v", roleName, err)
			continue
		}

		err = u.roleRepo.AddRoleToUser(role.ID, *params.UserID)
		if err != nil {
			// Log error and continue
			u.logger.Errorf("Failed to assign role %s to user %d: %v", roleName, *params.UserID, err)
		}
	}

	return map[string]interface{}{
		"userId": params.UserID,
		"roles":  params.Roles,
	}, nil
}

// SetRoleToPlanCommand assigns roles to a plan
func (u *RoleUsecase) SetRoleToPlanCommand(params *role.SetRoleToPlanCommand) (any, error) {
	// Check admin access
	// Note: In .NET this was done with gate.IsAdminAccess()

	// Check if plan exists
	_, err := u.planRepo.GetByID(*params.PlanID)
	if err != nil {
		return nil, err
	}

	// First, remove all existing roles from plan
	err = u.roleRepo.RemoveAllRolesFromPlan(*params.PlanID)
	if err != nil {
		return nil, err
	}

	// Then add the new roles
	for _, roleName := range params.Roles {
		role, err := u.roleRepo.GetByName(roleName)
		if err != nil {
			// Log error and continue
			u.logger.Errorf("Failed to find role with name %s: %v", roleName, err)
			continue
		}

		err = u.roleRepo.AddRoleToPlan(role.ID, *params.PlanID)
		if err != nil {
			// Log error and continue
			u.logger.Errorf("Failed to assign role %s to plan %d: %v", roleName, *params.PlanID, err)
		}
	}

	return map[string]interface{}{
		"planId": params.PlanID,
		"roles":  params.Roles,
	}, nil
}

// GetAllPermissionQuery gets all permissions with pagination
func (u *RoleUsecase) GetAllPermissionQuery(params *role.GetAllPermissionQuery) (any, error) {
	// Check admin access
	// Note: In .NET this was done with gate.IsAdminAccess()

	// Note: We need to implement a repository method to get all permissions
	// Since we don't have direct access to a permission repository in the container
	permissions, count, err := u.roleRepo.GetAllPermissions(params.PaginationRequestDto)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"items": permissions,
		"total": count,
	}, nil
}

// GetAllRoleQuery gets all roles with pagination
func (u *RoleUsecase) GetAllRoleQuery(params *role.GetAllRoleQuery) (any, error) {
	// Check admin access
	// Note: In .NET this was done with gate.IsAdminAccess()

	result, count, err := u.roleRepo.GetAll(params.PaginationRequestDto)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"items": result,
		"total": count,
	}, nil
}

// GetRolePermissionsQuery gets permissions for roles with pagination
func (u *RoleUsecase) GetRolePermissionsQuery(params *role.GetRolePermissionsQuery) (any, error) {
	// Check admin access
	// Note: In .NET this was done with gate.IsAdminAccess()

	// This is a placeholder - in a real implementation you'd query role-permission mappings
	// We need to add a method to the role repository to get permissions for a role
	permissions, count, err := u.roleRepo.GetRolePermissions(params.PaginationRequestDto)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"items": permissions,
		"total": count,
	}, nil
}
