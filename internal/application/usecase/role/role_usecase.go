package roleusecase

import (
	role2 "github.com/amirex128/new_site_builder/internal/application/dto/role"
	"github.com/amirex128/new_site_builder/internal/application/usecase"
	"github.com/amirex128/new_site_builder/internal/application/utils/resp"
	"github.com/amirex128/new_site_builder/internal/contract"
	repository2 "github.com/amirex128/new_site_builder/internal/contract/repository"
	"github.com/amirex128/new_site_builder/internal/domain"
	"time"
)

type RoleUsecase struct {
	*usecase.BaseUsecase
	roleRepo     repository2.IRoleRepository
	customerRepo repository2.ICustomerRepository
	userRepo     repository2.IUserRepository
	planRepo     repository2.IPlanRepository
}

func NewRoleUsecase(c contract.IContainer) *RoleUsecase {
	return &RoleUsecase{
		BaseUsecase: &usecase.BaseUsecase{
			Logger:      c.GetLogger(),
			AuthContext: c.GetAuthTransientService(),
		},
		roleRepo:     c.GetRoleRepo(),
		customerRepo: c.GetCustomerRepo(),
		userRepo:     c.GetUserRepo(),
		planRepo:     c.GetPlanRepo(),
	}
}

// CreateRoleCommand creates a new role
func (u *RoleUsecase) CreateRoleCommand(params *role2.CreateRoleCommand) (*resp.Response, error) {
	newRole := domain.Role{
		Name:      *params.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		IsDeleted: false,
	}
	roleID, err := u.roleRepo.Create(&newRole)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	if len(params.PermissionIDs) > 0 {
		for _, permissionID := range params.PermissionIDs {
			err = u.roleRepo.AddPermissionToRole(roleID, permissionID)
			if err != nil {
				continue
			}
		}
	}
	return resp.NewResponseData(resp.Created, resp.Data{"id": roleID, "name": newRole.Name}, "نقش با موفقیت ایجاد شد"), nil
}

// UpdateRoleCommand updates an existing role
func (u *RoleUsecase) UpdateRoleCommand(params *role2.UpdateRoleCommand) (*resp.Response, error) {
	existingRole, err := u.roleRepo.GetByID(*params.ID)
	if err != nil {
		return nil, resp.NewError(resp.NotFound, err.Error())
	}
	if params.Name != nil {
		existingRole.Name = *params.Name
	}
	existingRole.UpdatedAt = time.Now()
	err = u.roleRepo.Update(existingRole)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	if params.PermissionIDs != nil {
		err = u.roleRepo.RemoveAllPermissionsFromRole(*params.ID)
		if err != nil {
			return nil, resp.NewError(resp.Internal, err.Error())
		}
		for _, permissionID := range params.PermissionIDs {
			err = u.roleRepo.AddPermissionToRole(*params.ID, permissionID)
			if err != nil {
				continue
			}
		}
	}
	return resp.NewResponseData(resp.Updated, resp.Data{"role": existingRole}, "نقش با موفقیت بروزرسانی شد"), nil
}

// SetRoleToCustomerCommand assigns roles to a customer
func (u *RoleUsecase) SetRoleToCustomerCommand(params *role2.SetRoleToCustomerCommand) (*resp.Response, error) {
	_, err := u.customerRepo.GetByID(*params.CustomerID)
	if err != nil {
		return nil, resp.NewError(resp.NotFound, err.Error())
	}
	err = u.roleRepo.RemoveAllRolesFromCustomer(*params.CustomerID)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	for _, roleName := range params.Role {
		role, err := u.roleRepo.GetByName(roleName)
		if err != nil {
			continue
		}
		err = u.roleRepo.AddRoleToCustomer(role.ID, *params.CustomerID)
		if err != nil {
			continue
		}
	}
	return resp.NewResponseData(resp.Success, resp.Data{"customerId": params.CustomerID, "roles": params.Role}, "نقش‌ها با موفقیت به مشتری اختصاص داده شدند"), nil
}

// SetRoleToUserCommand assigns roles to a user
func (u *RoleUsecase) SetRoleToUserCommand(params *role2.SetRoleToUserCommand) (*resp.Response, error) {
	_, err := u.userRepo.GetByID(*params.UserID)
	if err != nil {
		return nil, resp.NewError(resp.NotFound, err.Error())
	}
	err = u.roleRepo.RemoveAllRolesFromUser(*params.UserID)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	for _, roleName := range params.Roles {
		role, err := u.roleRepo.GetByName(roleName)
		if err != nil {
			continue
		}
		err = u.roleRepo.AddRoleToUser(role.ID, *params.UserID)
		if err != nil {
			continue
		}
	}
	return resp.NewResponseData(resp.Success, resp.Data{"userId": params.UserID, "roles": params.Roles}, "نقش‌ها با موفقیت به کاربر اختصاص داده شدند"), nil
}

// SetRoleToPlanCommand assigns roles to a plan
func (u *RoleUsecase) SetRoleToPlanCommand(params *role2.SetRoleToPlanCommand) (*resp.Response, error) {
	_, err := u.planRepo.GetByID(*params.PlanID)
	if err != nil {
		return nil, resp.NewError(resp.NotFound, err.Error())
	}
	err = u.roleRepo.RemoveAllRolesFromPlan(*params.PlanID)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	for _, roleName := range params.Roles {
		role, err := u.roleRepo.GetByName(roleName)
		if err != nil {
			continue
		}
		err = u.roleRepo.AddRoleToPlan(role.ID, *params.PlanID)
		if err != nil {
			continue
		}
	}
	return resp.NewResponseData(resp.Success, resp.Data{"planId": params.PlanID, "roles": params.Roles}, "نقش‌ها با موفقیت به طرح اختصاص داده شدند"), nil
}

// GetAllPermissionQuery gets all permissions with pagination
func (u *RoleUsecase) GetAllPermissionQuery(params *role2.GetAllPermissionQuery) (*resp.Response, error) {
	permissions, err := u.roleRepo.GetAllPermissions()
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	return resp.NewResponseData(resp.Retrieved, resp.Data{"items": permissions}, "دسترسی‌ها با موفقیت دریافت شدند"), nil
}

// GetAllRoleQuery gets all roles with pagination
func (u *RoleUsecase) GetAllRoleQuery(params *role2.GetAllRoleQuery) (*resp.Response, error) {
	result, err := u.roleRepo.GetAll(params.PaginationRequestDto)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	count := result.TotalCount
	return resp.NewResponseData(resp.Retrieved, resp.Data{"items": result.Items, "total": count}, "نقش‌ها با موفقیت دریافت شدند"), nil
}

// GetRolePermissionsQuery gets permissions for roles with pagination
func (u *RoleUsecase) GetRolePermissionsQuery(params *role2.GetRolePermissionsQuery) (*resp.Response, error) {
	permissions, err := u.roleRepo.GetRolePermissions(params.RoleID)
	if err != nil {
		return nil, resp.NewError(resp.Internal, err.Error())
	}
	return resp.NewResponseData(resp.Retrieved, resp.Data{"items": permissions}, "دسترسی‌های نقش با موفقیت دریافت شدند"), nil
}
