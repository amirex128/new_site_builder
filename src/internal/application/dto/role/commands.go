package role

// CreateRoleCommand represents a command to create a new role
type CreateRoleCommand struct {
	Name          *string `json:"name" validate:"required,max=100" error:"required=نام نقش الزامی است|max=نام نقش نمی‌تواند بیشتر از 100 کاراکتر باشد"`
	PermissionIDs []int64 `json:"permissionIds" validate:"required" error:"required=دسترسی الزامی است"`
}

// SetRoleToCustomerCommand represents a command to set roles to a customer
type SetRoleToCustomerCommand struct {
	Role       []string `json:"role" validate:"required" error:"required=نقش الزامی است"`
	CustomerID *int64   `json:"customerId" validate:"required,gt=0" error:"required=مشتری الزامی است|gt=شناسه مشتری باید بزرگتر از 0 باشد"`
}

// SetRoleToPlanCommand represents a command to set roles to a plan
type SetRoleToPlanCommand struct {
	Roles  []string `json:"roles" validate:"required,min=1" error:"required=نقش‌ها الزامی هستند|min=حداقل یک نقش باید مشخص شود"`
	PlanID *int64   `json:"planId" validate:"required,gt=0" error:"required=پلن الزامی است|gt=شناسه پلن باید بزرگتر از 0 باشد"`
}

// SetRoleToUserCommand represents a command to set roles to a user
type SetRoleToUserCommand struct {
	Roles  []string `json:"roles" validate:"required" error:"required=نقش‌ها الزامی است"`
	UserID *int64   `json:"userId" validate:"required,gt=0" error:"required=کاربر الزامی است|gt=شناسه کاربر باید بزرگتر از 0 باشد"`
}

// UpdateRoleCommand represents a command to update a role
type UpdateRoleCommand struct {
	ID            *int64  `json:"id" validate:"required,gt=0" error:"required=نقش الزامی است|gt=شناسه نقش باید بزرگتر از 0 باشد"`
	Name          *string `json:"name,omitempty" validate:"omitempty,max=200" error:"max=نام نقش نمی‌تواند بیشتر از 200 کاراکتر باشد"`
	PermissionIDs []int64 `json:"permissionIds,omitempty" error:""`
}
