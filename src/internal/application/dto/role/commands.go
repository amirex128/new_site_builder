package role

// CreateRoleCommand represents a command to create a new role
type CreateRoleCommand struct {
	Name          *string `json:"name" validate:"required_text=1 100"`
	PermissionIDs []int64 `json:"permissionIds" validate:"array_number=1 100 1 0 false"`
}

// SetRoleToCustomerCommand represents a command to set roles to a customer
type SetRoleToCustomerCommand struct {
	Role       []string `json:"role" validate:"array_string=1 50 1 100"`
	CustomerID *int64   `json:"customerId" validate:"required,gt=0"`
}

// SetRoleToPlanCommand represents a command to set roles to a plan
type SetRoleToPlanCommand struct {
	Roles  []string `json:"roles" validate:"array_string=1 50 1 100"`
	PlanID *int64   `json:"planId" validate:"required,gt=0"`
}

// SetRoleToUserCommand represents a command to set roles to a user
type SetRoleToUserCommand struct {
	Roles  []string `json:"roles" validate:"array_string=1 50 1 100"`
	UserID *int64   `json:"userId" validate:"required,gt=0"`
}

// UpdateRoleCommand represents a command to update a role
type UpdateRoleCommand struct {
	ID            *int64  `json:"id" validate:"required,gt=0"`
	Name          *string `json:"name,omitempty" validate:"optional_text=1 200"`
	PermissionIDs []int64 `json:"permissionIds,omitempty" validate:"array_number_optional=0 100 1 0 false"`
}
